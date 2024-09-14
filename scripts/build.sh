#!/bin/bash

# Get the script directory
scriptdir=$(dirname "$(realpath "$0")")

# Function to check if Go is installed
check_go_installed() {
    if ! command -v go &> /dev/null; then
        echo "Go is not installed or not in your PATH. Please install Go before proceeding."
        exit 1
    fi
}

# Function to build the API binary
build_api() {
    srcdir=$(dirname "$scriptdir")

    # Change to the source directory
    cd "$srcdir" || { echo "Failed to change directory to $srcdir"; exit 1; }

    local output_dir="$srcdir/bin"
    local build_cmd_dir="$srcdir/cmd/api"

    echo "Building API binary..."

    # Ensure the output directory exists
    mkdir -p "$output_dir"

    # Check if the source directory exists
    if [ ! -d "$build_cmd_dir" ]; then
        echo "API source directory '$build_cmd_dir' does not exist. Check the path."
        exit 1
    fi

    # Build the binary
    go build -o "$output_dir/incrate" "$build_cmd_dir"
    if [ $? -ne 0 ]; then
        echo "Failed to build API binary. Check for errors in the build process."
        exit 1
    fi

    echo "API binary built successfully at $output_dir/incrate"
}

# Ensure Go is installed
check_go_installed

# Call the function to build the API binary
build_api