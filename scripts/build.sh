#!/usr/bin/env bash
set -euo pipefail

# Build script for marcli
# Runs the build command via go run

cd "$(dirname "$0")/.."
go run . build

