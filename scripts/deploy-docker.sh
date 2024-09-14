#!/bin/bash

scriptdir=$(dirname "$(realpath "$0")")

srcdir=$(dirname "$scriptdir")

# Function to check Docker user permissions
check_docker_permissions() {
    # Check if the user is part of the 'docker' group or has sudo privileges
    if ! docker info >/dev/null 2>&1; then
        echo "You do not have permission to use Docker."
        echo "Please run the script with 'sudo' or ensure your user is added to the 'docker' group."
        exit 1
    fi
}

# Function to deploy using Docker Compose
deploy_docker_compose() {
    local compose_file=$1
    local detach_mode=$2
    
    # Change to the source directory
    cd "$srcdir" || { echo "Failed to change directory to $srcdir"; exit 1; }

    # Set detach flag based on user input
    if [ "$detach_mode" == true ]; then
        docker compose -f "$compose_file" up --build -d
    else
        docker compose -f "$compose_file" up --build
    fi

    if [ $? -ne 0 ]; then
        echo "Failed to deploy using Docker Compose. Check for errors."
        exit 1
    fi
}

# Default values
compose_file="$srcdir/docker-compose.yml"
detach_mode=false

# Parse command line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -f|--file) compose_file="$2"; shift ;;
        -d|--detach) detach_flag=true ;;
        *) echo "Unknown parameter: $1"; exit 1 ;;
    esac
    shift
done

# Ensure Docker permissions
check_docker_permissions

# Deploy with Docker Compose
deploy_docker_compose "$compose_file" "$detach_mode"