// Command xsdgen generates Go structs from Travelport UAPI XSD schemas.
//
// It is a purpose-built, multi-package generator. Every XSD schema namespace
// becomes its own Go package, named after the domain without a version suffix
// for single-version domains (e.g. air_v55_0 -> package air; contract upgrades
// keep import paths stable); only the multi-version common family keeps its
// version in the package name (common32/common37/common55/...). Enumeration
// simpleTypes are emitted into a per-domain "enums" sub-package so the scalar
// types and the closed-set enumerations do not bloat a single file.
//
// For each type it:
//   - preserves xs:sequence element order (streaming parse),
//   - embeds xs:extension base types (local or cross-package) so derived
//     structs keep their base fields and their original XML namespaces,
//   - flattens xs:group / xs:attributeGroup,
//   - marks minOccurs=0 fields as pointers and maxOccurs>1 as slices,
//   - models xs:choice members as pointer fields with omitempty,
//   - emits xml tags carrying the full namespace URI,
//   - emits snake_case json tags,
//   - documents each type with its XSD type name.
package main

import (
	"encoding/xml"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// ---------------------------------------------------------------------------
// Ordered XSD node tree (streaming parse preserves document order)
// ---------------------------------------------------------------------------

type Node struct {
	Local    string
	Space    string
	Attrs    map[string]string
	Children []*Node
	Text     string
}

func (n *Node) attr(name string) string { return n.Attrs[name] }

func (n *Node) child(local string) *Node {
	for _, c := range n.Children {
		if c.Local == local {
			return c
		}
	}
	return nil
}

func (n *Node) children(local string) []*Node {
	var out []*Node
	for _, c := range n.Children {
		if c.Local == local {
			out = append(out, c)
		}
	}
	return out
}

func (n *Node) text() string {
	if strings.TrimSpace(n.Text) != "" {
		return strings.TrimSpace(n.Text)
	}
	var b strings.Builder
	for _, c := range n.Children {
		if t := c.text(); t != "" {
			b.WriteString(t)
		}
	}
	return strings.TrimSpace(b.String())
}

func parseXSD(path string) (*Node, map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	dec := xml.NewDecoder(f)
	prefixes := map[string]string{}

	var root *Node
	stack := []*Node{}

	flush := func(start xml.StartElement) {
		for _, a := range start.Attr {
			if a.Name.Space == "xmlns" {
				prefixes[a.Name.Local] = a.Value
			} else if a.Name.Local == "xmlns" && a.Name.Space == "" {
				prefixes[""] = a.Value
			}
		}
	}

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			flush(t)
			node := &Node{Local: t.Name.Local, Space: t.Name.Space, Attrs: map[string]string{}}
			for _, a := range t.Attr {
				if a.Name.Space == "xmlns" || a.Name.Local == "xmlns" {
					continue
				}
				key := a.Name.Local
				if a.Name.Space != "" {
					if uri, ok := prefixes[a.Name.Space]; ok {
						key = "{" + uri + "}" + a.Name.Local
					} else {
						key = a.Name.Space + ":" + a.Name.Local
					}
				}
				node.Attrs[key] = a.Value
			}
			if root == nil {
				root = node
			} else {
				stack[len(stack)-1].Children = append(stack[len(stack)-1].Children, node)
			}
			stack = append(stack, node)
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			if len(stack) > 0 {
				stack[len(stack)-1].Text += string(t)
			}
		}
	}
	if root == nil {
		return nil, nil, fmt.Errorf("no root element in %s", path)
	}
	return root, prefixes, nil
}

// ---------------------------------------------------------------------------
// Schema model
// ---------------------------------------------------------------------------

const xsNS = "http://www.w3.org/2001/XMLSchema"

type fieldDef struct {
	goName  string // exported Go field name (or embedded expr)
	xmlName string // XML local element/attribute name
	xmlNS   string // XML namespace URI for this field
	goType  string // Go type expression (without leading *)
	pointer bool
	slice   bool
	attr    bool
	embed   bool // embedded (anonymous) base struct field
	comment string
}

type complexTypeDef struct {
	name   string
	ns     string
	node   *Node
	doc    string
	fields []fieldDef
}

type attrGroupDef struct {
	name string
	node *Node // schema node; attrs resolved lazily at expansion
}

type groupDef struct {
	name      string
	particles []*Node
}

type simpleDef struct {
	name     string
	ns       string
	isEnum   bool
	baseGo   string // builtin Go type for the restriction base
	enumVals []string
	doc      string
}

type typeKey struct{ ns, local string }

type pkgInfo struct {
	name       string // go package name + output dir base, e.g. "air54"
	ns         string
	dir        string
	importPath string
	enumsName  string // e.g. "air54enums"
	enumsDir   string
	enumsPath  string
	// group points to the Go package that hosts the structs of this namespace.
	// XSD allows namespaces to reference each other recursively
	// (air ↔ rail ↔ universal) while Go forbids import cycles, so such
	// namespaces must be merged into one Go package. nil means this package
	// hosts itself.
	group *pkgInfo
}

// goPkg returns the Go package hosting the structs of this namespace.
func (p *pkgInfo) goPkg() *pkgInfo {
	if p == nil {
		return nil
	}
	if p.group != nil {
		return p.group
	}
	return p
}

type world struct {
	nsToPkg    map[string]*pkgInfo
	complex    map[typeKey]*complexTypeDef
	simple     map[typeKey]*simpleDef
	attrGroups map[typeKey]*attrGroupDef
	groups     map[typeKey]*groupDef
	elemTypes  map[typeKey]string // top-level element name -> its type qname
	// elemSimpleGo maps a top-level element declared with an INLINE anonymous
	// <xs:simpleType> (e.g. SystemTime, Payload) to the builtin Go type of its
	// restriction base. Such elements have no @type and no complexType, so a
	// ref to them would otherwise degrade to interface{}.
	elemSimpleGo map[typeKey]string
	localIndex   map[string]typeKey // local name -> first defining typeKey (chameleon/include fallback)
	mods         []*pkgInfo         // packages in registry order
	// extBases records complexTypes that appear as an xs:extension base. Such
	// types are embedded (flattened) into derived structs and must NOT carry an
	// XMLName field, otherwise Go would marshal a spurious wrapper element.
	extBases map[typeKey]bool
	// rootElems records top-level <xs:element> declarations. Only these types may
	// carry an XMLName field: encoding/xml REJECTS (hard error) a struct whose
	// XMLName tag disagrees with the parent field's tag, and 762 nested fields
	// reference a type under a different element name (e.g. field AirSegment of
	// type TypeBaseAirSegment). Nested elements get their name from the field tag.
	rootElems map[typeKey]bool
	// synth maps an inline anonymous <xs:complexType> node to the synthetic named
	// complexType generated for it, so repeated resolution stays idempotent.
	synth map[*Node]typeKey
	// importedBy records, for each package namespace, the set of other package
	// namespaces that reference at least one of its types. Used to drop orphan
	// "common" packages that nothing else depends on.
	importedBy map[string]map[string]bool
	// structDeps records struct-level cross-namespace type references
	// (fromNS -> set of referenced namespaces). Enum-only references are
	// excluded: enums live in leaf subpackages and can never form a cycle.
	// It drives auto-detection of cyclic namespaces that must be merged.
	structDeps map[string]map[string]bool
	// typeRename holds Go type names that had to be disambiguated because two
	// merged namespaces declare the same local type name.
	typeRename map[typeKey]string
	// attrForm records, per target namespace, whether attributes are
	// namespace-qualified by default (XSD attributeFormDefault="qualified").
	// All 58 Travelport XSDs use "unqualified": locally declared attributes
	// must not carry a namespace, otherwise encoding/xml emits prefix:Attr on
	// Marshal (which the GDS always rejects) and, on Unmarshal, cannot match
	// the unqualified wire attributes and silently drops every attribute value
	// (most business data lives on attributes).
	// Only attributes that are explicitly form="qualified" or global
	// attributes referenced via ref carry a namespace.
	attrForm map[string]bool
}

// goTypeName returns the Go name of a type in a namespace (rewritten on
// name collision inside a merged package).
func (w *world) goTypeName(ns, local string) string {
	if n, ok := w.typeRename[typeKey{ns, local}]; ok {
		return n
	}
	return exportName(local)
}

// markImport records that fromNS references a type owned by importedNS.
func (w *world) markImport(importedNS, fromNS string) {
	if importedNS == "" || fromNS == "" || importedNS == fromNS {
		return
	}
	if w.importedBy[importedNS] == nil {
		w.importedBy[importedNS] = map[string]bool{}
	}
	w.importedBy[importedNS][fromNS] = true
}

// ---------------------------------------------------------------------------
// Built-in XSD type mapping
// ---------------------------------------------------------------------------

var builtinTypes = map[string]string{
	"string": "string", "normalizedString": "string", "token": "string",
	"NMTOKEN": "string", "language": "string", "anyURI": "string",
	"ID": "string", "IDREF": "string", "QName": "string", "Name": "string", "NCName": "string",
	"boolean": "bool", "decimal": "float64", "float": "float64", "double": "float64",
	"integer": "int64", "int": "int", "long": "int64", "short": "int16", "byte": "int8",
	"unsignedInt": "uint", "unsignedLong": "uint64", "unsignedShort": "uint16", "unsignedByte": "uint8",
	"nonNegativeInteger": "int64", "positiveInteger": "int64", "negativeInteger": "int64", "nonPositiveInteger": "int64",
	// Date/time values always map to string, passing through the XSD lexical
	// form verbatim (xs:date looks like 2026-09-01, xs:dateTime like
	// 2026-09-01T12:00:00).
	// They were previously mapped to time.Time, but XSD simpleTypes become
	// named Go types (type TypeDate time.Time), and named types do not
	// inherit time.Time's Marshal/Unmarshal methods while the generator does
	// not supply codecs: JSON decoding then failed outright (cannot unmarshal
	// string into TypeDate) and XML encoding silently emitted an empty
	// element <CheckinDate></CheckinDate>, sending empty dates to the GDS.
	// string also stays consistent with the existing time/duration handling
	// and with the caller-facing JSON contract.
	"dateTime": "string", "date": "string", "time": "string", "duration": "string",
	"base64Binary": "string", "hexBinary": "string",
}

func builtinGo(local string) string {
	if t, ok := builtinTypes[local]; ok {
		return t
	}
	return "string"
}

func splitQName(val string) (uri, local string) {
	if strings.HasPrefix(val, "{") {
		if i := strings.Index(val, "}"); i >= 0 {
			return val[1:i], val[i+1:]
		}
	}
	return "", val
}

// ---------------------------------------------------------------------------
// Reference resolution (cross-package aware)
// ---------------------------------------------------------------------------

// findComplex resolves a type key, falling back to the global local-name index
// (XSD chameleon/include semantics: an unprefixed ref may name a type declared
// in an included/imported schema of a different namespace).
func (w *world) findComplex(uri, local string) *complexTypeDef {
	if c, ok := w.complex[typeKey{uri, local}]; ok {
		return c
	}
	if alt, ok := w.localIndex[local]; ok {
		if c, ok := w.complex[alt]; ok {
			return c
		}
	}
	return nil
}

func (w *world) findSimple(uri, local string) *simpleDef {
	if s, ok := w.simple[typeKey{uri, local}]; ok {
		return s
	}
	if alt, ok := w.localIndex[local]; ok {
		if s, ok := w.simple[alt]; ok {
			return s
		}
	}
	return nil
}

// resolveRef returns the Go type expression for an XSD type reference and, if it
// lives in a different package, the package that must be imported.
func (w *world) resolveRef(val, fromNS string) (expr string, imp *pkgInfo) {
	uri, local := splitQName(val)
	if uri == "" {
		if t, ok := builtinTypes[local]; ok {
			return t, nil
		}
		if alt, ok := w.localIndex[local]; ok {
			uri, local = alt.ns, alt.local
		} else {
			return "interface{}", nil
		}
	}
	if uri == xsNS {
		return builtinGo(local), nil
	}
	if c := w.findComplex(uri, local); c != nil {
		return w.pkgExpr(c.ns, local, fromNS), w.impPkg(c.ns, fromNS)
	}
	if s := w.findSimple(uri, local); s != nil {
		if s.isEnum {
			ep := w.nsToPkg[s.ns]
			if ep == nil {
				return "string", nil
			}
			return ep.enumsName + "." + exportName(local), w.impEnums(ep, fromNS)
		}
		return w.pkgExpr(s.ns, local, fromNS), w.impPkg(s.ns, fromNS)
	}
	// Fallback: not resolvable in any namespace, but the name is an XSD
	// builtin (schemas sometimes write unprefixed base="string"); treat it as
	// a builtin instead of degrading to interface{}.
	if t, ok := builtinTypes[local]; ok {
		return t, nil
	}
	return "interface{}", nil
}

func (w *world) pkgExpr(uri, local, fromNS string) string {
	owner := w.nsToPkg[uri].goPkg()
	from := w.nsToPkg[fromNS].goPkg()
	name := w.goTypeName(uri, local)
	if owner == nil || owner == from {
		return name
	}
	return owner.name + "." + name
}

func (w *world) impPkg(uri, fromNS string) *pkgInfo {
	owner := w.nsToPkg[uri]
	from := w.nsToPkg[fromNS]
	if owner == nil || from == nil || owner == from {
		return nil
	}
	w.markImport(uri, fromNS)
	if w.structDeps[fromNS] == nil {
		w.structDeps[fromNS] = map[string]bool{}
	}
	w.structDeps[fromNS][uri] = true
	if owner.goPkg() == from.goPkg() {
		return nil // same merged package, no import needed
	}
	return owner.goPkg()
}

// impEnums always returns the enums subpackage (it is always a distinct package
// from its parent domain package).
func (w *world) impEnums(ep *pkgInfo, fromNS string) *pkgInfo {
	if ep == nil {
		return nil
	}
	w.markImport(ep.ns, fromNS)
	return &pkgInfo{name: ep.enumsName, importPath: ep.enumsPath}
}

// ---------------------------------------------------------------------------
// Collect schema model
// ---------------------------------------------------------------------------

// qualifyTree rewrites type/ref/base attribute values into "{uri}local" form
// using the prefix->namespace map of the file that declared them.
func qualifyTree(n *Node, filePrefixNS map[string]string, targetNS string) {
	for _, k := range []string{"type", "ref", "base"} {
		if v, ok := n.Attrs[k]; ok && v != "" && !strings.HasPrefix(v, "{") {
			n.Attrs[k] = qualifyVal(v, filePrefixNS, targetNS)
		}
	}
	for _, c := range n.Children {
		qualifyTree(c, filePrefixNS, targetNS)
	}
}

func qualifyVal(val string, filePrefixNS map[string]string, targetNS string) string {
	if val == "" {
		return ""
	}
	if strings.HasPrefix(val, "{") {
		return val
	}
	if i := strings.Index(val, ":"); i >= 0 {
		pre := val[:i]
		local := val[i+1:]
		if uri, ok := filePrefixNS[pre]; ok && uri != "" {
			return "{" + uri + "}" + local
		}
		return "{UNRESOLVED}" + val
	}
	// An unprefixed QName resolves against the default namespace (xmlns=)
	// per XML rules. Travelport schemas declare xmlns=targetNamespace, so
	// ref="Name" names the top-level element Name of this schema, not the
	// xs:Name builtin — checking builtins first would misclassify it into the
	// XMLSchema namespace, fail resolution, and degrade it to interface{}.
	if uri, ok := filePrefixNS[""]; ok && uri != "" {
		return "{" + uri + "}" + val
	}
	if _, ok := builtinTypes[val]; ok {
		return "{" + xsNS + "}" + val
	}
	return "{" + targetNS + "}" + val
}

func (w *world) collect(schema *Node, filePrefixNS map[string]string) {
	ns := schema.attr("targetNamespace")
	pkg := w.nsToPkg[ns]
	if pkg == nil {
		return // unknown namespace (e.g. missing uprofileCommon) -> degrade
	}
	// Records whether attributes of this namespace are qualified by default.
	// Only "qualified" counts as qualified; the default (including
	// "unqualified") is treated as unqualified — the actual semantics of
	// Travelport XSDs.
	w.attrForm[ns] = schema.attr("attributeFormDefault") == "qualified"
	qualifyTree(schema, filePrefixNS, ns)

	for _, ct := range schema.children("complexType") {
		name := ct.attr("name")
		if name == "" {
			continue
		}
		k := typeKey{ns, name}
		if _, dup := w.complex[k]; dup {
			continue
		}
		w.complex[k] = &complexTypeDef{name: name, ns: ns, node: ct, doc: ctDoc(ct)}
		if _, ok := w.localIndex[name]; !ok {
			w.localIndex[name] = k
		}
	}
	for _, st := range schema.children("simpleType") {
		name := st.attr("name")
		if name == "" {
			continue
		}
		k := typeKey{ns, name}
		if _, dup := w.simple[k]; dup {
			continue
		}
		w.simple[k] = w.collectSimple(name, ns, st)
		if _, ok := w.localIndex[name]; !ok {
			w.localIndex[name] = k
		}
	}
	for _, ag := range schema.children("attributeGroup") {
		if n := ag.attr("name"); n != "" {
			k := typeKey{ns, n}
			if _, dup := w.attrGroups[k]; dup {
				continue
			}
			w.attrGroups[k] = &attrGroupDef{name: n, node: ag}
		}
	}
	for _, gr := range schema.children("group") {
		if n := gr.attr("name"); n != "" {
			k := typeKey{ns, n}
			if _, dup := w.groups[k]; dup {
				continue
			}
			w.groups[k] = &groupDef{name: n, particles: directParticles(gr)}
		}
	}
	// Top-level elements: a typed element registers name->type so element refs
	// can resolve to the referenced element's actual type.
	for _, el := range schema.children("element") {
		name := el.attr("name")
		if name == "" {
			continue
		}
		k := typeKey{ns, name}
		w.rootElems[k] = true
		if t := el.attr("type"); t != "" {
			if _, dup := w.elemTypes[k]; !dup {
				w.elemTypes[k] = t
			}
		}
		if st := el.child("simpleType"); st != nil {
			if _, dup := w.elemSimpleGo[k]; !dup {
				w.elemSimpleGo[k] = w.collectSimple(name, ns, st).baseGo
			}
		}
		ct := el.child("complexType")
		if ct == nil {
			continue
		}
		if _, dup := w.complex[k]; dup {
			continue
		}
		w.complex[k] = &complexTypeDef{name: name, ns: ns, node: ct, doc: elDoc(el)}
	}
}

func (w *world) collectSimple(name, ns string, st *Node) *simpleDef {
	sd := &simpleDef{name: name, ns: ns, doc: ctDoc(st)}
	rest := st.child("restriction")
	if rest == nil {
		// list/union not modeled; treat as string alias.
		sd.baseGo = "string"
		return sd
	}
	base := rest.attr("base")
	if base != "" {
		if uri, local := splitQName(base); uri == xsNS {
			sd.baseGo = builtinGo(local)
		} else {
			sd.baseGo = "string"
		}
	} else {
		sd.baseGo = "string"
	}
	hasEnum := false
	for _, en := range rest.children("enumeration") {
		hasEnum = true
		sd.enumVals = append(sd.enumVals, en.attr("value"))
	}
	sd.isEnum = hasEnum
	return sd
}

func directParticles(n *Node) []*Node {
	var out []*Node
	for _, c := range n.Children {
		switch c.Local {
		case "element", "choice", "sequence", "group", "any":
			out = append(out, c)
		}
	}
	return out
}

func ctDoc(n *Node) string {
	if a := n.child("annotation"); a != nil {
		if d := a.child("documentation"); d != nil {
			return sanitizeDoc(d.text())
		}
	}
	return ""
}

func elDoc(n *Node) string {
	if a := n.child("annotation"); a != nil {
		if d := a.child("documentation"); d != nil {
			return sanitizeDoc(d.text())
		}
	}
	return ""
}

func sanitizeDoc(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}

// ---------------------------------------------------------------------------
// Field resolution (with extension embedding + cross-package imports)
// ---------------------------------------------------------------------------

func (w *world) resolveFields(def *complexTypeDef, imports map[string]*pkgInfo) {
	seen := map[string]bool{}
	var out []fieldDef

	var add func(fd fieldDef)
	add = func(fd fieldDef) {
		if fd.goName == "" && !fd.embed {
			return
		}
		if seen[fd.goName] {
			return
		}
		if fd.embed {
			// dedup embedded base by its expr
			seen[fd.goType] = true
		} else {
			seen[fd.goName] = true
		}
		out = append(out, fd)
	}

	w.collectComplex(def.node, def.ns, add, imports)
	def.fields = out
}

func (w *world) collectComplex(ct *Node, fromNS string, add func(fieldDef), imports map[string]*pkgInfo) {
	if sc := ct.child("simpleContent"); sc != nil {
		if ext := sc.child("extension"); ext != nil {
			base := ext.attr("base")
			if base != "" {
				gt, _ := w.resolveRef(base, fromNS)
				add(fieldDef{goName: "Value", xmlName: ",chardata", goType: gt, comment: "XSD simpleContent value"})
			}
			for _, a := range w.collectAttrs(ext, fromNS, imports) {
				add(a)
			}
		}
		return
	}

	if cc := ct.child("complexContent"); cc != nil {
		if ext := cc.child("extension"); ext != nil {
			base := ext.attr("base")
			if base != "" {
				if uri, local := splitQName(base); uri != "" {
					w.extBases[typeKey{uri, local}] = true
				} else if alt, ok := w.localIndex[local]; ok {
					w.extBases[alt] = true
				}
			}
			if base != "" {
				uri, local := splitQName(base)
				if uri == xsNS {
					gt, _ := w.resolveRef(base, fromNS)
					add(fieldDef{goName: "Value", xmlName: ",chardata", goType: gt, comment: "XSD simpleContent value"})
				} else if c := w.findComplex(uri, local); c != nil {
					expr, imp := w.pkgExpr(c.ns, local, fromNS), w.impPkg(c.ns, fromNS)
					if imp != nil {
						imports[imp.importPath] = imp
					}
					add(fieldDef{embed: true, goName: expr, goType: expr, comment: "XSD extension base " + local})
				} else if s := w.findSimple(uri, local); s != nil && !s.isEnum {
					gt, _ := w.resolveRef(base, fromNS)
					add(fieldDef{goName: "Value", xmlName: ",chardata", goType: gt, comment: "XSD simpleContent value"})
				} else {
					add(fieldDef{goName: "Value", xmlName: ",chardata", goType: "interface{}", comment: "XSD extension base (unresolved)"})
				}
			}
			w.collectContent(ext, fromNS, add, imports)
		}
		if rest := cc.child("restriction"); rest != nil {
			w.collectContent(rest, fromNS, add, imports)
		}
		return
	}

	w.collectContent(ct, fromNS, add, imports)
}

func (w *world) collectContent(n *Node, fromNS string, add func(fieldDef), imports map[string]*pkgInfo) {
	for _, a := range w.collectAttrs(n, fromNS, imports) {
		add(a)
	}
	for _, p := range directParticles(n) {
		w.collectParticle(p, fromNS, false, add, imports)
	}
}

func (w *world) collectParticle(p *Node, fromNS string, forceOptional bool, add func(fieldDef), imports map[string]*pkgInfo) {
	switch p.Local {
	case "element":
		fd := w.fieldFromElement(p, fromNS, forceOptional, imports)
		add(fd)
	case "sequence":
		for _, c := range directParticles(p) {
			w.collectParticle(c, fromNS, forceOptional, add, imports)
		}
	case "choice":
		for _, c := range directParticles(p) {
			switch c.Local {
			case "element":
				add(w.fieldFromElement(c, fromNS, true, imports))
			case "sequence":
				for _, cc := range directParticles(c) {
					w.collectParticle(cc, fromNS, true, add, imports)
				}
			case "choice":
				w.collectParticle(c, fromNS, true, add, imports)
			case "group":
				w.expandGroup(c, fromNS, true, add, imports)
			}
		}
	case "group":
		w.expandGroup(p, fromNS, forceOptional, add, imports)
	case "any":
		add(fieldDef{goName: "Any", xmlName: ",any", goType: "[]byte", comment: "XSD wildcard any", xmlNS: fromNS})
	}
}

func (w *world) expandGroup(gr *Node, fromNS string, forceOptional bool, add func(fieldDef), imports map[string]*pkgInfo) {
	ref := gr.attr("ref")
	if ref == "" {
		for _, c := range directParticles(gr) {
			w.collectParticle(c, fromNS, forceOptional, add, imports)
		}
		return
	}
	// Same as collectAttrs: groups must be looked up in the referenced
	// namespace, otherwise cross-namespace <xs:group ref="common:xxx"/> would
	// be silently dropped. An empty uri means a same-namespace reference.
	refNS, local := splitQName(ref)
	if refNS == "" {
		refNS = fromNS
	}
	if gdef, ok := w.groups[typeKey{refNS, local}]; ok {
		minOccurs, _ := strconv.Atoi(gr.attr("minOccurs"))
		optional := forceOptional || minOccurs == 0
		for _, p := range gdef.particles {
			if p.Local == "element" {
				add(w.fieldFromElement(p, fromNS, optional, imports))
			} else {
				w.collectParticle(p, fromNS, forceOptional, add, imports)
			}
		}
	}
}

func (w *world) collectAttrs(n *Node, fromNS string, imports map[string]*pkgInfo) []fieldDef {
	var out []fieldDef
	for _, a := range n.children("attribute") {
		out = append(out, w.attrFromNode(a, fromNS, imports))
	}
	for _, ag := range n.children("attributeGroup") {
		if ref := ag.attr("ref"); ref != "" {
			// qualifyTree normalized refs to {uri}local; lookups must key on
			// the REFERENCED namespace. An earlier implementation dropped uri
			// and always looked up fromNS, so cross-namespace refs like
			// ref="common:providerReservation" were never found and silently
			// dropped (HotelRetrieveReq lost its required
			// ProviderCode/ProviderLocatorCode and the GDS only rejected the
			// message after it was sent).
			// An empty uri means a same-namespace reference; fall back to fromNS.
			refNS, local := splitQName(ref)
			if refNS == "" {
				refNS = fromNS
			}
			if grp, ok := w.attrGroups[typeKey{refNS, local}]; ok && grp.node != nil {
				// Keep passing fromNS (not refNS) here: fromNS doubles as the
				// "destination package" that decides whether the type
				// expression needs package qualification and an import.
				// Passing refNS would degrade attribute types from common into
				// bare identifiers and miss the import.
				// Attribute qualification is unaffected — all 58 schemas use
				// attributeFormDefault="unqualified", so both choices agree
				// (xmlNS stays empty).
				for _, a := range w.collectAttrs(grp.node, fromNS, imports) {
					out = append(out, a)
				}
			}
		}
	}
	return out
}

func (w *world) attrFromNode(a *Node, fromNS string, imports map[string]*pkgInfo) fieldDef {
	ref := a.attr("ref")
	name := a.attr("name")
	xmlNS := fromNS
	if name == "" && ref != "" {
		// Global attributes referenced via ref: always live in their
		// targetNamespace (ruri) and are qualified.
		ruri, rlocal := splitQName(ref)
		name = rlocal
		if ruri != "" {
			xmlNS = ruri
		}
	} else if name != "" {
		// Locally declared attributes: qualification follows
		// attributeFormDefault and the per-attribute form override.
		// When neither is "qualified" the attribute is unqualified (no
		// namespace) — the real semantics of Travelport XSDs; anything else
		// breaks encoding/xml on both send and receive.
		qualified := w.attrForm[fromNS] || a.attr("form") == "qualified"
		if !qualified {
			xmlNS = ""
		}
	}
	use := a.attr("use")
	optional := use != "required"
	typ := a.attr("type")
	if typ == "" && ref != "" {
		ruri, rlocal := splitQName(ref)
		if t, ok := w.elemTypes[typeKey{ruri, rlocal}]; ok {
			typ = t
		}
	}
	// Locally declared attributes may carry an inline anonymous
	// <xs:simpleType> (e.g. MaxResults of ReferenceDataSearchModifiers:
	// <xs:attribute name="MaxResults"><xs:simpleType>
	// <xs:restriction base="xs:integer">…). XSD attributes are always
	// simpleTypes, and the base type must be mapped faithfully
	// (integer→int64); otherwise the fallback below degrades it to string
	// and JSON numbers fail to unmarshal ("cannot unmarshal number into Go
	// struct field … of type string").
	// collectSimple is reused so this shares one code path with inline
	// element simpleTypes.
	if typ == "" {
		if st := a.child("simpleType"); st != nil {
			gt := w.collectSimple(name, fromNS, st).baseGo
			ptr := optional
			return fieldDef{
				goName:  exportName(name),
				xmlName: name,
				xmlNS:   xmlNS,
				goType:  gt,
				pointer: ptr,
				attr:    true,
				comment: fmt.Sprintf("XSD attribute %s", xmlName(name)),
			}
		}
	}
	gt, imp := w.resolveRef(typ, fromNS)
	pointer := optional
	if imp != nil {
		imports[imp.importPath] = imp
	}
	if typ == "" || gt == "interface{}" {
		// XSD attributes can only be simpleTypes, so whether the declaration
		// is missing or unresolvable across namespaces, fall back to string
		// (interface{} can neither be marshaled nor unmarshaled by
		// encoding/xml).
		gt, pointer = "string", optional
	}
	if optional {
		pointer = true
	}
	return fieldDef{
		goName:  exportName(name),
		xmlName: name,
		xmlNS:   xmlNS,
		goType:  gt,
		pointer: pointer,
		attr:    true,
		comment: fmt.Sprintf("XSD attribute %s", xmlName(name)),
	}
}

// synthInlineType synthesizes and registers a named complexType for an
// inline anonymous <xs:complexType> so it can emit a struct like any other
// type instead of degrading to interface{} (encoding/xml can neither marshal
// nor unmarshal interface{}, which would crash requests / yield nil
// responses).
// Synthesized types only occur in nested positions, so they get no XMLName;
// the element name comes from the parent field's xml tag.
func (w *world) synthInlineType(fromNS, elemName string, ct *Node) string {
	if k, ok := w.synth[ct]; ok {
		return w.pkgExpr(k.ns, k.local, fromNS)
	}
	base := exportName(elemName) + "Inline"
	local := base
	for i := 2; ; i++ {
		if _, exists := w.complex[typeKey{fromNS, local}]; !exists {
			break
		}
		local = fmt.Sprintf("%s%d", base, i)
	}
	k := typeKey{fromNS, local}
	w.complex[k] = &complexTypeDef{
		name: local,
		ns:   fromNS,
		node: ct,
		doc:  fmt.Sprintf("Synthesized from the inline anonymous complexType of element %s.", elemName),
	}
	w.synth[ct] = k
	return w.pkgExpr(fromNS, local, fromNS)
}

func (w *world) fieldFromElement(el *Node, fromNS string, optionalOverride bool, imports map[string]*pkgInfo) fieldDef {
	ref := el.attr("ref")
	name := el.attr("name")
	xmlNS := fromNS
	if name == "" && ref != "" {
		ruri, rlocal := splitQName(ref)
		name = rlocal
		if ruri != "" {
			xmlNS = ruri
		}
	}
	minOccurs, _ := strconv.Atoi(el.attr("minOccurs"))
	if el.attr("minOccurs") == "" {
		minOccurs = 1
	}
	maxOccurs := el.attr("maxOccurs")
	optional := minOccurs == 0 || optionalOverride

	pointer := optional
	slice := false
	if maxOccurs == "unbounded" || maxOccurs == "2" || maxOccurs == "3" || maxOccurs == "4" {
		slice = true
		pointer = false
	} else if n, err := strconv.Atoi(maxOccurs); err == nil && n > 1 {
		slice = true
		pointer = false
	}

	ct := el.child("complexType")
	if ct != nil {
		if ext := ct.complexContentExtension(); ext != "" {
			gt, _ := w.resolveRef(ext, fromNS)
			return fieldDef{goName: exportName(name), xmlName: name, xmlNS: xmlNS, goType: gt, pointer: pointer, slice: slice, comment: "XSD element " + xmlName(name)}
		}
		if simpleContent := ct.simpleContentExtension(); simpleContent != "" {
			gt, _ := w.resolveRef(simpleContent, fromNS)
			return fieldDef{goName: exportName(name), xmlName: name, xmlNS: xmlNS, goType: gt, pointer: pointer, slice: slice, comment: "XSD element " + xmlName(name)}
		}
		// Inline anonymous complexType (no extension/simpleContent):
		// synthesize a named type instead of degrading to interface{},
		// which would crash xml.Marshal and lose data in xml.Unmarshal.
		synth := w.synthInlineType(fromNS, name, ct)
		return fieldDef{goName: exportName(name), xmlName: name, xmlNS: xmlNS, goType: synth, pointer: pointer, slice: slice, comment: "XSD element " + xmlName(name)}
	}

	// Resolve the element's type. A bare ref resolves to the referenced
	// element's declared type (tracked in elemTypes). Many Travelport elements
	// (e.g. BillingPointOfSaleInfo, Keyword, HotelProperty) are declared as a
	// top-level <xs:element> with an INLINE <xs:complexType> and no @type, so
	// they are registered as complexTypes keyed by their element name but have
	// no elemTypes entry. Fall back to a global complexType lookup by the ref's
	// (namespace, local) so these resolve to real Go types instead of interface{}.
	// Inline anonymous <xs:simpleType> (restriction of xs:string/xs:dateTime
	// etc.): map the restriction base directly to its Go builtin instead of
	// degrading to interface{}.
	if st := el.child("simpleType"); st != nil {
		gt := w.collectSimple(name, fromNS, st).baseGo
		if gt == "" {
			gt = "string"
		}
		if optional && !slice {
			pointer = true
		}
		return fieldDef{goName: exportName(name), xmlName: name, xmlNS: xmlNS, goType: gt, pointer: pointer, slice: slice, comment: "XSD element " + xmlName(name)}
	}

	typ := el.attr("type")
	if typ == "" && ref != "" {
		ruri, rlocal := splitQName(ref)
		if t, ok := w.elemTypes[typeKey{ruri, rlocal}]; ok {
			typ = t
		} else if c := w.findComplex(ruri, rlocal); c != nil {
			typ = "{" + c.ns + "}" + c.name
		} else if gt, ok := w.elemSimpleGo[typeKey{ruri, rlocal}]; ok {
			// References a "top-level element + inline simpleType" declaration
			// (e.g. ref="SystemTime").
			if gt == "" {
				gt = "string"
			}
			if optional && !slice {
				pointer = true
			}
			return fieldDef{goName: exportName(name), xmlName: name, xmlNS: xmlNS, goType: gt, pointer: pointer, slice: slice, comment: "XSD element " + xmlName(name)}
		}
	}
	gt, imp := w.resolveRef(typ, fromNS)
	if imp != nil {
		imports[imp.importPath] = imp
	}
	if typ == "" {
		// Neither @type nor an inline type (xs:anyType); in practice these
		// are all plain-text elements, treated as string.
		gt = "string"
	}
	if optional && !slice {
		pointer = true
	}
	return fieldDef{goName: exportName(name), xmlName: name, xmlNS: xmlNS, goType: gt, pointer: pointer, slice: slice, comment: "XSD element " + xmlName(name)}
}

// ---------------------------------------------------------------------------
// Emit
// ---------------------------------------------------------------------------

func (w *world) xmlTag(f fieldDef) string {
	if f.xmlName == ",chardata" {
		return ",chardata"
	}
	if f.xmlName == ",any" {
		return ",any"
	}
	ns := f.xmlNS
	if f.attr {
		if ns == "" {
			return f.xmlName + ",attr"
		}
		return fmt.Sprintf("%s %s,attr", ns, f.xmlName)
	}
	tag := fmt.Sprintf("%s %s", ns, f.xmlName)
	if f.pointer || f.slice {
		tag += ",omitempty"
	}
	return tag
}

// groupNamespaces returns the set of all namespaces whose structs are hosted
// by pkg.
func (w *world) groupNamespaces(pkg *pkgInfo) map[string]bool {
	out := map[string]bool{}
	for _, p := range w.mods {
		if p.goPkg() == pkg {
			out[p.ns] = true
		}
	}
	return out
}

// mergeCyclicNamespaces auto-detects mutually recursive namespaces (strongly
// connected components) and merges them into a single Go package. XSD allows
// references like air ↔ rail ↔ universal while Go forbids import cycles, so
// they must be merged; the representative of a merge is the package with the
// most types, keeping existing import paths (e.g. util54 depending on air54)
// working.
func (w *world) mergeCyclicNamespaces() [][]*pkgInfo {
	idx := map[string]int{}
	low := map[string]int{}
	onStack := map[string]bool{}
	var stack []string
	counter := 0
	var sccs [][]string

	var strong func(v string)
	strong = func(v string) {
		idx[v], low[v] = counter, counter
		counter++
		stack = append(stack, v)
		onStack[v] = true
		deps := make([]string, 0, len(w.structDeps[v]))
		for d := range w.structDeps[v] {
			deps = append(deps, d)
		}
		sort.Strings(deps)
		for _, d := range deps {
			if w.nsToPkg[d] == nil {
				continue
			}
			if _, seen := idx[d]; !seen {
				strong(d)
				if low[d] < low[v] {
					low[v] = low[d]
				}
			} else if onStack[d] {
				if idx[d] < low[v] {
					low[v] = idx[d]
				}
			}
		}
		if low[v] == idx[v] {
			var comp []string
			for {
				n := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[n] = false
				comp = append(comp, n)
				if n == v {
					break
				}
			}
			if len(comp) > 1 {
				sort.Strings(comp)
				sccs = append(sccs, comp)
			}
		}
	}
	for _, p := range w.mods {
		if _, seen := idx[p.ns]; !seen {
			strong(p.ns)
		}
	}

	var merged [][]*pkgInfo
	for _, comp := range sccs {
		// Representative = the package with the most types (name order as
		// tiebreaker); other members no longer emit separately after merging.
		var rep *pkgInfo
		repCount := -1
		for _, ns := range comp {
			p := w.nsToPkg[ns]
			n := countComplex(w, p)
			if n > repCount || (n == repCount && p.name < rep.name) {
				rep, repCount = p, n
			}
		}
		group := []*pkgInfo{rep}
		for _, ns := range comp {
			if p := w.nsToPkg[ns]; p != rep {
				p.group = rep
				group = append(group, p)
			}
		}
		merged = append(merged, group)
	}
	return merged
}

// resolveGroupCollisions resolves type-name collisions inside a merged
// package: the representative keeps its names, while colliding types from
// other members get a package prefix (e.g. rail54's Characteristic ->
// Rail54Characteristic).
func (w *world) resolveGroupCollisions() {
	byGoPkg := map[*pkgInfo][]*pkgInfo{}
	for _, p := range w.mods {
		gp := p.goPkg()
		byGoPkg[gp] = append(byGoPkg[gp], p)
	}
	for gp, members := range byGoPkg {
		if len(members) < 2 {
			continue
		}
		// Representative first, the rest in registration order for stability.
		sort.SliceStable(members, func(i, j int) bool { return members[i] == gp && members[j] != gp })
		taken := map[string]bool{}
		for _, m := range members {
			var locals []typeKey
			for k := range w.complex {
				if k.ns == m.ns {
					locals = append(locals, k)
				}
			}
			for k, s := range w.simple {
				if !s.isEnum && k.ns == m.ns {
					locals = append(locals, typeKey{k.ns, k.local})
				}
			}
			sort.Slice(locals, func(i, j int) bool {
				return w.goTypeName(locals[i].ns, locals[i].local) < w.goTypeName(locals[j].ns, locals[j].local)
			})
			for _, k := range locals {
				name := exportName(k.local)
				if taken[name] {
					name = exportName(m.name) + exportName(k.local)
					w.typeRename[typeKey{k.ns, k.local}] = name
				}
				taken[name] = true
			}
		}
	}
}

// injectMethod generates InjectInfrastructure on request base types: the
// server side only back-fills the call's trace identifier (TraceId) when the
// request body leaves it empty. Fields such as BillingPointOfSaleInfo,
// TargetBranch and ProviderCode are no longer injected by code — callers
// (API users) must provide them explicitly in the request body, because they
// are business/authorization credentials that must not be hardcoded or
// auto-filled from runtime configuration.
//
// Backfill semantics (aligned with observability best practice): the gateway
// respects a TraceId business value already present in the request body and
// only falls back to the global trace_id when it is absent, keeping upstream
// echo and troubleshooting chains working. The trace_id travels independently
// in the X-Trace-Id HTTP header rather than the request body, separate from
// the body's business TraceId field.
func (w *world) injectMethod(recv string, hasTrace, traceIsPtr bool) string {
	var sb strings.Builder
	sb.WriteString("// InjectInfrastructure back-fills the call's global trace identifier (trace_id)\n")
	sb.WriteString("// when the request body's TraceId is empty. If the caller provided TraceId\n")
	sb.WriteString("// explicitly, its business value is respected and never overwritten.\n")
	sb.WriteString("// Billing point of sale, target branch and similar fields are NOT injected\n")
	sb.WriteString("// here; callers must provide them explicitly in the request.\n")
	sb.WriteString(fmt.Sprintf("func (b *%s) InjectInfrastructure(traceID string) {\n", recv))

	if hasTrace {
		sb.WriteString("\tif traceID != \"\" && b.TraceId == nil {\n")
		if traceIsPtr {
			sb.WriteString("\t\tb.TraceId = &traceID\n")
		} else {
			sb.WriteString("\t\tb.TraceId = traceID\n")
		}
		sb.WriteString("\t}\n")
	}
	sb.WriteString("}\n\n")
	return sb.String()
}

// emitPackage writes the structs + scalar simpleType aliases for a package.
func (w *world) emitPackage(pkg *pkgInfo) (string, bool) {
	imports := map[string]*pkgInfo{}
	var structs, scalars []string
	needsTime := false
	needsXML := false

	// One Go package may host several namespaces (mutually recursive XSD
	// namespaces get merged); group marks every namespace this Go package
	// actually hosts (including merged-in rail/universal).
	group := w.groupNamespaces(pkg)
	var cnames, snames []typeKey
	for k := range w.complex {
		if group[k.ns] {
			cnames = append(cnames, k)
		}
	}
	for k, s := range w.simple {
		if !s.isEnum && group[k.ns] {
			snames = append(snames, k)
		}
	}
	byGoName := func(list []typeKey) func(i, j int) bool {
		return func(i, j int) bool {
			ni, nj := w.goTypeName(list[i].ns, list[i].local), w.goTypeName(list[j].ns, list[j].local)
			if ni != nj {
				return ni < nj
			}
			return list[i].ns < list[j].ns
		}
	}
	sort.Slice(cnames, byGoName(cnames))
	sort.Slice(snames, byGoName(snames))

	for _, k := range snames {
		sd := w.simple[k]
		gn := w.goTypeName(k.ns, k.local)
		scalars = append(scalars, fmt.Sprintf("// %s corresponds to XSD simpleType %q.\ntype %s %s\n\n", k.local, k.local, gn, sd.baseGo))
	}

	for _, k := range cnames {
		local, ns := k.local, k.ns
		def := w.complex[k]
		w.resolveFields(def, imports)
		gn := w.goTypeName(ns, local)
		// Only types declared as top-level <xs:element> get an XMLName
		// (carrying the namespace, for the SOAP root node). Nested types
		// never do: encoding/xml hard-errors when XMLName disagrees with the
		// parent field tag's element name, and many field element names
		// differ from their type name (e.g. field AirSegment has type
		// TypeBaseAirSegment).
		// xs:extension bases are flattened into derived structs and equally
		// cannot carry an XMLName.
		isRoot := w.rootElems[k] && !w.extBases[k]
		var sb strings.Builder
		if def.doc != "" {
			sb.WriteString(fmt.Sprintf("// %s corresponds to XSD complexType %q.\n// %s\n", gn, local, def.doc))
		} else {
			sb.WriteString(fmt.Sprintf("// %s corresponds to XSD complexType %q.\n", gn, local))
		}
		sb.WriteString(fmt.Sprintf("type %s struct {\n", gn))
		if isRoot {
			needsXML = true
			sb.WriteString(fmt.Sprintf("\tXMLName xml.Name `xml:\"%s %s\" json:\"-\"` // SOAP root node (namespace-qualified)\n", ns, local))
		}
		var billing *fieldDef
		traceIsPtr := false
		hasTrace := false
		for i, f := range def.fields {
			if f.goName == "BillingPointOfSaleInfo" && f.goType != "interface{}" && !f.slice {
				billing = &def.fields[i]
			}
			if f.goName == "TraceId" {
				hasTrace, traceIsPtr = true, f.pointer
			}
			if f.embed {
				sb.WriteString(fmt.Sprintf("\t%s\n", f.goType))
				if strings.Contains(f.goType, "time.Time") {
					needsTime = true
				}
				continue
			}
			tag := w.xmlTag(f)
			json := camelCase(f.goName)
			if f.pointer || f.slice {
				json += ",omitempty"
			}
			if f.xmlName == ",chardata" {
				json = "value,omitempty"
			}
			if f.xmlName == ",any" {
				json = "-,omitempty"
			}
			full := fmt.Sprintf("`xml:\"%s\" json:\"%s\"`", tag, json)
			typ := f.goType
			if f.slice {
				typ = "[]" + typ
			} else if f.pointer {
				typ = "*" + typ
			}
			if strings.Contains(typ, "time.Time") {
				needsTime = true
			}
			if f.comment != "" {
				sb.WriteString(fmt.Sprintf("\t%s %s %s // %s\n", f.goName, typ, full, f.comment))
			} else {
				sb.WriteString(fmt.Sprintf("\t%s %s %s\n", f.goName, typ, full))
			}
		}
		sb.WriteString("}\n\n")
		structs = append(structs, sb.String())
		if local == "BaseCoreReq" && billing != nil {
			structs = append(structs, w.injectMethod(gn, hasTrace, traceIsPtr))
		}
	}

	if len(structs) == 0 && len(scalars) == 0 {
		return "", false
	}

	var b strings.Builder
	b.WriteString("// Code generated by tools/airxsdgen. DO NOT EDIT.\n\n")
	b.WriteString("package " + pkg.name + "\n\n")

	if needsTime || needsXML || len(imports) > 0 {
		b.WriteString("import (\n")
		if needsXML {
			b.WriteString("\t\"encoding/xml\"\n")
		}
		if needsTime {
			b.WriteString("\t\"time\"\n")
		}
		// stable import order
		paths := make([]string, 0, len(imports))
		for p := range imports {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		for _, p := range paths {
			b.WriteString(fmt.Sprintf("\t%q\n", p))
		}
		b.WriteString(")\n\n")
	}
	b.WriteString(strings.Join(scalars, ""))
	b.WriteString(strings.Join(structs, ""))
	return b.String(), true
}

// emitEnums writes the enumeration simpleTypes for a package's enums subpackage.
func (w *world) emitEnums(pkg *pkgInfo) (string, bool) {
	imports := map[string]*pkgInfo{}
	var enames []string
	for k, s := range w.simple {
		if k.ns == pkg.ns && s.isEnum {
			enames = append(enames, k.local)
		}
	}
	sort.Strings(enames)
	if len(enames) == 0 {
		return "", false
	}

	var blocks []string
	for _, local := range enames {
		sd := w.simple[typeKey{pkg.ns, local}]
		gn := exportName(local)
		var sb strings.Builder
		doc := sd.doc
		if doc != "" {
			sb.WriteString(fmt.Sprintf("// %s corresponds to XSD simpleType %q (enumeration).\n// %s\n", gn, local, doc))
		} else {
			sb.WriteString(fmt.Sprintf("// %s corresponds to XSD simpleType %q (enumeration).\n", gn, local))
		}
		sb.WriteString(fmt.Sprintf("type %s %s\n\n", gn, sd.baseGo))
		if len(sd.enumVals) > 0 {
			sb.WriteString("const (\n")
			seen := map[string]bool{}
			for _, v := range sd.enumVals {
				cn := enumConstName(gn, v, seen)
				if sd.baseGo == "string" {
					sb.WriteString(fmt.Sprintf("\t%s %s = %s\n", cn, gn, strconv.Quote(v)))
				} else {
					sb.WriteString(fmt.Sprintf("\t%s %s = %s\n", cn, gn, v))
				}
			}
			sb.WriteString(")\n\n")
		}
		blocks = append(blocks, sb.String())
	}

	var b strings.Builder
	b.WriteString("// Code generated by tools/airxsdgen. DO NOT EDIT.\n\n")
	b.WriteString("package " + pkg.enumsName + "\n\n")
	if len(imports) > 0 {
		b.WriteString("import (\n")
		paths := make([]string, 0, len(imports))
		for p := range imports {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		for _, p := range paths {
			b.WriteString(fmt.Sprintf("\t%q\n", p))
		}
		b.WriteString(")\n\n")
	}
	b.WriteString(strings.Join(blocks, ""))
	return b.String(), true
}

func enumConstName(typeName, value string, seen map[string]bool) string {
	var sb strings.Builder
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(r)
		} else {
			sb.WriteByte('_')
		}
	}
	suffix := sb.String()
	if suffix == "" {
		suffix = "Empty"
	}
	runes := []rune(suffix)
	if unicode.IsDigit(runes[0]) {
		suffix = "_" + suffix
	} else if unicode.IsLower(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
		suffix = string(runes)
	}
	name := typeName + suffix
	for seen[name] {
		name = name + "_"
	}
	seen[name] = true
	return name
}

func xmlName(local string) string { return local }

func (n *Node) complexContentExtension() string {
	if cc := n.child("complexContent"); cc != nil {
		if ext := cc.child("extension"); ext != nil {
			return ext.attr("base")
		}
	}
	return ""
}

func (n *Node) simpleContentExtension() string {
	if sc := n.child("simpleContent"); sc != nil {
		if ext := sc.child("extension"); ext != nil {
			return ext.attr("base")
		}
	}
	return ""
}

func exportName(local string) string {
	if local == "" {
		return ""
	}
	runes := []rune(local)
	if !unicode.IsUpper(runes[0]) {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return string(runes)
}

func snakeCase(name string) string {
	if name == "" {
		return ""
	}
	var b strings.Builder
	runes := []rune(name)
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1])) {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// camelCase converts an exported Go field name to camelCase (lowercase first
// letter) for the public JSON contract.
// E.g. BillingPointOfSaleInfo -> billingPointOfSaleInfo, CacheName ->
// cacheName.
func camelCase(name string) string {
	if name == "" {
		return ""
	}
	runes := []rune(name)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// ---------------------------------------------------------------------------
// Registry + main
// ---------------------------------------------------------------------------

type mod struct {
	dir  string // wsdl subdir
	ns   string // data namespace
	name string // go package base name
}

func schemaMods() []mod {
	const base = "http://www.travelport.com/schema/"
	// Package naming convention (ADR-004): single-version domains use
	// versionless package names (air_v55_0 → air) so import paths stay
	// stable across contract upgrades; only the multi-version common family
	// keeps its version (common32/33/37/55…, each legacy domain pinned to
	// its own common version).
	return []mod{
		{"air_v55_0", base + "air_v55_0", "air"},
		{"hotel_v55_0", base + "hotel_v55_0", "hotel"},
		{"rail_v55_0", base + "rail_v55_0", "rail"},
		{"universal_v55_0", base + "universal_v55_0", "universal"},
		{"cruise_v55_0", base + "cruise_v55_0", "cruise"},
		{"gdsQueue_v55_0", base + "gdsQueue_v55_0", "gdsqueue"},
		{"passive_v55_0", base + "passive_v55_0", "passive"},
		{"vehicle_v55_0", base + "vehicle_v55_0", "vehicle"},
		{"sharedBooking_v55_0", base + "sharedBooking_v55_0", "sharedbooking"},
		{"util_v55_0", base + "util_v55_0", "util"},
		{"common_v55_0", base + "common_v55_0", "common55"},
		{"common_v32_0", base + "common_v32_0", "common32"},
		{"common_v33_0", base + "common_v33_0", "common33"},
		{"common_v34_0", base + "common_v34_0", "common34"},
		{"common_v37_0", base + "common_v37_0", "common37"},
		{"common_v38_0", base + "common_v38_0", "common38"},
		{"system_v32_0", base + "system_v32_0", "system"},
		{"terminal_v33_0", base + "terminal_v33_0", "terminal"},
		{"uprofile_v37_0", base + "uprofile_v37_0", "uprofile"},
		{"sharedUprofile_v20_0", base + "sharedUprofile_v20_0", "shareduprofile"},
		// Contract dependency of sharedUprofile_v20_0: the common: prefix in
		// its XSD actually points at the uprofileCommon_v30_0 namespace.
		{"uprofileCommon_v30_0", base + "uprofileCommon_v30_0", "uprofilecommon"},
		{"SessionContext_v1_0", "http://www.travelport.com/soa/common/security/SessionContext_v1_0", "sessioncontext"},
	}
}

func main() {
	w := &world{
		nsToPkg:      map[string]*pkgInfo{},
		complex:      map[typeKey]*complexTypeDef{},
		simple:       map[typeKey]*simpleDef{},
		attrGroups:   map[typeKey]*attrGroupDef{},
		groups:       map[typeKey]*groupDef{},
		elemTypes:    map[typeKey]string{},
		elemSimpleGo: map[typeKey]string{},
		localIndex:   map[string]typeKey{},
		extBases:     map[typeKey]bool{},
		rootElems:    map[typeKey]bool{},
		synth:        map[*Node]typeKey{},
		importedBy:   map[string]map[string]bool{},
		structDeps:   map[string]map[string]bool{},
		typeRename:   map[typeKey]string{},
		attrForm:     map[string]bool{},
	}

	const modRoot = "github.com/shuiyihan12/uapi-go/pkg/generated"
	for _, m := range schemaMods() {
		p := &pkgInfo{
			name:       m.name,
			ns:         m.ns,
			dir:        filepath.Join("pkg/generated", m.name),
			importPath: modRoot + "/" + m.name,
			enumsName:  m.name + "enums",
			enumsDir:   filepath.Join("pkg/generated", m.name, "enums"),
			enumsPath:  modRoot + "/" + m.name + "/enums",
		}
		w.nsToPkg[m.ns] = p
		w.mods = append(w.mods, p)
	}

	var totalComplex, totalEnums int
	for _, m := range schemaMods() {
		files, _ := filepath.Glob(filepath.Join("wsdl", m.dir, "*.xsd"))
		if len(files) == 0 {
			continue
		}
		for _, f := range files {
			root, prefixes, err := parseXSD(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse %s: %v\n", f, err)
				continue
			}
			schema := root
			if schema.Local != "schema" {
				continue
			}
			filePrefixNS := map[string]string{"xs": xsNS}
			for k, v := range prefixes {
				if v != "" {
					filePrefixNS[k] = v // includes the default prefix "" (xmlns=); unprefixed QNames rely on it
				}
			}
			w.collect(schema, filePrefixNS)
		}
	}

	// Pre-pass: resolve the fields of every complexType once, for two reasons:
	// 1) All cross-package references (importedBy) must be registered before
	//    the emit loop decides on "orphan common packages". Otherwise a
	//    package like common32, referenced only by the later-emitted system32,
	//    would be misjudged as an orphan, skipped, and never emitted.
	// 2) All synthetic types from inline anonymous complexTypes must be
	//    registered before emit snapshots the type list. Synthetics can nest
	//    further inline types, so iterate until no new types appear
	//    (fixed-point iteration).
	for {
		keys := make([]typeKey, 0, len(w.complex))
		for k := range w.complex {
			keys = append(keys, k)
		}
		// Sort to make inline-type numbering deterministic. synthInlineType
		// assigns "XInline", "XInline2", ... in the order inline types are
		// first encountered; an unsorted (map-iteration-order) keys slice
		// yields a different numbering on every run, making the generated
		// output non-reproducible.
		sort.Slice(keys, func(i, j int) bool {
			if keys[i].ns != keys[j].ns {
				return keys[i].ns < keys[j].ns
			}
			return keys[i].local < keys[j].local
		})
		before := len(w.complex)
		for _, k := range keys {
			w.resolveFields(w.complex[k], map[string]*pkgInfo{})
		}
		if len(w.complex) == before {
			break
		}
	}

	// Auto-merge mutually recursive namespaces (Go forbids import cycles)
	// and resolve name collisions after merging.
	for _, g := range w.mergeCyclicNamespaces() {
		names := make([]string, 0, len(g))
		for _, p := range g {
			names = append(names, p.name)
		}
		fmt.Printf("merged cyclic namespaces into package %s: %s\n", g[0].name, strings.Join(names, ", "))
	}
	w.resolveGroupCollisions()

	for _, p := range w.mods {
		// Skip orphan common packages that no other package references (e.g.
		// common34/common38): no domain package in this repository uses their
		// XSD version, so emitting them would only produce dead code.
		if strings.HasPrefix(p.name, "common") && len(w.importedBy[p.ns]) == 0 {
			continue
		}
		if p.goPkg() != p {
			// Structs were merged into the representative; clean up any stale
			// file, keeping only the enums subpackage.
			_ = os.Remove(filepath.Join(p.dir, p.name+".go"))
			w.writeEnums(p, &totalEnums)
			continue
		}
		src, ok := w.emitPackage(p)
		if !ok {
			continue
		}
		formatted, err := format.Source([]byte(src))
		if err != nil {
			_ = os.WriteFile("/tmp/bad_"+p.name+".go", []byte(src), 0o644)
			fmt.Fprintf(os.Stderr, "format %s: %v\n", p.name, err)
			os.Exit(1)
		}
		if err := os.MkdirAll(p.dir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", p.dir, err)
			os.Exit(1)
		}
		out := filepath.Join(p.dir, p.name+".go")
		if err := os.WriteFile(out, formatted, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
			os.Exit(1)
		}
		w.writeEnums(p, &totalEnums)
		totalComplex += countComplex(w, p)
	}
	fmt.Printf("generated: %d complexType structs across %d packages, %d enums packages\n", totalComplex, len(w.mods), totalEnums)
}

// writeEnums emits the enums subpackage of a namespace. Enums subpackages
// always exist per namespace (they depend on no other package and can never
// form a cycle), even when the structs were merged into another package.
func (w *world) writeEnums(p *pkgInfo, totalEnums *int) {
	esrc, ok := w.emitEnums(p)
	if !ok {
		return
	}
	efmt, err := format.Source([]byte(esrc))
	if err != nil {
		fmt.Fprintf(os.Stderr, "format enums %s: %v\n", p.enumsName, err)
		os.Exit(1)
	}
	if err := os.MkdirAll(p.enumsDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", p.enumsDir, err)
		os.Exit(1)
	}
	eout := filepath.Join(p.enumsDir, "enums.go")
	if err := os.WriteFile(eout, efmt, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", eout, err)
		os.Exit(1)
	}
	*totalEnums += countEnums(w, p)
}

func countComplex(w *world, p *pkgInfo) int {
	n := 0
	for k := range w.complex {
		if k.ns == p.ns {
			n++
		}
	}
	return n
}

func countEnums(w *world, p *pkgInfo) int {
	n := 0
	for k, s := range w.simple {
		if k.ns == p.ns && s.isEnum {
			n++
		}
	}
	return n
}
