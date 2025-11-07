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

generate_firewall_rules() {
  local firewallPorts="$1"
  local tmpfile="$2"
  
  if [[ -z "$firewallPorts" ]] || [[ "$firewallPorts" == "null" ]]; then
    return
  fi
  
  echo "  # Default firewall rules" >> "$tmpfile"
  echo "  firewall_user_rules = [" >> "$tmpfile"
  echo "    # Required system ports" >> "$tmpfile"
  echo "    {" >> "$tmpfile"
  echo "      \"type\"     = \"input\"" >> "$tmpfile"
  echo "      \"port\"     = \"22\"" >> "$tmpfile"
  echo "      \"protocol\" = \"tcp\"" >> "$tmpfile"
  echo "      \"targets\"  = [\"0.0.0.0/0\", \"::/0\"]" >> "$tmpfile"
  echo "    }," >> "$tmpfile"
  echo "    {" >> "$tmpfile"
  echo "      \"type\"     = \"input\"" >> "$tmpfile"
  echo "      \"port\"     = \"4242\"" >> "$tmpfile"
  echo "      \"protocol\" = \"udp\"" >> "$tmpfile"
  echo "      \"targets\"  = [\"0.0.0.0/0\", \"::/0\"]" >> "$tmpfile"
  echo "    }," >> "$tmpfile"
  
  # Parse firewall ports and add them
  IFS=',' read -ra PORTS <<< "$firewallPorts"
  if [[ ${#PORTS[@]} -gt 0 ]]; then
    echo "    # Application ports" >> "$tmpfile"
    for port in "${PORTS[@]}"; do
      # Trim whitespace using bash parameter expansion
      port="${port## }"
      port="${port%% }"
      [[ -z "$port" ]] && continue
      # Determine protocol (default to tcp)
      protocol="tcp"
      echo "    {" >> "$tmpfile"
      echo "      \"type\"     = \"input\"" >> "$tmpfile"
      echo "      \"port\"     = \"$port\"" >> "$tmpfile"
      echo "      \"protocol\" = \"$protocol\"" >> "$tmpfile"
      echo "      \"targets\"  = [\"0.0.0.0/0\", \"::/0\"]" >> "$tmpfile"
      echo "    }," >> "$tmpfile"
    done
  fi
  
  # Remove trailing comma from last element
  sed -i '' '$ s/,$//' "$tmpfile"
  
  echo "  ]" >> "$tmpfile"
}

while read -r template; do
  decoded=$(echo "$template" | base64 --decode)

  # Skip templates from "Full Stack" category
  [[ $(echo "$decoded" | jq -r '.category') == "Full Stack" ]] && continue

  resourceName=$(clean_string $(echo "$decoded" | jq -r '.title'))
  documentationName=$(echo "$decoded" | jq -r '.title')
  defaultVersion=$(echo "$decoded" | jq -r '.dockerhub_default_tag')
  firewallPorts=$(echo "$decoded" | jq -r '.firewallPorts')

  dirPath="./examples/resources/elestio_$resourceName/"
  mkdir -p "$dirPath"

  import=$(cat elestio-templates/examples/import.sh | sed "s/\[TEMPLATE_RESOURCE_NAME\]/$resourceName/g" | sed "s/\[TEMPLATE_DOCUMENTATION_NAME\]/$documentationName/g")
  echo "$import" > "$dirPath/import.sh"

  # Generate resource file with firewall rules
  tmpResourceFile=$(mktemp)
  cat elestio-templates/examples/resource.tf | \
    sed "s/\[TEMPLATE_RESOURCE_NAME\]/$resourceName/g" | \
    sed "s/\[TEMPLATE_DOCUMENTATION_NAME\]/$documentationName/g" | \
    sed "s/\[TEMPLATE_DEFAULT_VERSION\]/$defaultVersion/g" > "$tmpResourceFile"
  
  # Insert firewall rules if we have them
  if [[ -n "$firewallPorts" ]] && [[ "$firewallPorts" != "null" ]]; then
    tmpFirewallFile=$(mktemp)
    generate_firewall_rules "$firewallPorts" "$tmpFirewallFile"
    
    # Replace placeholder with firewall rules
    awk '/\[TEMPLATE_FIREWALL_RULES\]/ { system("cat '"$tmpFirewallFile"'"); next } 1' "$tmpResourceFile" > "$dirPath/resource.tf"
    rm -f "$tmpFirewallFile"
  else
    # Remove placeholder if no firewall rules
    sed '/\[TEMPLATE_FIREWALL_RULES\]/d' "$tmpResourceFile" > "$dirPath/resource.tf"
  fi
  
  rm -f "$tmpResourceFile"

done < <(jq -r '.templates[] | @base64' <<< "$templates")