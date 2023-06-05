#!/bin/bash

set -euo pipefail

hashfile() {
  local file="$1"
  local html_file="$2"
  local directory
  directory="$(dirname "${file}")"
  local filename_with_extension
  filename_with_extension="$(basename "${file}")"
  local filename_without_extension="${filename_with_extension%%.*}"
  local extension="${file#*.}"

  filehash="$(md5sum "${file}" | cut -c1-8)"
  mv "${file}" "${directory}/${filename_without_extension}-${filehash}.${extension}"
  sed -i "s/$filename_with_extension/$filename_without_extension-$filehash.$extension/g" "${html_file}"
}

hashfile "public/scripts.js" "public/index.html"
hashfile "public/style.css" "public/index.html"
