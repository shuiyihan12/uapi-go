# syntax=docker/dockerfile:1

# RUNTIME_IMAGE must be declared before the first FROM (global scope) so the
# FROM lines can reference it. Defaults to the distroless static image; on
# networks that cannot reach gcr.io, point it at a trusted mirror via build
# arg, or pull it from a reachable network and retag it locally (BuildKit
# does not force remote resolution for locally present images):
#   docker build --build-arg RUNTIME_IMAGE=<mirror>/distroless/static-debian12:nonroot .
#   (docker compose passes it via the UAPI_RUNTIME_IMAGE env var; see docker-compose.yml)
# Note: the runtime image is the container's trust anchor; use trusted sources only.
ARG RUNTIME_IMAGE=gcr.io/distroless/static-debian12:nonroot

# ============================================================================
# Build stage: compile static Go binaries (daemon + healthcheck helper)
# ============================================================================
FROM golang:1.25-alpine AS builder

# TARGETOS / TARGETARCH are injected automatically by buildx --platform for
# multi-arch images (linux/amd64, linux/arm64); CGO_ENABLED=0 ensures a
# fully static link.
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

ENV CGO_ENABLED=0 \
    GOOS=${TARGETOS:-linux} \
    GOARCH=${TARGETARCH}

WORKDIR /src

# Copy the module manifests first: with go.sum unchanged the cache layer
# hits and the dependency download is skipped.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Then copy the sources and build. trimpath strips local paths, -s -w strips
# the symbol table: a smaller, safer image.
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" \
    -o /out/uapi-daemon ./cmd/daemon && \
    go build -trimpath -ldflags "-s -w" -o /out/healthcheck ./cmd/healthcheck

# ============================================================================
# Runtime stage: distroless static image (no shell / package manager, ships
# CA certificates and tzdata, runs as non-root user 65532). The standard Go
# production setup with minimal attack surface and image size.
# ============================================================================
FROM ${RUNTIME_IMAGE}

COPY --from=builder /out/uapi-daemon /uapi-daemon
COPY --from=builder /out/healthcheck /healthcheck

# The daemon serves both the business API and the ops endpoints
# (/health /ready /stats /metrics).
EXPOSE 8080

# The container-level health check probes the process self-check (/ready),
# not upstream reachability (so upstream flaps cannot get the container
# killed). The server port can be overridden via PORT; healthcheck
# reads the same variable.
ENV PORT=8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/healthcheck"]

USER nonroot:nonroot

ENTRYPOINT ["/uapi-daemon"]
