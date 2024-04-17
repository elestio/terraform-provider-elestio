#!/bin/bash

api_endpoint="https://api.elest.io/api/servers/getTemplates?id=-1"
output_file="internal/provider/templates.json"

# Check if jq is installed
if ! [ -x "$(command -v jq)" ]; then
  echo "Error: jq is not installed. Please install jq and try again."
  exit 1
fi

# Make HTTP request and save response to a variable
response=$(curl -s "$api_endpoint")

# Check if HTTP request was successful
if [ $? -ne 0 ]; then
  echo "Error: HTTP request failed. Please check your internet connection and try again."
  exit 1
fi

# Check if response is valid JSON
if ! echo "$response" | jq . > /dev/null 2>&1; then
  echo "Error: Response is not valid JSON. Please check the API endpoint and try again."
  exit 1
fi

# Rename "instances" to "templates"
response=$(echo "$response" | jq '. + {"templates":.instances}|del(.instances)')

# Remove templates with "Full Stack" category
response=$(echo "$response" | jq '. + {"templates":.templates|map(select(.category != "Full Stack"))}')

# # Add "https:" in front of all image URLs
response=$(echo "$response" | jq '. + {"templates":.templates|map(.mainImage |= sub("^"; "https:"))}')

# Output the modified JSON to a local file
echo "$response" > "$output_file"

echo "JSON file created successfully! $output_file"