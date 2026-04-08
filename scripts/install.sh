#!/usr/bin/env bash
set -euo pipefail

REPO="masudparvezsajjad/website-blocker"
BINARY_NAME="mps-blocker"
VERSION="${VERSION:-latest}"

OS="$(uname -s)"
ARCH="$(uname -m)"

if [[ "$OS" != "Darwin" ]]; then
  echo "This installer currently supports macOS only."
  exit 1
fi

case "$ARCH" in
  arm64|aarch64)
    ASSET_PRIMARY="mps-blocker_darwin_arm64.tar.gz"
    ASSET_LEGACY="adult-blocker_darwin_arm64.tar.gz"
    ;;
  x86_64)
    ASSET_PRIMARY="mps-blocker_darwin_amd64.tar.gz"
    ASSET_LEGACY="adult-blocker_darwin_amd64.tar.gz"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

release_url_for_asset() {
  local asset="$1"
  if [[ "$VERSION" == "latest" ]]; then
    echo "https://github.com/${REPO}/releases/latest/download/${asset}"
  else
    echo "https://github.com/${REPO}/releases/download/${VERSION}/${asset}"
  fi
}

download_asset() {
  local asset="$1"
  curl -fsSL "$(release_url_for_asset "$asset")" -o "$TMP_DIR/${asset}"
}

ASSET=""
echo "Downloading release for ${ARCH}..."
if download_asset "$ASSET_PRIMARY" 2>/dev/null; then
  ASSET="$ASSET_PRIMARY"
elif download_asset "$ASSET_LEGACY" 2>/dev/null; then
  ASSET="$ASSET_LEGACY"
  echo "Note: using legacy archive ${ASSET_LEGACY} (binary installed as ${BINARY_NAME})." >&2
else
  echo "Error: could not download a release for this Mac." >&2
  echo "Tried: ${ASSET_PRIMARY}, ${ASSET_LEGACY}" >&2
  echo "See https://github.com/${REPO}/releases or build from source (README)." >&2
  exit 1
fi

echo "Extracting ${ASSET}..."
tar -xzf "$TMP_DIR/${ASSET}" -C "$TMP_DIR"

BIN_SRC=""
if [[ -f "$TMP_DIR/${BINARY_NAME}" ]]; then
  BIN_SRC="$TMP_DIR/${BINARY_NAME}"
elif [[ -f "$TMP_DIR/blocker" ]]; then
  BIN_SRC="$TMP_DIR/blocker"
else
  echo "Error: archive had no ${BINARY_NAME} or blocker executable." >&2
  ls -la "$TMP_DIR" >&2
  exit 1
fi

INSTALL_DIR="/usr/local/bin"
if [[ -d "/opt/homebrew/bin" && "$ARCH" == "arm64" ]]; then
  INSTALL_DIR="/opt/homebrew/bin"
fi

echo "Installing to ${INSTALL_DIR}..."
sudo install -m 0755 "$BIN_SRC" "${INSTALL_DIR}/${BINARY_NAME}"

echo
echo "Installed successfully."
echo "Run:"
echo "  sudo ${BINARY_NAME} install"
echo "  sudo ${BINARY_NAME} enable"
echo "  sudo ${BINARY_NAME} daemon"
