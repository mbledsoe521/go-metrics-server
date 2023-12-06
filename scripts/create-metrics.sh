#!/bin/bash

# Define the number of key-value pairs
num_pairs=20

# Define the output file path
output_file="/data/metrics_from_special_app.txt"

# Ensure the output directory exists
mkdir -p "$(dirname "$output_file")"

# Open the output file for writing
exec > "$output_file"

# Generate and write key-value pairs
for i in $(seq "$num_pairs"); do
  key="CPU$i"
  value=$((RANDOM % 101))
  printf "%s=%s\n" "$key" "$value"
done

