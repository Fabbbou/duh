#!/bin/bash
set -e

# Duh installation script
# Usage: curl -sSL https://raw.githubusercontent.com/Fabbbou/duh/main/install.sh | sh

GITHUB_REPO="Fabbbou/duh"  # Replace with actual username
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="duh"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    printf "${BLUE}[INFO]${NC} %s\n" "$1" >&2
}

log_success() {
    printf "${GREEN}[SUCCESS]${NC} %s\n" "$1" >&2
}

log_warning() {
    printf "${YELLOW}[WARNING]${NC} %s\n" "$1" >&2
}

log_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1" >&2
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)
    
    # Normalize OS
    case $os in
        linux*)
            os="linux"
            ;;
        darwin*)
            os="darwin"
            ;;
        mingw*|msys*|cygwin*)
            os="windows"
            ;;
        *)
            log_error "Unsupported OS: $os"
            exit 1
            ;;
    esac
    
    # Normalize architecture
    case $arch in
        x86_64|amd64)
            arch="amd64"
            ;;
        arm64|aarch64)
            arch="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest release version from GitHub
get_latest_version() {
    log_info "Fetching latest release version..."
    
    local latest_url="https://api.github.com/repos/${GITHUB_REPO}/releases/latest"
    local version
    
    if command -v curl >/dev/null 2>&1; then
        version=$(curl -s "$latest_url" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        version=$(wget -qO- "$latest_url" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
    else
        log_error "Neither curl nor wget is available"
        exit 1
    fi
    
    if [ -z "$version" ]; then
        log_error "Failed to get latest version"
        exit 1
    fi
    
    echo "$version"
}

# Download file with progress
download_file() {
    local url="$1"
    local output="$2"
    
    # Ensure output directory exists
    local output_dir=$(dirname "$output")
    if [ ! -d "$output_dir" ]; then
        log_error "Output directory does not exist: $output_dir"
        return 1
    fi
    
    log_info "Downloading from: $url"
    log_info "Saving to: $output"
    
    if command -v curl >/dev/null 2>&1; then
        echo "Using curl for download..."
        echo "url : $url output: $output"
        if ! curl -L --fail --progress-bar "$url" -o "$output"; then
            log_error "Download failed with curl"
            return 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget --progress=bar:force -O "$output" "$url"; then
            log_error "Download failed with wget"
            return 1
        fi
    else
        log_error "Neither curl nor wget is available"
        return 1
    fi
    
    # Verify file was downloaded
    if [ ! -f "$output" ]; then
        log_error "Downloaded file does not exist: $output"
        return 1
    fi
    
    log_success "Downloaded successfully: $(ls -lh "$output" | awk '{print $5}')"
}

# Verify checksum
verify_checksum() {
    local file="$1"
    local checksum_file="$2"
    local filename=$(basename "$file")
    
    log_info "Verifying checksum..."
    
    if ! command -v sha256sum >/dev/null 2>&1; then
        log_warning "sha256sum not available, skipping checksum verification"
        return 0
    fi
    
    # Extract checksum for our file
    local expected_checksum=$(grep "$filename" "$checksum_file" 2>/dev/null | cut -d' ' -f1)
    
    if [ -z "$expected_checksum" ]; then
        log_warning "Checksum not found for $filename, skipping verification"
        return 0
    fi
    
    # Calculate actual checksum
    local actual_checksum=$(sha256sum "$file" | cut -d' ' -f1)
    
    if [ "$expected_checksum" = "$actual_checksum" ]; then
        log_success "Checksum verification passed"
        return 0
    else
        log_error "Checksum verification failed!"
        log_error "Expected: $expected_checksum"
        log_error "Actual:   $actual_checksum"
        exit 1
    fi
}

# Check if running as root for install location
check_install_permissions() {
    # Skip permission check on Windows - we'll handle it during install
    if [[ $(detect_platform) == *"windows"* ]]; then
        return 0
    fi
    
    if [ ! -w "$INSTALL_DIR" ]; then
        if [ "$EUID" -ne 0 ]; then
            log_error "Installation requires write access to $INSTALL_DIR"
            log_error "Please run with sudo or choose a different install location"
            log_error "Example: INSTALL_DIR=\$HOME/.local/bin $0"
            exit 1
        fi
    fi
}

# Main installation function
install_duh() {
    log_info "Starting Duh installation..."
    
    # Detect platform
    local platform=$(detect_platform)
    log_info "Detected platform: $platform"
    
    # Get latest version
    local version=$(get_latest_version)
    log_info "Latest version: $version"
    
    # Determine binary name
    local binary_suffix=""
    if [[ $platform == *"windows"* ]]; then
        binary_suffix=".exe"
        # On Windows, provide sensible defaults only if INSTALL_DIR is not customized
        if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
            # User didn't customize INSTALL_DIR, use Windows-appropriate default
            if [ -n "$PROGRAMDATA" ]; then
                INSTALL_DIR="$PROGRAMDATA/duh"
            elif [ -n "$USERPROFILE" ]; then
                INSTALL_DIR="$USERPROFILE/.local/bin"
            else
                INSTALL_DIR="$HOME/.local/bin"
            fi
            log_warning "Windows detected: Installing to $INSTALL_DIR"
            log_warning "You may need to add $INSTALL_DIR to your PATH manually"
        else
            log_info "Using custom install directory: $INSTALL_DIR"
        fi
    fi
    
    local binary_name="duh-${platform}${binary_suffix}"
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${binary_name}"
    local checksum_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/checksums.txt"
    
    # Create temporary directory
    local tmp_dir
    if [[ $platform == *"windows"* ]]; then
        # Use current directory for Windows - /tmp has issues with executables in Git Bash
        tmp_dir="./duh-install-$$"
    else
        tmp_dir=$(mktemp -d)
    fi
    
    # Ensure temp directory exists and is writable
    if ! mkdir -p "$tmp_dir"; then
        log_error "Failed to create temporary directory: $tmp_dir"
        exit 1
    fi
    
    if [ ! -w "$tmp_dir" ]; then
        log_error "Temporary directory is not writable: $tmp_dir"
        exit 1
    fi
    
    log_info "Using temporary directory: $tmp_dir"
    local binary_path="$tmp_dir/$binary_name"
    local checksum_path="$tmp_dir/checksums.txt"
    
    # Cleanup function
    cleanup() {
        if [ -n "$tmp_dir" ] && [ -d "$tmp_dir" ]; then
            log_info "Cleaning up temporary directory: $tmp_dir"
            rm -rf "$tmp_dir"
        fi
    }
    trap cleanup EXIT
    
    log_info "Downloading $binary_name..."
    download_file "$download_url" "$binary_path"
    
    log_info "Downloading checksums..."
    download_file "$checksum_url" "$checksum_path"
    
    # Verify checksum
    verify_checksum "$binary_path" "$checksum_path"
    
    # Check install permissions
    check_install_permissions
    
    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"
    
    # Install binary
    local install_path="$INSTALL_DIR/$BINARY_NAME$binary_suffix"
    log_info "Installing to $install_path..."
    
    cp "$binary_path" "$install_path"
    chmod +x "$install_path"
    
    log_success "Duh $version installed successfully!"
    log_info "Location: $install_path"
    
    # Manual cleanup for better user feedback
    cleanup
    log_info "Location: $install_path"
    
    # Verify installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        log_success "Duh is now available in your PATH"
        log_info "Run 'duh --help' to get started"
    else
        log_warning "Duh installed but not found in PATH"
        log_warning "You may need to add $INSTALL_DIR to your PATH"
        log_warning "Or restart your terminal"
    fi
}

# Run installation
install_duh