#!/bin/bash
set -e

# Duh installation script
# Usage: curl -sSL https://raw.githubusercontent.com/Fabbbou/duh/main/install.sh | sh

GITHUB_REPO="Fabbbou/duh"  # Replace with actual username
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
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

# Check if duh is currently installed
is_duh_installed() {
    local install_path="$1"
    [ -f "$install_path" ]
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
    local install_dir="$1"
    local platform=$(detect_platform)
    
    # Skip permission check on Windows - we'll handle it during install
    case "$platform" in
        *windows*)
            return 0
            ;;
    esac
    
    # Check if install directory exists
    if [ -d "$install_dir" ]; then
        # Directory exists - check if writable
        if [ ! -w "$install_dir" ]; then
            log_error "No write permission for $install_dir"
            local uid=$(id -u 2>/dev/null || echo "0")
            if [ "$uid" -ne 0 ]; then
                log_error "Please run with sudo or choose a different install location"
                log_error "Example: INSTALL_DIR=\$HOME/.local/bin $0"
            fi
            exit 1
        fi
    else
        # Directory doesn't exist - check parent directory
        local parent_dir=$(dirname "$install_dir")
        if [ ! -d "$parent_dir" ]; then
            log_error "Parent directory does not exist: $parent_dir"
            exit 1
        fi
        if [ ! -w "$parent_dir" ]; then
            log_error "No write permission for parent directory: $parent_dir"
            local uid=$(id -u 2>/dev/null || echo "0")
            if [ "$uid" -ne 0 ]; then
                log_error "Please run with sudo or choose a different install location"
                log_error "Example: INSTALL_DIR=\$HOME/.local/bin $0"
            fi
            exit 1
        fi
    fi
}

# Check if temp/download directory is writable
check_temp_permissions() {
    local temp_dir="$1"
    local parent_dir=$(dirname "$temp_dir")
    
    if [ ! -d "$parent_dir" ]; then
        log_error "Parent directory does not exist: $parent_dir"
        exit 1
    fi
    
    if [ ! -w "$parent_dir" ]; then
        log_error "No write permission for temporary directory location: $parent_dir"
        log_error "Cannot create temporary files for installation"
        exit 1
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
    
    # Check install directory permissions BEFORE proceeding
    log_info "Checking installation directory permissions: $INSTALL_DIR"
    check_install_permissions "$INSTALL_DIR"
    
    # Determine binary name
    local binary_suffix=""
    case "$platform" in
        *windows*)
            binary_suffix=".exe"
            ;;
    esac
    
    local binary_name="duh-${platform}${binary_suffix}"
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${binary_name}"
    local checksum_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/checksums.txt"
    
    # Create temporary directory
    local tmp_dir
    case "$platform" in
        *windows*)
            # Use current directory for Windows - /tmp has issues with executables in Git Bash
            tmp_dir="./duh-install-$$"
            ;;
        *)
            tmp_dir=$(mktemp -d)
            ;;
    esac
    
    # Check temp directory parent permissions BEFORE trying to create
    log_info "Checking temporary directory permissions"
    check_temp_permissions "$tmp_dir"
    
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
    
    # Create install directory if it doesn't exist
    if ! mkdir -p "$INSTALL_DIR"; then
        log_error "Failed to create installation directory: $INSTALL_DIR"
        exit 1
    fi
    
    # Install binary
    local install_path="$INSTALL_DIR/$BINARY_NAME$binary_suffix"
    
    # Check if this is an update or fresh install
    local is_update=false
    
    if is_duh_installed "$install_path"; then
        log_info "Updating existing Duh installation"
        is_update=true
    else
        log_info "Installing Duh $version to $install_path"
    fi
    
    cp "$binary_path" "$install_path"
    chmod +x "$install_path"
    
    # Success message based on whether this was an update
    if [ "$is_update" = true ]; then
        log_success "Duh successfully updated to $version!"
    else
        log_success "Duh $version installed successfully!"
    fi
    log_info "Location: $install_path"
    
    # Manual cleanup for better user feedback
    cleanup
    
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