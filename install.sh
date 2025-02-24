#!/bin/bash

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture to Go's naming convention
case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Download the appropriate binary
BINARY_URL="https://github.com/yourusername/kubex/releases/latest/download/kubex-${OS}-${ARCH}"
echo "Downloading kubex for ${OS}-${ARCH}..."
curl -L "$BINARY_URL" -o kubex || { echo "Download failed"; exit 1; }

# Make the binary executable
chmod +x kubex

# Move the binary to /usr/local/bin
echo "Installing kubex..."
sudo mv kubex /usr/local/bin/ || { echo "Installation failed"; exit 1; }

echo "kubex installed successfully!"
