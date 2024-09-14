#!/bin/bash

scriptdir=$(dirname "$(realpath "$0")")

# Function to build Docker image
build_image() {
    srcdir=$(dirname "$scriptdir")

    # Change to the source directory
    cd "$srcdir" || { echo "Failed to change directory to $srcdir"; exit 1; }

    # Build Docker image
    docker build -f "$srcdir/Dockerfile" .
}

# Function to check Docker user permissions
check_docker_permissions() {
    # Check if the user is part of the 'docker' group or has sudo privileges
    if ! docker info >/dev/null 2>&1; then
        echo "You do not have permission to use Docker."
        echo "Please run the script with 'sudo' or ensure your user is added to the 'docker' group."
        exit 1
    fi
}

# Check Docker permissions
check_docker_permissions

# Call the function to build the Docker image
build_image