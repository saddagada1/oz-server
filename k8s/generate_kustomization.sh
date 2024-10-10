#!/bin/bash

# Ensure all necessary tools are installed
command -v envsubst >/dev/null 2>&1 || { echo >&2 "envsubst is required but it's not installed. Aborting."; exit 1; }
command -v sed >/dev/null 2>&1 || { echo >&2 "sed is required but it's not installed. Aborting."; exit 1; }
command -v base64 >/dev/null 2>&1 || { echo >&2 "base64 is required but it's not installed. Aborting."; exit 1; }

# Read and base64 encode the POSTGRES_PASSWORD
postgres_password=$(grep '^POSTGRES_PASSWORD=' .env | cut -d'=' -f2- | tr -d "'")
encoded_password=$(echo -n $postgres_password | base64)

# Export environment variables
export POSTGRES_PASSWORD_BASE64=$encoded_password

# Read environment variables from .env file and substitute them in the template
env $(cat .env | grep -v '^#' | sed 's/\([^=]*\)=\(.*\)/\1="\2"/' | xargs) envsubst < kustomization-template.yaml > kustomization.yaml

echo "Kustomization file generated successfully."
