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
    local delay=0.1
    local spinstr='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
    while kill -0 "$pid" 2>/dev/null; do
        for i in $(seq 0 9); do
            printf "\r  ${BLUE}%s${NC} %s" "$(echo "$spinstr" | cut -c$((i+1)))" "$2"
            sleep $delay
        done
    done
    printf "\r  ${GREEN}✓${NC} %s\n" "$2"
}

progress_bar() {
    local current=$1
    local total=$2
    local width=40
    local percentage=$((current * 100 / total))
    local filled=$((current * width / total))
    local empty=$((width - filled))
    
    printf "\r  ["
    printf "%${filled}s" | tr ' ' '='
    if [ $filled -lt $width ]; then
        printf ">"
        printf "%$((empty - 1))s" | tr ' ' ' '
    fi
    printf "] %3d%%" "$percentage"
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

download_with_progress() {
    local url=$1
    local output=$2
    
    if command -v curl >/dev/null 2>&1; then
        curl -fSL --progress-bar "$url" -o "$output" 2>&1
    elif command -v wget >/dev/null 2>&1; then
        wget --progress=bar:force "$url" -O "$output" 2>&1
    else
        echo "${RED}Error: curl or wget required${NC}"
        exit 1
    fi
}

download_and_install() {
    VERSION=$1
    OS=$2
    ARCH=$3

    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        echo "${RED}Error: Unsupported platform: OS=$(uname -s), Arch=$(uname -m)${NC}"
        exit 1
    fi

    if [ "$OS" = "windows" ]; then
        ARCHIVE="${BINARY_NAME}-${OS}-${ARCH}.zip"
    else
        ARCHIVE="${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
    
    printf "  ${BLUE}→${NC} Downloading ${BOLD}%s %s${NC} for %s/%s\n" "$BINARY_NAME" "$VERSION" "$OS" "$ARCH"
    
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    (download_with_progress "$DOWNLOAD_URL" "${TMP_DIR}/${ARCHIVE}") &
    DOWNLOAD_PID=$!
    
    spinner $DOWNLOAD_PID "Downloading..."
    
    wait $DOWNLOAD_PID || {
        printf "\r  ${RED}✗${NC} Download failed\n"
        echo "${RED}Error: Failed to download ${DOWNLOAD_URL}${NC}"
        exit 1
    }

    (
        cd "$TMP_DIR"
        if [ "$OS" = "windows" ]; then
            unzip -q "$ARCHIVE"
        else
            tar -xzf "$ARCHIVE"
        fi
    ) &
    EXTRACT_PID=$!
    
    spinner $EXTRACT_PID "Extracting..."
    
    wait $EXTRACT_PID || {
        printf "\r  ${RED}✗${NC} Extraction failed\n"
        exit 1
    }

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
    INSTALL_PID=$!
    
    spinner $INSTALL_PID "Installing to ${INSTALL_DIR}..."
    
    wait $INSTALL_PID || {
        printf "\r  ${RED}✗${NC} Installation failed\n"
        exit 1
    }

    echo ""
    printf "  ${GREEN}${BOLD}Successfully installed %s %s!${NC}\n" "$BINARY_NAME" "$VERSION"
    echo ""
    
    if [ "$INSTALL_DIR" = "$HOME/.local/bin" ]; then
        if ! echo "$PATH" | grep -q "$HOME/.local/bin"; then
            printf "  ${YELLOW}Note:${NC} Add to your PATH:\n"
            case "$SHELL" in
                */bash)
                    printf "    ${BLUE}echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc${NC}\n"
                    printf "    ${BLUE}source ~/.bashrc${NC}\n"
                    ;;
                */zsh)
                    printf "    ${BLUE}echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc${NC}\n"
                    printf "    ${BLUE}source ~/.zshrc${NC}\n"
                    ;;
                *)
                    printf "    Add ${BLUE}%s${NC} to your PATH\n" "$INSTALL_DIR"
                    ;;
            esac
            echo ""
        fi
    fi

    printf "  ${GREEN}→${NC} Run ${BOLD}'haft --help'${NC} to get started\n"
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
        printf "  ${BLUE}→${NC} Fetching latest version...\r"
        VERSION=$(get_latest_version)
        printf "  ${GREEN}✓${NC} Latest version: ${BOLD}%s${NC}\n" "$VERSION"
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
