#!/bin/bash

# Script to update static JavaScript libraries
# Downloads latest versions from CDN and updates local files

STATIC_DIR="static"
mkdir -p "$STATIC_DIR"

echo "Updating static JavaScript libraries..."

# Function to download and update a file if it's different
update_file() {
    local url=$1
    local file=$2
    local temp_file="${file}.tmp"
    
    echo "Checking $file..."
    
    # Download to temp file
    if curl -s -L "$url" -o "$temp_file"; then
        # Check if file exists and compare
        if [ -f "$file" ]; then
            if cmp -s "$file" "$temp_file"; then
                echo "  $file is up to date"
                rm "$temp_file"
            else
                echo "  Updating $file"
                mv "$temp_file" "$file"
            fi
        else
            echo "  Creating $file"
            mv "$temp_file" "$file"
        fi
    else
        echo "  Failed to download $file"
        rm -f "$temp_file"
    fi
}

# Update HTMx
update_file "https://cdn.jsdelivr.net/npm/htmx.org@latest/dist/htmx.min.js" "$STATIC_DIR/htmx.min.js"

# Update Alpine.js
update_file "https://cdn.jsdelivr.net/npm/alpinejs@latest/dist/alpine.min.js" "$STATIC_DIR/alpine.js"

# Update xterm.js
update_file "https://cdn.jsdelivr.net/npm/xterm@latest/lib/xterm.js" "$STATIC_DIR/xterm.js"

# Update xterm.css
update_file "https://cdn.jsdelivr.net/npm/xterm@latest/css/xterm.css" "$STATIC_DIR/xterm.css"

echo "Done!"

