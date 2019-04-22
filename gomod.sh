#!/usr/bin/env bash
set -eu

touch go.mod

PROJECT_NAME="paul-nelson-baker"
CURRENT_DIR=$(basename $(pwd))

#require github.com/aws/aws-lambda-go v1.6.0
CONTENT=$(cat <<-EOD
module github.com/${PROJECT_NAME}/${CURRENT_DIR}

EOD
)

echo "$CONTENT" > go.mod
