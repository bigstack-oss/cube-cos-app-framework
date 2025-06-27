#!/bin/bash

# JSON file containing image data
json_file="images.json"

# Host to mirror images to
CUBE_VIP="10.32.10.180"

# Loop through each image in the JSON file and pull it
jq -c '.[]' "$json_file" | while read -r image; do
    fqdn=$(echo "$image" | jq -r '.fqdn')
    name=$(echo "$image" | jq -r '.name')
    space=$(echo "$image" | jq -r '.space')
    tag=$(echo "$image" | jq -r '.tag')
    
    skopeo copy docker-archive:pkgs/${name}-${tag}.tar docker://${CUBE_VIP}:5080/${space}/${name}:${tag} --dest-tls-verify=false --preserve-digests
done
