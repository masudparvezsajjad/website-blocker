#!/usr/bin/env bash
set -euo pipefail

REPO="masudparvzsajjad/website-blocker"
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
    ASSET="mps-blocker_darwin_arm64.tar.gz"
    ;;
  x86_64)
    ASSET="mps-blocker_darwin_amd64.tar.gz"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

if [[ "$VERSION" == "latest" ]]; then
  DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
else
  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
fi

echo "Downloading ${ASSET}..."
curl -fsSL "$DOWNLOAD_URL" -o "$TMP_DIR/${ASSET}"

echo "Extracting..."
tar -xzf "$TMP_DIR/${ASSET}" -C "$TMP_DIR"

INSTALL_DIR="/usr/local/bin"
if [[ -d "/opt/homebrew/bin" && "$ARCH" == "arm64" ]]; then
  INSTALL_DIR="/opt/homebrew/bin"
fi

echo "Installing to ${INSTALL_DIR}..."
sudo install -m 0755 "$TMP_DIR/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"

echo
echo "Installed successfully."
echo "Run:"
echo "  sudo ${BINARY_NAME} install"
echo "  sudo ${BINARY_NAME} enable"
echo "  sudo ${BINARY_NAME} daemon"