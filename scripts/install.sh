#!/bin/sh

# Burnless installer for macOS and Linux

# Usage: curl -fsSL https://raw.githubusercontent.com/Custos-com/Burnless/main/scripts/install.sh | sh

set -e

REPO="Custos-com/Burnless"

BINARY="burnless"

INSTALL_DIR="/usr/local/bin"

# print banner

printf '\033[0;34m'

cat << 'BANNER'

  ██████╗ ██╗   ██╗██████╗ ███╗   ██╗██╗     ███████╗███████╗███████╗

  ██╔══██╗██║   ██║██╔══██╗████╗  ██║██║     ██╔════╝██╔════╝██╔════╝

  ██████╔╝██║   ██║██████╔╝██╔██╗ ██║██║     █████╗  ███████╗███████╗

  ██╔══██╗██║   ██║██╔══██╗██║╚██╗██║██║     ██╔══╝  ╚════██║╚════██║

  ██████╔╝╚██████╔╝██║  ██║██║ ╚████║███████╗███████╗███████║███████║

  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝╚══════╝╚══════╝╚══════╝╚══════╝

BANNER

printf '\033[0m'

printf '  Stop burning your error budget. Stop burning out your team.\n\n'

# detect OS

OS=$(uname -s | tr '[:upper:]' '[:lower:]')

case "$OS" in

  linux)  OS="linux" ;;

  darwin) OS="darwin" ;;

  *)

    echo "Unsupported OS: $OS"

    exit 1

    ;;

esac

# detect architecture

ARCH=$(uname -m)

case "$ARCH" in

  x86_64)        ARCH="amd64" ;;

  aarch64|arm64) ARCH="arm64" ;;

  *)

    echo "Unsupported architecture: $ARCH"

    exit 1

    ;;

esac

# get latest release tag from GitHub API

echo "  Finding latest Burnless release..."

LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | python3 -c "import sys,json; print(json.load(sys.stdin)['tag_name'])")

if [ -z "$LATEST" ]; then

  echo "Error: could not find latest release"

  exit 1

fi

echo "  Installing Burnless ${LATEST} for ${OS}/${ARCH}..."

# build download URL

FILENAME="${BINARY}_${OS}_${ARCH}.tar.gz"

URL="https://github.com/${REPO}/releases/download/${LATEST}/${FILENAME}"

# download and extract

TMP_DIR=$(mktemp -d)

trap "rm -rf $TMP_DIR" EXIT

curl -fsSL "$URL" -o "${TMP_DIR}/${FILENAME}"

tar -xzf "${TMP_DIR}/${FILENAME}" -C "$TMP_DIR"

# install binary

chmod +x "${TMP_DIR}/${BINARY}"

if [ -w "$INSTALL_DIR" ]; then

  mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"

else

  echo "  Installing to ${INSTALL_DIR} (may require sudo password)..."

  sudo mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"

fi

# success

printf '\n\033[0;32m'

echo "  ✓ Burnless ${LATEST} installed successfully!"

printf '\033[0m\n'

echo "  Get started:"

echo "    burnless init        # create your first sre.yaml"

echo "    burnless toil log    # log a toil event"

echo "    burnless --help      # see all commands"

echo ""

echo "  Docs: https://github.com/Custos-com/Burnless"

echo ""

