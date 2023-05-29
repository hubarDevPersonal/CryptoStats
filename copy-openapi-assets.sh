#!/usr/bin/env bash

# Drop existing folders
rm -rf out/api/ > /dev/null 2>&1 || true

# Create directories if they don't exist
mkdir -p out/api/

# Copy OpenAPI assets
cp -f api/gen/http/openapi* out/api/

cp -f api/www/index.html out/api/index.html
