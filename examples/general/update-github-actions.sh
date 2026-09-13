#!/bin/bash

# Title: Update GitHub Actions to use latest major versions

# This script updates common GitHub Actions to their latest major versions
# Useful for keeping CI/CD pipelines up to date and secure
#
# Usage: Place in your repository root and run with multi-gitter
# It will update actions/checkout, actions/setup-go, actions/setup-node, etc.
#
# Example:
#   multi-gitter run ./update-github-actions.sh -O my-org -m "chore: update GitHub Actions to latest versions"

set -e

# Check if .github/workflows directory exists
if [ ! -d ".github/workflows" ]; then
    echo "No .github/workflows directory found, skipping"
    exit 0
fi

# Update common GitHub Actions to their latest major versions
# Using sed to replace old action versions with new ones

# actions/checkout v3 -> v4
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: actions\/checkout@v3/uses: actions\/checkout@v4/g' 2>/dev/null || true

# actions/setup-go v4 -> v5
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: actions\/setup-go@v4/uses: actions\/setup-go@v5/g' 2>/dev/null || true

# actions/setup-go v3 -> v5 (skip v4)
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: actions\/setup-go@v3/uses: actions\/setup-go@v5/g' 2>/dev/null || true

# actions/setup-node v3 -> v4
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: actions\/setup-node@v3/uses: actions\/setup-node@v4/g' 2>/dev/null || true

# actions/setup-python v4 -> v5
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: actions\/setup-python@v4/uses: actions\/setup-python@v5/g' 2>/dev/null || true

# actions/setup-python v3 -> v5
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: actions\/setup-python@v3/uses: actions\/setup-python@v5/g' 2>/dev/null || true

# actions/upload-artifact v3 -> v4
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: actions\/upload-artifact@v3/uses: actions\/upload-artifact@v4/g' 2>/dev/null || true

# actions/download-artifact v3 -> v4
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: actions\/download-artifact@v3/uses: actions\/download-artifact@v4/g' 2>/dev/null || true

# docker/login-action v2 -> v3
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: docker\/login-action@v2/uses: docker\/login-action@v3/g' 2>/dev/null || true

# docker/setup-buildx-action v2 -> v3
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: docker\/setup-buildx-action@v2/uses: docker\/setup-buildx-action@v3/g' 2>/dev/null || true

# docker/build-push-action v4 -> v5
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: docker\/build-push-action@v4/uses: docker\/build-push-action@v5/g' 2>/dev/null || true

# docker/build-push-action v3 -> v5
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: docker\/build-push-action@v3/uses: docker\/build-push-action@v5/g' 2>/dev/null || true

# actions/cache v3 -> v4
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: actions\/cache@v3/uses: actions\/cache@v4/g' 2>/dev/null || true

# codecov/codecov-action v3 -> v4
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: codecov\/codecov-action@v3/uses: codecov\/codecov-action@v4/g' 2>/dev/null || true

# sigstore/cosign-installer v2/v3 -> v3
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: sigstore\/cosign-installer@v2/uses: sigstore\/cosign-installer@v3/g' 2>/dev/null || true

# goreleaser/goreleaser-action v4 -> v5
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: goreleaser\/goreleaser-action@v4/uses: goreleaser\/goreleaser-action@v5/g' 2>/dev/null || true

# goreleaser/goreleaser-action v3 -> v5
find .github/workflows -name "*.yml" -o -name "*.yaml" | \
    xargs sed -i 's/uses: goreleaser\/goreleaser-action@v3/uses: goreleaser\/goreleaser-action@v5/g' 2>/dev/null || true

echo "GitHub Actions updated successfully!"
