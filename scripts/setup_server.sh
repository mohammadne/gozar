#!/bin/bash

output_directory="./outputs"

# gozar_version="v0.2.0"
# podman run -d --name gozar ghcr.io/mohammadne/gozar:$gozar_version
# podman exec -it gozar /app/entrypoint executer
# podman cp gozar:/app/outputs $output_directory
# podman rm -f gozar -t 0

server_xray_directory="/root/xray"

ssh gozar-server -t 'rm -rf '"$server_xray_directory"''
scp -r "$output_directory/server" gozar-server:$server_xray_directory
ssh gozar-server -t 'cd '"$server_xray_directory"' && docker compose -f compose.yml up -d'
