#!/bin/sh
set -e

REPO="KashifKhn/haft"
BINARY_NAME="haft"

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'
BOLD='\033[1m'

spinner() {
    local pid=$1
    local msg=$2
    local spin='⣾⣽⣻⢿⡿⣟⣯⣷'
    local i=0
    
    while kill -0 "$pid" 2>/dev/null; do
        i=$(( (i + 1) % 8 ))
        printf "\r  ${BLUE}%s${NC} %s" "$(echo "$spin" | cut -c$((i+1)))" "$msg"
        sleep 0.1
    done
    printf "\r  ${GREEN}✓${NC} %s\n" "$msg"
}

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
        printf "  ${RED}✗${NC} Unsupported platform: OS=$(uname -s), Arch=$(uname -m)\n"
        exit 1
    fi

    if [ "$OS" = "windows" ]; then
        ARCHIVE="${BINARY_NAME}-${OS}-${ARCH}.zip"
    else
        ARCHIVE="${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
    
    printf "  ${BLUE}→${NC} Platform: ${BOLD}%s/%s${NC}\n" "$OS" "$ARCH"
    
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${ARCHIVE}" &
    spinner $! "Downloading ${BINARY_NAME} ${VERSION}..."

    (
        cd "$TMP_DIR"
        if [ "$OS" = "windows" ]; then
            unzip -q "$ARCHIVE" 2>/dev/null
        else
            tar -xzf "$ARCHIVE" 2>/dev/null
        fi
    ) &
    spinner $! "Extracting..."

    cd "$TMP_DIR"
    
    if [ "$OS" = "windows" ]; then
        BINARY="${BINARY_NAME}-${OS}-${ARCH}.exe"
    else
        BINARY="${BINARY_NAME}-${OS}-${ARCH}"
    fi

    INSTALL_DIR="/usr/local/bin"
    
    if [ ! -w "$INSTALL_DIR" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    fi

    (
        if [ "$OS" = "windows" ]; then
            mv "$BINARY" "${INSTALL_DIR}/${BINARY_NAME}.exe"
        else
            mv "$BINARY" "${INSTALL_DIR}/${BINARY_NAME}"
            chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
        fi
    ) &
    spinner $! "Installing to ${INSTALL_DIR}..."

    echo ""
    printf "  ${GREEN}${BOLD}✓ Successfully installed %s %s${NC}\n" "$BINARY_NAME" "$VERSION"
    echo ""
    
    if [ "$INSTALL_DIR" = "$HOME/.local/bin" ]; then
        case "$PATH" in
            *"$HOME/.local/bin"*) ;;
            *)
                printf "  ${YELLOW}!${NC} Add to your PATH:\n"
                case "$SHELL" in
                    */bash)
                        printf "    ${BLUE}echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc && source ~/.bashrc${NC}\n"
                        ;;
                    */zsh)
                        printf "    ${BLUE}echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc && source ~/.zshrc${NC}\n"
                        ;;
                    *)
                        printf "    Add ${BLUE}%s${NC} to your PATH\n" "$INSTALL_DIR"
                        ;;
                esac
                echo ""
                ;;
        esac
    fi

    printf "  ${GREEN}→${NC} Run ${BOLD}haft --help${NC} to get started\n"
    echo ""
}

main() {
    echo ""
    printf "${BLUE}${BOLD}"
    echo "  _   _    _    _____ _____ "
    echo " | | | |  / \  |  ___|_   _|"
    echo " | |_| | / _ \ | |_    | |  "
    echo " |  _  |/ ___ \|  _|   | |  "
    echo " |_| |_/_/   \_\_|     |_|  "
    printf "${NC}"
    echo ""
    printf "  ${BOLD}The Spring Boot CLI${NC}\n"
    echo ""

    VERSION=${1:-}
    
    if [ -z "$VERSION" ]; then
        printf "  ${BLUE}⣾${NC} Fetching latest version...\r"
        VERSION=$(get_latest_version)
        printf "  ${GREEN}✓${NC} Latest version: ${BOLD}%s${NC}        \n" "$VERSION"
    fi
    
    if [ -z "$VERSION" ]; then
        printf "  ${RED}✗${NC} Could not determine latest version\n"
        exit 1
    fi

    OS=$(get_os)
    ARCH=$(get_arch)

    echo ""
    download_and_install "$VERSION" "$OS" "$ARCH"
}

main "$@"
