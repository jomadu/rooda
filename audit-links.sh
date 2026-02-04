#!/usr/bin/env bash
set -euo pipefail

broken=0

check_link() {
    local source_file="$1"
    local link_target="$2"
    local line_num="$3"
    
    # Remove anchor if present
    local file_path="${link_target%%#*}"
    
    # Skip external links
    if echo "$file_path" | grep -qE '^https?://'; then
        return 0
    fi
    
    # Resolve relative to source file's directory
    local source_dir
    source_dir="$(dirname "$source_file")"
    local resolved_path="$source_dir/$file_path"
    
    if [[ ! -f "$resolved_path" ]]; then
        echo "BROKEN: $source_file:$line_num -> $file_path (resolved: $resolved_path)"
        ((broken++))
    fi
}

# Find all markdown files
while IFS= read -r file; do
    line_num=0
    while IFS= read -r line; do
        ((line_num++))
        # Extract markdown links: [text](path)
        echo "$line" | grep -oE '\[[^]]+\]\([^)]+\)' | while read -r match; do
            link_target=$(echo "$match" | sed -E 's/\[([^]]+)\]\(([^)]+)\)/\2/')
            check_link "$file" "$link_target" "$line_num"
        done
    done < "$file"
done < <(find . -name "*.md" -not -path "./.beads/*" -not -path "./node_modules/*")

if [[ $broken -gt 0 ]]; then
    echo "Found $broken broken link(s)"
    exit 1
fi

echo "All links valid"
exit 0
