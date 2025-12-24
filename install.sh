#!/bin/sh
set -e

REPO="KashifKhn/haft"
BINARY_NAME="haft"

get_latest_version() {
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | 
        grep '"tag_name":' | 
        sed -E 's/.*"([^"]+)".*/\1/'
}

get_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)          echo "unknown" ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64" ;;
        arm64|aarch64)  echo "arm64" ;;
        *)              echo "unknown" ;;
    esac
}

download_and_install() {
    VERSION=$1
    OS=$2
    ARCH=$3

    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        echo "Error: Unsupported platform: OS=$(uname -s), Arch=$(uname -m)"
        exit 1
    fi

    if [ "$OS" = "windows" ]; then
        ARCHIVE="${BINARY_NAME}-${OS}-${ARCH}.zip"
    else
        ARCHIVE="${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
    
    echo "Downloading ${BINARY_NAME} ${VERSION} for ${OS}/${ARCH}..."
    
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    if ! curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${ARCHIVE}"; then
        echo "Error: Failed to download ${DOWNLOAD_URL}"
        exit 1
    fi

    echo "Extracting..."
    cd "$TMP_DIR"
    
    if [ "$OS" = "windows" ]; then
        unzip -q "$ARCHIVE"
        BINARY="${BINARY_NAME}-${OS}-${ARCH}.exe"
    else
        tar -xzf "$ARCHIVE"
        BINARY="${BINARY_NAME}-${OS}-${ARCH}"
    fi

    INSTALL_DIR="/usr/local/bin"
    
    if [ ! -w "$INSTALL_DIR" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
        echo "Installing to ${INSTALL_DIR} (no sudo access)"
    fi

    echo "Installing to ${INSTALL_DIR}..."
    
    if [ "$OS" = "windows" ]; then
        mv "$BINARY" "${INSTALL_DIR}/${BINARY_NAME}.exe"
    else
        mv "$BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    echo ""
    echo "Successfully installed ${BINARY_NAME} ${VERSION}!"
    echo ""
    
    if [ "$INSTALL_DIR" = "$HOME/.local/bin" ]; then
        case "$SHELL" in
            */bash)
                echo "Add to your PATH by running:"
                echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc"
                echo "  source ~/.bashrc"
                ;;
            */zsh)
                echo "Add to your PATH by running:"
                echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc"
                echo "  source ~/.zshrc"
                ;;
            *)
                echo "Add ${INSTALL_DIR} to your PATH"
                ;;
        esac
        echo ""
    fi

    echo "Run 'haft --help' to get started"
}

main() {
    echo ""
    echo "  _   _    _    _____ _____ "
    echo " | | | |  / \  |  ___|_   _|"
    echo " | |_| | / _ \ | |_    | |  "
    echo " |  _  |/ ___ \|  _|   | |  "
    echo " |_| |_/_/   \_\_|     |_|  "
    echo ""
    echo "The Spring Boot CLI"
    echo ""

    VERSION=${1:-$(get_latest_version)}
    
    if [ -z "$VERSION" ]; then
        echo "Error: Could not determine latest version"
        exit 1
    fi

    OS=$(get_os)
    ARCH=$(get_arch)

    download_and_install "$VERSION" "$OS" "$ARCH"
}

main "$@"
