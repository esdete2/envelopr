#!/bin/sh
set -e

# Configuration
BINARY_NAME="envelopr"
GITHUB_REPO="esdete2/envelopr"
INSTALL_DIR="/usr/local/bin"

# Detect OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

if [ "$OS" = "Darwin" ]; then
    if [ "$ARCH" = "arm64" ]; then
        PACKAGE="envelopr_Darwin_arm64.tar.gz"
    else
        PACKAGE="envelopr_Darwin_x86_64.tar.gz"
    fi
elif [ "$OS" = "Linux" ]; then
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
        PACKAGE="envelopr_Linux_arm64.tar.gz"
    else
        PACKAGE="envelopr_Linux_x86_64.tar.gz"
    fi
else
    echo "Unsupported operating system: $OS"
    exit 1
fi

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf $TMP_DIR' EXIT

# Download latest release
echo "Downloading $PACKAGE..."
RELEASE_URL="https://github.com/$GITHUB_REPO/releases/latest/download/$PACKAGE"
curl -sL "$RELEASE_URL" | tar xz -C "$TMP_DIR"

# Install binary
echo "Installing to $INSTALL_DIR..."
if [ ! -w "$INSTALL_DIR" ]; then
    sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

echo "Successfully installed $BINARY_NAME"
echo "Run 'envelopr --help' to get started"