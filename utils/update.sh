#!/bin/bash

# Generate slides.json by scanning YYYY/MM/slug/index.html structure

echo "Scanning slides..."

entries=()

while IFS= read -r file; do
  [ -z "$file" ] && continue

  year=$(echo "$file" | cut -d'/' -f2)
  month=$(echo "$file" | cut -d'/' -f3)
  slug=$(echo "$file" | cut -d'/' -f4)

  # Extract title from <title> tag
  title=$(grep -oP '(?<=<title>).*(?=</title>)' "$file" 2>/dev/null | head -1)
  [ -z "$title" ] && title="$slug"

  entries+=("  { \"year\": \"$year\", \"month\": \"$month\", \"slug\": \"$slug\", \"title\": \"$title\" }")
done < <(find . -path "./.git" -prune -o -type f -name "index.html" -print 2>/dev/null | \
  grep -E "^./[0-9]{4}/[0-9]{2}/[^/]+/index.html$" | sort -r)

# Write JSON
echo "[" > slides.json
for i in "${!entries[@]}"; do
  if [ $i -lt $((${#entries[@]} - 1)) ]; then
    echo "${entries[$i]}," >> slides.json
  else
    echo "${entries[$i]}" >> slides.json
  fi
done
echo "]" >> slides.json

echo "Generated slides.json:"
cat slides.json
