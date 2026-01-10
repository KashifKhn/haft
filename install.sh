#!/bin/sh
set -e

REPO="KashifKhn/haft"
BINARY_NAME="haft"
DOCS_URL="https://haft.kashifkhan.dev"

CYAN='\033[0;36m'
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
MUTED='\033[0;2m'
BOLD='\033[1m'
NC='\033[0m'

requested_version=""
no_modify_path=false

usage() {
    cat <<EOF
Haft Installer - The missing Spring Boot CLI

Usage: install.sh [options]

Options:
    -h, --help              Display this help message
    -v, --version <version> Install a specific version (e.g., 0.5.0)
        --no-modify-path    Don't modify shell config files

Examples:
    curl -fsSL https://raw.githubusercontent.com/$REPO/main/install.sh | sh
    curl -fsSL ... | sh -s -- --version 0.5.0
    curl -fsSL ... | sh -s -- --no-modify-path
EOF
}

print_logo() {
    printf "\n"
    printf "${CYAN}   ▄▄${NC}   ${BOLD}██  ██   ████   ████████  ████████${NC}   ${CYAN}█${NC}\n"
    printf "${CYAN}  ▀▀▀▄${NC}  ${BOLD}██  ██  ██  ██  ██           ██${NC}      ${CYAN}█${NC}\n"
    printf "${CYAN}     █${NC}  ${BOLD}██████  ██████  █████        ██${NC}      ${CYAN}█${NC}\n"
    printf "${CYAN}  ▄▄▄▀${NC}  ${BOLD}██  ██  ██  ██  ██           ██${NC}      ${CYAN}█${NC}\n"
    printf "${CYAN}   ▀▀${NC}   ${BOLD}██  ██  ██  ██  ██           ██${NC}      ${CYAN}█${NC}\n"
    printf "\n"
    printf "              ${MUTED}The missing Spring Boot CLI${NC}\n"
    printf "            ${MUTED}scaffolding at terminal speed${NC}\n"
    printf "\n"
}

print_success() {
    printf "  ${GREEN}✓${NC} $1\n"
}

print_error() {
    printf "  ${RED}✗${NC} $1\n"
}

print_info() {
    printf "  ${CYAN}→${NC} $1\n"
}

print_warning() {
    printf "  ${YELLOW}!${NC} $1\n"
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

detect_rosetta() {
    if [ "$(get_os)" = "darwin" ] && [ "$(get_arch)" = "amd64" ]; then
        rosetta_flag=$(sysctl -n sysctl.proc_translated 2>/dev/null || echo 0)
        if [ "$rosetta_flag" = "1" ]; then
            echo "arm64"
            return
        fi
    fi
    get_arch
}

get_latest_version() {
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | 
        grep '"tag_name":' | 
        sed -E 's/.*"v?([^"]+)".*/\1/'
}

check_existing_version() {
    haft_bin=""
    
    if command -v haft >/dev/null 2>&1; then
        haft_bin=$(command -v haft)
    elif [ -x "/usr/local/bin/haft" ]; then
        haft_bin="/usr/local/bin/haft"
    elif [ -x "$HOME/.local/bin/haft" ]; then
        haft_bin="$HOME/.local/bin/haft"
    fi
    
    if [ -n "$haft_bin" ] && [ -x "$haft_bin" ]; then
        tmp_version_file=$(mktemp)
        "$haft_bin" version > "$tmp_version_file" 2>&1 || true
        if [ -s "$tmp_version_file" ]; then
            installed_version=$(sed -E 's/.*version v?([0-9.]+).*/\1/' "$tmp_version_file")
            rm -f "$tmp_version_file"
            if [ -n "$installed_version" ]; then
                echo "$installed_version"
                return
            fi
        fi
        rm -f "$tmp_version_file"
    fi
    echo ""
}

print_progress() {
    bytes="$1"
    total="$2"
    
    if [ "$total" -le 0 ] 2>/dev/null; then
        return 0
    fi
    
    width=50
    percent=$((bytes * 100 / total))
    if [ "$percent" -gt 100 ]; then
        percent=100
    fi
    
    filled=$((percent * width / 100))
    empty=$((width - filled))
    
    bar=""
    i=0
    while [ $i -lt $filled ]; do
        bar="${bar}■"
        i=$((i + 1))
    done
    i=0
    while [ $i -lt $empty ]; do
        bar="${bar}･"
        i=$((i + 1))
    done
    
    printf "\r  ${CYAN}%s${NC} ${MUTED}%3d%%${NC}" "$bar" "$percent"
}

download_with_progress() {
    url="$1"
    output="$2"
    
    if [ ! -t 2 ]; then
        curl -fsSL "$url" -o "$output"
        return $?
    fi
    
    printf "\033[?25l"
    
    trap 'printf "\033[?25h"' EXIT INT TERM
    
    total_size=$(curl -sI "$url" | grep -i "content-length" | tail -1 | tr -d '\r' | awk '{print $2}')
    
    if [ -z "$total_size" ] || [ "$total_size" = "0" ]; then
        curl -fsSL "$url" -o "$output"
        ret=$?
        printf "\033[?25h"
        echo ""
        return $ret
    fi
    
    curl -fsSL "$url" -o "$output" 2>/dev/null &
    curl_pid=$!
    
    while kill -0 "$curl_pid" 2>/dev/null; do
        if [ -f "$output" ]; then
            current_size=$(wc -c < "$output" 2>/dev/null | tr -d ' ' || echo "0")
            print_progress "$current_size" "$total_size"
        fi
        sleep 0.1
    done
    
    wait $curl_pid
    ret=$?
    
    if [ $ret -eq 0 ]; then
        print_progress "$total_size" "$total_size"
    fi
    
    printf "\033[?25h"
    echo ""
    
    return $ret
}

add_to_path() {
    config_file="$1"
    path_command="$2"
    
    if [ ! -f "$config_file" ]; then
        return 1
    fi
    
    if grep -q "$INSTALL_DIR" "$config_file" 2>/dev/null; then
        return 0
    fi
    
    if [ -w "$config_file" ]; then
        echo "" >> "$config_file"
        echo "# haft" >> "$config_file"
        echo "$path_command" >> "$config_file"
        print_success "Added to PATH in $config_file"
        return 0
    fi
    
    return 1
}

configure_path() {
    if [ "$no_modify_path" = true ]; then
        return
    fi
    
    case "$PATH" in
        *"$INSTALL_DIR"*) 
            return
            ;;
    esac
    
    current_shell=$(basename "$SHELL" 2>/dev/null || echo "sh")
    
    path_added=false
    
    case "$current_shell" in
        fish)
            if add_to_path "$HOME/.config/fish/config.fish" "fish_add_path $INSTALL_DIR"; then
                path_added=true
            fi
            ;;
        zsh)
            for config_file in "$HOME/.zshrc" "$HOME/.zshenv"; do
                if [ -f "$config_file" ]; then
                    if add_to_path "$config_file" "export PATH=\"$INSTALL_DIR:\$PATH\""; then
                        path_added=true
                        break
                    fi
                fi
            done
            ;;
        bash)
            for config_file in "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.profile"; do
                if [ -f "$config_file" ]; then
                    if add_to_path "$config_file" "export PATH=\"$INSTALL_DIR:\$PATH\""; then
                        path_added=true
                        break
                    fi
                fi
            done
            ;;
        *)
            if [ -f "$HOME/.profile" ]; then
                if add_to_path "$HOME/.profile" "export PATH=\"$INSTALL_DIR:\$PATH\""; then
                    path_added=true
                fi
            fi
            ;;
    esac
    
    if [ "$path_added" = false ]; then
        print_warning "Add to your PATH manually:"
        printf "    ${MUTED}export PATH=\"%s:\$PATH\"${NC}\n" "$INSTALL_DIR"
    fi
}

download_and_install() {
    VERSION="$1"
    OS="$2"
    ARCH="$3"

    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        print_error "Unsupported platform: OS=$(uname -s), Arch=$(uname -m)"
        exit 1
    fi

    if [ "$OS" = "windows" ]; then
        ARCHIVE="${BINARY_NAME}-${OS}-${ARCH}.zip"
    else
        ARCHIVE="${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARCHIVE}"
    
    print_info "Platform: ${BOLD}${OS}/${ARCH}${NC}"
    echo ""
    
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    if ! download_with_progress "$DOWNLOAD_URL" "${TMP_DIR}/${ARCHIVE}"; then
        print_error "Failed to download ${BINARY_NAME} v${VERSION}"
        print_info "Check releases: https://github.com/${REPO}/releases"
        exit 1
    fi

    cd "$TMP_DIR"
    if [ "$OS" = "windows" ]; then
        unzip -q "$ARCHIVE" 2>/dev/null || {
            print_error "Failed to extract archive"
            exit 1
        }
    else
        tar -xzf "$ARCHIVE" 2>/dev/null || {
            print_error "Failed to extract archive"
            exit 1
        }
    fi

    if [ "$OS" = "windows" ]; then
        BINARY="${BINARY_NAME}-${OS}-${ARCH}.exe"
        TARGET_BINARY="${BINARY_NAME}.exe"
    else
        BINARY="${BINARY_NAME}-${OS}-${ARCH}"
        TARGET_BINARY="${BINARY_NAME}"
    fi

    INSTALL_DIR="/usr/local/bin"
    
    if [ ! -w "$INSTALL_DIR" ]; then
        INSTALL_DIR="$HOME/.local/bin"
        mkdir -p "$INSTALL_DIR"
    fi
    
    export INSTALL_DIR

    if [ "$OS" = "windows" ]; then
        mv "$BINARY" "${INSTALL_DIR}/${TARGET_BINARY}"
    else
        mv "$BINARY" "${INSTALL_DIR}/${TARGET_BINARY}"
        chmod +x "${INSTALL_DIR}/${TARGET_BINARY}"
    fi

    print_success "Installed to ${INSTALL_DIR}/${TARGET_BINARY}"
}

while [ $# -gt 0 ]; do
    case "$1" in
        -h|--help)
            usage
            exit 0
            ;;
        -v|--version)
            if [ -n "${2:-}" ]; then
                requested_version="$2"
                requested_version="${requested_version#v}"
                shift 2
            else
                print_error "--version requires a version argument"
                exit 1
            fi
            ;;
        --no-modify-path)
            no_modify_path=true
            shift
            ;;
        *)
            print_warning "Unknown option '$1'"
            shift
            ;;
    esac
done

main() {
    print_logo

    if [ -n "$requested_version" ]; then
        VERSION="$requested_version"
        print_info "Installing version: ${BOLD}v${VERSION}${NC}"
    else
        printf "  ${MUTED}Fetching latest version...${NC}\r"
        VERSION=$(get_latest_version)
        if [ -z "$VERSION" ]; then
            echo ""
            print_error "Could not determine latest version"
            exit 1
        fi
        printf "                                    \r"
        print_success "Latest version: ${BOLD}v${VERSION}${NC}"
    fi

    existing_version=$(check_existing_version)
    if [ -n "$existing_version" ]; then
        clean_existing=$(echo "$existing_version" | sed 's/^v//')
        clean_version=$(echo "$VERSION" | sed 's/^v//')
        if [ "$clean_existing" = "$clean_version" ]; then
            print_info "Version v${VERSION} already installed"
            exit 0
        else
            print_info "Upgrading from v${clean_existing} to v${clean_version}"
        fi
    fi

    OS=$(get_os)
    ARCH=$(detect_rosetta)

    echo ""
    download_and_install "$VERSION" "$OS" "$ARCH"

    configure_path

    if [ -n "${GITHUB_ACTIONS-}" ] && [ "${GITHUB_ACTIONS}" = "true" ]; then
        echo "$INSTALL_DIR" >> "$GITHUB_PATH"
        print_success "Added to \$GITHUB_PATH"
    fi

    echo ""
    print_info "Run ${BOLD}haft --help${NC} to get started"
    print_info "Docs: ${MUTED}${DOCS_URL}${NC}"
    echo ""
}

main
