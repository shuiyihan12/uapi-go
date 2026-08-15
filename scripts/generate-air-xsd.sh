#!/bin/bash
#
# Generate Go structs for ALL Travelport UAPI XSD schemas.
#
# This uses the purpose-built multi-package generator tools/airxsdgen (not
# xgen). It emits one Go package per XSD schema namespace (e.g. air_v54_0 ->
# package air54, common_v54_0 -> package common54), with enumeration
# simpleTypes placed in a per-domain "enums" sub-package. Cross-package type
# references become imported, package-qualified Go types.
#
# Each package carries:
#   - xml tags with the full namespace URI,
#   - snake_case json tags,
#   - xs:sequence order preserved,
#   - xs:choice members as pointer + omitempty,
#   - minOccurs=0 fields as pointers, maxOccurs>1 as slices,
#   - xs:extension bases embedded (local or cross-package).
#
# Usage: scripts/generate-air-xsd.sh

set -euo pipefail

GREEN='\033[0;32m'; NC='\033[0m'
log_info() { echo -e "\033[0;34m[INFO]\033[0m $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }

cd "$(dirname "$0")/.."

log_info "Generating XSD structs for all Travelport UAPI schemas"
go run ./tools/airxsdgen

log_info "Building all generated packages"
go build ./pkg/generated/...

log_info "Vetting all generated packages"
go vet ./pkg/generated/...

log_success "Generated $(find pkg/generated -name '*.go' | wc -l | tr -d ' ') Go files under pkg/generated"
