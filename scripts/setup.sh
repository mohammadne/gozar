#!/bin/bash

output_directory="./outputs"

# gozar_version="v0.2.0"
# podman run -d --name gozar ghcr.io/mohammadne/gozar:$gozar_version
# podman exec -it gozar /app/entrypoint executer
# podman cp gozar:/app/outputs $output_directory
# podman rm -f gozar -t 0

server_xray_directory="/root/xray"

ssh proxy-server -t 'rm -rf '"$server_xray_directory"''
scp -r "$output_directory/server" proxy-server:$server_xray_directory
ssh proxy-server -t 'cd '"$server_xray_directory"' && docker compose -f compose.todo.yml up -d'

# cd "$output_directory/client" && podman compose -f compose-notls.yml up -d
podman run --restart=always -d --name "xray-notls" -p 10808:10808 -v "./notls.json:/usr/local/xray/config.json:z" "xray" run -c config.json

podman run --restart=always -d --name "xray-reality" -p 10809:10809 -v "./reality.json:/usr/local/xray/config.json:z" "xray" run -c config.json
