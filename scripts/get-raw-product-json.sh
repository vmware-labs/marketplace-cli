#!/usr/bin/env bash

# This script wraps around the CLI to get the JSON structure of a product directly from the Marketplace without
# trying to unmarshal it into the structure.
# This is helpful in diagnosing times when the Marketplace adds a field that the Marketplace CLI does not know about yet.

if [[ -z "$1" ]]; then
  echo "USAGE: $0 <product slug>"
  exit 1
fi

mkpcli curl "/api/v1/products/${1}?increaseViewCount=false&isSlug=true"
