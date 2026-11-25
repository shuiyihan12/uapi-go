package universal

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/shuiyihan12/uapi-go/pkg/client"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

type wsdlOperation struct {
	Name  string `xml:"name,attr"`
	Input struct {
		Message string `xml:"message,attr"`
	} `xml:"input"`
	Output struct {
		Message string `xml:"message,attr"`
	} `xml:"output"`
	SOAP struct {
		Action string `xml:"soapAction,attr"`
	} `xml:"operation"`
}

type wsdlDefinitions struct {
	Messages []struct {
		Name  string `xml:"name,attr"`
		Parts []struct {
			Element string `xml:"element,attr"`
		} `xml:"part"`
	} `xml:"message"`
	PortTypes []struct {
		Name       string          `xml:"name,attr"`
		Operations []wsdlOperation `xml:"operation"`
	} `xml:"portType"`
	Bindings []struct {
		Name       string          `xml:"name,attr"`
		Type       string          `xml:"type,attr"`
		Operations []wsdlOperation `xml:"operation"`
	} `xml:"binding"`
	Services []struct {
		Ports []struct {
			Binding string `xml:"binding,attr"`
			Address struct {
				Location string `xml:"location,attr"`
			} `xml:"address"`
		} `xml:"port"`
	} `xml:"service"`
}

type wsdlRoute struct {
	Operation, Request, Response, Service, Action string
}

func wsdlLocalName(name string) string {
	return name[strings.LastIndexByte(name, ':')+1:]
}

// Resolve message -> portType -> binding -> service/address from the source
// WSDL, independently of the production routing table and package names.
func loadWSDLRoutes(t *testing.T) map[string]wsdlRoute {
	t.Helper()
	read := func(file string) wsdlDefinitions {
		data, err := os.ReadFile("../../../wsdl/universal_v55_0/" + file)
		if err != nil {
			t.Fatal(err)
		}
		var defs wsdlDefinitions
		if err := xml.Unmarshal(data, &defs); err != nil {
			t.Fatal(err)
		}
		return defs
	}
	abstract, concrete := read("UniversalRecordAbstract.wsdl"), read("UniversalRecord.wsdl")
	elements := map[string]string{}
	for _, message := range abstract.Messages {
		if len(message.Parts) == 1 {
			elements[message.Name] = wsdlLocalName(message.Parts[0].Element)
		}
	}
	ports := map[string][]wsdlOperation{}
	for _, port := range abstract.PortTypes {
		ports[port.Name] = port.Operations
	}
	routes := map[string]wsdlRoute{}
	for _, service := range concrete.Services {
		for _, port := range service.Ports {
			for _, binding := range concrete.Bindings {
				if binding.Name != wsdlLocalName(port.Binding) {
					continue
				}
				address, err := url.Parse(port.Address.Location)
				if err != nil {
					t.Fatal(err)
				}
				for _, operation := range ports[wsdlLocalName(binding.Type)] {
					request := elements[wsdlLocalName(operation.Input.Message)]
					response := elements[wsdlLocalName(operation.Output.Message)]
					var action string
					for _, bound := range binding.Operations {
						if bound.Name == operation.Name {
							action = bound.SOAP.Action
						}
					}
					if request == "" || response == "" || action == "" {
						t.Fatalf("incomplete WSDL binding %s", binding.Name)
					}
					routes[request] = wsdlRoute{Operation: strings.TrimSuffix(request, "Req"), Request: request, Response: response, Service: path.Base(address.Path), Action: action}
				}
			}
		}
	}
	if len(routes) != reflect.TypeOf((*UniversalServicePort)(nil)).Elem().NumMethod()-1 {
		t.Fatalf("WSDL ports do not cover the service interface: %d routes", len(routes))
	}
	return routes
}

type routingTransportFunc func(*http.Request) (*http.Response, error)

func (f routingTransportFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestEveryUniversalOperationUsesItsWSDLPort(t *testing.T) {
	routes := loadWSDLRoutes(t)
	for _, region := range []string{"", "apac", "emea", "americas"} {
		t.Run("region="+region, func(t *testing.T) {
			base := "https://configured.example.invalid/custom/uAPI"
			if region != "" {
				base = "https://" + region + ".universal-api.travelport.com/B2BGateway/connect/uAPI"
			}
			var calls atomic.Int64
			transport := routingTransportFunc(func(req *http.Request) (*http.Response, error) {
				calls.Add(1)
				var envelope struct {
					Body struct {
						Contents []byte `xml:",innerxml"`
					} `xml:"Body"`
				}
				if err := xml.NewDecoder(req.Body).Decode(&envelope); err != nil {
					return nil, err
				}
				var root struct {
					XMLName xml.Name
					Trace   string `xml:"TraceId,attr"`
				}
				if err := xml.Unmarshal(envelope.Body.Contents, &root); err != nil {
					return nil, err
				}
				route, ok := routes[root.XMLName.Local]
				if !ok {
					return nil, fmt.Errorf("request missing from WSDL: %s", root.XMLName.Local)
				}
				if got, want := req.URL.String(), base+"/"+route.Service; got != want {
					return nil, fmt.Errorf("%s: URL %s, WSDL requires %s", route.Operation, got, want)
				}
				// B2BGateway requires the SOAPAction header to be present but
				// empty; the WSDL's localhost URI is stale and never sent.
				if values, ok := req.Header["Soapaction"]; !ok || len(values) != 1 || values[0] != "" {
					return nil, fmt.Errorf("%s: SOAPAction must be present and empty, got %q (present=%v)", route.Operation, values, ok)
				}
				if req.Header.Get("Authorization") != "Basic test-only" || req.Header.Get("X-Trace-Id") != "route-"+route.Operation || root.Trace != req.Header.Get("X-Trace-Id") {
					return nil, fmt.Errorf("%s: credentials/trace lost", route.Operation)
				}
				body := fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><%s xmlns="%s"/></soap:Body></soap:Envelope>`, route.Response, root.XMLName.Space)
				return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"text/xml"}}, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
			})
			svc, err := NewUniversalService(client.SOAPConfig{
				BaseEndpoint: "https://configured.example.invalid/custom/uAPI", ServiceName: "UniversalRecordService", Transport: transport,
			}, logging.Noop())
			if err != nil {
				t.Fatal(err)
			}
			defer svc.Close()
			var wg sync.WaitGroup
			errors := make(chan error, len(routes))
			// Mix all operations on the same client concurrently. Endpoint selection
			// must not mutate a shared serviceName/SOAPAction or create cross-call races.
			for _, route := range routes {
				method := reflect.ValueOf(svc).MethodByName(route.Operation)
				if !method.IsValid() {
					t.Fatalf("no SDK method for WSDL request %s", route.Request)
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					ctx := trace.WithTraceID(context.Background(), "route-"+route.Operation)
					ctx = requestctx.WithRegion(requestctx.WithAuthorization(ctx, "Basic test-only"), region)
					request := reflect.New(method.Type().In(1).Elem())
					result := method.Call([]reflect.Value{reflect.ValueOf(ctx), request})
					if !result[1].IsNil() {
						errors <- result[1].Interface().(error)
					}
				}()
			}
			wg.Wait()
			close(errors)
			for err := range errors {
				t.Error(err)
			}
			if calls.Load() != int64(len(routes)) {
				t.Fatalf("wanted one call per operation, got %d", calls.Load())
			}
		})
	}
}

func TestUnknownUniversalOperationDoesNotFallbackToWrongService(t *testing.T) {
	if _, err := operationService("UnknownOperation"); err == nil {
		t.Fatal("unknown operation must not silently route to UniversalRecordService")
	}
}
