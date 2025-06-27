#!/bin/bash

# JSON file containing image data
json_file="images.json"

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "jq is required but not installed. Please install jq and try again."
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Docker is required but not installed. Please install Docker and try again."
    exit 1
fi

# Loop through each image in the JSON file and pull it
jq -c '.[]' "$json_file" | while read -r image; do
    # Parse JSON object fields
    fqdn=$(echo "$image" | jq -r '.fqdn')
    name=$(echo "$image" | jq -r '.name')
    tag=$(echo "$image" | jq -r '.tag')
    
    echo "Processing image: $fqdn"
    echo "Name: $name, Tag: $tag"

    # Pull the Docker image
    docker pull --platform linux/amd64 "$fqdn"

    # Save the image to a tar file
    output_file="${name}-${tag}.tar"
    docker save "$fqdn" -o "pkgs/$output_file"
done
