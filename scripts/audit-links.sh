#!/usr/bin/env bash

tmpfile=$(mktemp)
trap 'rm -f "$tmpfile"' EXIT

# Find all markdown files and check links
find . -name "*.md" -not -path "./.beads/*" | while read -r file; do
    # Use perl to extract markdown links properly
    perl -ne 'while (/\[([^\]]+)\]\(([^)]+)\)/g) { print "$.:$2\n"; }' "$file" | while IFS=: read -r line_num link; do
        # Check external links
        if echo "$link" | grep -qE '^https?://'; then
            if ! curl -s -f -L --max-time 10 "$link" > /dev/null 2>&1; then
                echo "BROKEN: $file:$line_num -> $link (external link unreachable)" | tee -a "$tmpfile"
            fi
            continue
        fi
        
        # Remove anchor
        file_path="${link%%#*}"
        
        # Resolve relative to source file's directory
        source_dir=$(dirname "$file")
        resolved_path="$source_dir/$file_path"
        
        if [[ ! -f "$resolved_path" ]]; then
            echo "BROKEN: $file:$line_num -> $file_path (resolved: $resolved_path)" | tee -a "$tmpfile"
        fi
    done
done

if [[ -s "$tmpfile" ]]; then
    broken=$(wc -l < "$tmpfile" | tr -d ' ')
    echo ""
    echo "Found $broken broken link(s)"
    exit 1
fi

echo "All links valid"
