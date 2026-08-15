#!/bin/bash

# Build script
set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check the Go environment
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed"
        exit 1
    fi
    log_success "Go version: $(go version)"
}

# Install dependencies
install_deps() {
    log_info "Installing dependencies..."
    export GOPATH=$HOME/go-workspace
    go mod tidy
    go mod download
    log_success "Dependencies installed"
}

# Generate WSDL code
generate_wsdl() {
    log_info "Generating WSDL code..."
    if ! ./scripts/generate-air-xsd.sh; then
        log_error "WSDL code generation failed"
        exit 1
    fi
    log_success "WSDL code generated"
}

# Build all commands
build_commands() {
    log_info "Building commands..."

    local commands=(
        "cmd/generate"
        "cmd/generator"
        "cmd/daemon"
    )

    for cmd in "${commands[@]}"; do
        log_info "Building $cmd..."
        if go build -o "bin/$(basename $cmd)" "./$cmd"; then
            log_success "Built $cmd"
        else
            log_error "Failed to build $cmd"
            exit 1
        fi
    done
}

# Remove stale executables from bin/
clean_stale_binaries() {
    local valid_bins=(
        "generate"
        "generator"
        "daemon"
    )

    [ -d "bin" ] || return 0

    log_info "Cleaning stale binaries in bin/..."
    for f in bin/*; do
        [ -e "$f" ] || continue
        local b
        local keep="false"
        b="$(basename "$f")"

        for v in "${valid_bins[@]}"; do
            if [ "$b" = "$v" ]; then
                keep="true"
                break
            fi
        done

        if [ "$keep" = "false" ]; then
            rm -f "$f"
            log_info "Removed stale binary: $b"
        fi
    done
}

# Run tests
run_tests() {
    log_info "Running tests..."
    if go test ./test/...; then
        log_success "All tests passed"
    else
        log_error "Some tests failed"
        exit 1
    fi
}

# Run linting
run_lint() {
    log_info "Running linting..."

    # gofmt check
    if [ "$(gofmt -l . | grep -v vendor | grep -v pkg/generated | wc -l)" -gt 0 ]; then
        log_error "Code is not formatted. Run: gofmt -w ."
        gofmt -l . | grep -v vendor | grep -v pkg/generated
        exit 1
    fi

    # go vet check
    if ! go vet ./...; then
        log_error "go vet failed"
        exit 1
    fi

    log_success "Linting passed"
}

# Create a release package
create_release() {
    local version=${1:-"latest"}
    log_info "Creating release package for version $version..."

    local release_dir="release/uapi-go-$version"
    mkdir -p "$release_dir"

    # Copy binaries
    cp -r bin "$release_dir/"

    # Copy docs and configuration
    cp README.md "$release_dir/"
    cp -r scripts "$release_dir/"
    cp -r wsdl "$release_dir/"

    # Create the archive
    cd release
    tar -czf "uapi-go-$version.tar.gz" "uapi-go-$version"
    log_success "Release package created: release/uapi-go-$version.tar.gz"
    cd ..
}

# Main
main() {
    echo "UAPI Go Build Script"
    echo "==================="
    echo ""

    case "${1:-}" in
        "deps")
            check_go
            install_deps
            ;;
        "wsdl")
            generate_wsdl
            ;;
        "build")
            check_go
            install_deps
            mkdir -p bin
            clean_stale_binaries
            build_commands
            ;;
        "test")
            check_go
            install_deps
            run_tests
            ;;
        "lint")
            check_go
            run_lint
            ;;
        "all"|"")
            check_go
            install_deps
            generate_wsdl
            mkdir -p bin
            clean_stale_binaries
            build_commands
            run_tests
            run_lint
            log_success "Build completed successfully!"
            ;;
        "release")
            main "all"
            create_release "${2:-latest}"
            ;;
        "clean")
            log_info "Cleaning build artifacts..."
            rm -rf bin/
            rm -rf release/
            rm -rf pkg/generated/*
            log_success "Clean completed"
            ;;
        "help"|"-h"|"--help")
            echo "Usage: $0 [command]"
            echo ""
            echo "Commands:"
            echo "  deps       Install dependencies"
            echo "  wsdl       Generate WSDL code"
            echo "  build      Build all commands"
            echo "  test       Run tests"
            echo "  lint       Run linting"
            echo "  all        Full build (default)"
            echo "  release    Create release package"
            echo "  clean      Clean build artifacts"
            echo "  help       Show this help"
            echo ""
            ;;
        *)
            log_error "Unknown command: $1"
            echo "Run '$0 help' for usage information"
            exit 1
            ;;
    esac
}

main "$@"
