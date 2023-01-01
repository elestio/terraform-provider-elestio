#!/bin/bash

templates_file="internal/provider/templates.json"

# Check if jq is installed
if ! [ -x "$(command -v jq)" ]; then
  echo "Error: jq is not installed. Please install jq and try again."
  exit 1
fi

clean_string () {
  string="$1"
  string=$(echo "$string" | tr '[:upper:]' '[:lower:]' | sed 's/[^[:alnum:]]/_/g' | tr -s _)
  echo "$string"
}

templates=$(cat $templates_file)

while read -r template; do
  decoded=$(echo "$template" | base64 --decode)

  # Skip templates from "Full Stack" category
  [[ $(echo "$decoded" | jq -r '.category') == "Full Stack" ]] && continue

  resourceName=$(clean_string $(echo "$decoded" | jq -r '.title'))
  documentationName=$(echo "$decoded" | jq -r '.title')
  defaultVersion=$(echo "$decoded" | jq -r '.dockerhub_default_tag')

  dirPath="./examples/resources/elestio_$resourceName/"
  mkdir -p "$dirPath"

  import=$(cat elestio-templates/examples/import.sh | sed "s/\[TEMPLATE_RESOURCE_NAME\]/$resourceName/g" | sed "s/\[TEMPLATE_DOCUMENTATION_NAME\]/$documentationName/g")
  echo "$import" > "$dirPath/import.sh"

  resource=$(cat elestio-templates/examples/resource.tf | sed "s/\[TEMPLATE_RESOURCE_NAME\]/$resourceName/g" | sed "s/\[TEMPLATE_DOCUMENTATION_NAME\]/$documentationName/g" | sed "s/\[TEMPLATE_DEFAULT_VERSION\]/$defaultVersion/g")
  echo "$resource" > "$dirPath/resource.tf"

done < <(jq -r '.templates[] | @base64' <<< "$templates")