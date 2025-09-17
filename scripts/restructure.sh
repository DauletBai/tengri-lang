#!/bin/bash

# --- "tenge" Project Rebranding and Restructuring Script ---

set -e # Exit immediately if a command fails

echo "[REBRAND] Starting project rename: tenge -> tenge"

# 1. Rename files and directories
# We use 'find' to locate all files and dirs with 'tenge' in their name
# and rename them using a loop.
echo "[REBRAND] Renaming files and directories..."
find . -depth -name '*tenge*' | while read -r file; do
    new_file=$(echo "$file" | sed 's/tenge/tenge/g')
    mv "$file" "$new_file"
    echo "  Renamed: $file -> $new_file"
done

# 2. Rename file extensions from .tng to .tng
echo "[REBRAND] Renaming file extensions .tng -> .tng..."
find ./benchmarks -depth -name '*.tng' | while read -r file; do
    new_file="${file%.tng}.tng"
    mv "$file" "$new_file"
    echo "  Renamed: $file -> $new_file"
done

# 3. Update file contents
# Use 'grep' to find all files containing the string "tenge" (case-insensitive)
# and use 'sed' to replace it with "tenge".
echo "[REBRAND] Updating file contents..."
grep -rli 'tenge' . --exclude-dir=.git --exclude-dir=.bin | while read -r file; do
    sed -i '' 's/tenge/tenge/gI' "$file"
    echo "  Updated: $file"
done

# 4. Update .tng extension to .tng in contents
echo "[REBRAND] Updating .tng references to .tng..."
grep -rli '\.tng' . --exclude-dir=.git --exclude-dir=.bin | while read -r file; do
    sed -i '' 's/\.tng/\.tng/g' "$file"
    echo "  Updated: $file"
done


echo "[REBRAND] Project rebranding to 'tenge' is complete."
echo "NOTE: You may need to rename the root 'tenge-lang' directory manually."