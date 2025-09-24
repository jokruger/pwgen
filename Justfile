set shell := ["bash","-eu","-o","pipefail","-c"]

BINARY_NAME := "pwgen"
MODULE      := "github.com/jokruger/pwgen"
PKG         := "./"
OUT_DIR     := "bin"
DIST_DIR    := "dist"

# Derived version from git (fallback 'dev')
VERSION := `git describe --tags --always --dirty=-dev 2>/dev/null || echo dev`

# Platform targets (space separated)
PLATFORMS := "linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64"

# Centralized linker flags (extend to inject version once you add a variable)
LDFLAGS := "-s -w"

default: build

build:
  mkdir -p {{OUT_DIR}}
  echo "Building {{BINARY_NAME}} (version {{VERSION}}) for host..."
  go build -trimpath -ldflags "{{LDFLAGS}}" -o {{OUT_DIR}}/{{BINARY_NAME}} {{PKG}}
  echo "Built -> {{OUT_DIR}}/{{BINARY_NAME}}"

clean:
  rm -rf {{OUT_DIR}} {{DIST_DIR}}

dist: clean
  mkdir -p {{DIST_DIR}}
  echo "Building distribution artifacts (version {{VERSION}})..."
  for target in {{PLATFORMS}}; do \
    GOOS="${target%/*}"; \
    GOARCH="${target#*/}"; \
    EXT=""; \
    if [ "$GOOS" = "windows" ]; then EXT=".exe"; fi; \
    NAME="{{BINARY_NAME}}-{{VERSION}}-$GOOS-$GOARCH"; \
    OUT_PATH="{{DIST_DIR}}/$NAME"; \
    echo " - $NAME"; \
    mkdir -p "$OUT_PATH"; \
    GOOS="$GOOS" GOARCH="$GOARCH" go build -trimpath -ldflags "{{LDFLAGS}}" -o "$OUT_PATH/{{BINARY_NAME}}$EXT" {{PKG}}; \
    ( cd "{{DIST_DIR}}" && \
      if [ "$GOOS" = "windows" ]; then \
        zip -qr "$NAME.zip" "$NAME"; \
      else \
        tar -C . -czf "$NAME.tar.gz" "$NAME"; \
      fi ); \
    if command -v shasum >/dev/null 2>&1; then \
      shasum -a 256 "{{DIST_DIR}}/$NAME."* | sed 's#{{DIST_DIR}}/##' >> "{{DIST_DIR}}/SHA256SUMS.tmp"; \
    elif command -v sha256sum >/dev/null 2>&1; then \
      sha256sum "{{DIST_DIR}}/$NAME."* | sed 's#{{DIST_DIR}}/##' >> "{{DIST_DIR}}/SHA256SUMS.tmp"; \
    fi; \
  done
  if [ -f "{{DIST_DIR}}/SHA256SUMS.tmp" ]; then \
    sort "{{DIST_DIR}}/SHA256SUMS.tmp" > "{{DIST_DIR}}/SHA256SUMS"; \
    rm "{{DIST_DIR}}/SHA256SUMS.tmp"; \
  fi
  echo "Artifacts:"
  ls -1 {{DIST_DIR}}

install:
  go install {{MODULE}}@latest

# Use +args to preserve arguments more literally (recommended over *args for passthrough)
run +args: build
  ./{{OUT_DIR}}/{{BINARY_NAME}} {{args}}

version:
  @echo {{VERSION}}
