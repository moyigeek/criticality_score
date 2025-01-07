#!/bin/bash

CONFIG_PATH=$1

if [ -z "$CONFIG_PATH" ]; then
  echo "Error: config.json path is required as the first argument."
  exit 1
fi

DB_HOST=$(jq -r '.host' "$CONFIG_PATH")
DB_PORT=$(jq -r '.port' "$CONFIG_PATH")
DB_NAME=$(jq -r '.database' "$CONFIG_PATH")
DB_USER=$(jq -r '.user' "$CONFIG_PATH")
DB_PASSWORD=$(jq -r '.password' "$CONFIG_PATH")

export PGPASSWORD="$DB_PASSWORD"

DEBIAN_QUERY="SELECT package FROM debian_packages;"
DEEPIN_QUERY="SELECT package FROM deepin_packages;"
UBUNTU_QUERY="SELECT package FROM ubuntu_packages;"

debian_packages=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "$DEBIAN_QUERY" | awk '{print $1}')
deepin_packages=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "$DEEPIN_QUERY" | awk '{print $1}')
ubuntu_packages=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "$UBUNTU_QUERY" | awk '{print $1}')

declare -A debian_set
declare -A deepin_set
declare -A ubuntu_set
while IFS= read -r package; do
  debian_set["$package"]=true
done <<< "$debian_packages"

while IFS= read -r package; do
  deepin_set["$package"]=true
done <<< "$deepin_packages"

while IFS= read -r package; do
  ubuntu_set["$package"]=true
done <<< "$ubuntu_packages"

common_deepin_count=0
common_ubuntu_count=0

for debian_pkg in "${!debian_set[@]}"; do
  common_count=0
  if [[ ${deepin_set["$debian_pkg"]} ]]; then
    common_count=$((common_count + 1))
  fi
  if [[ ${ubuntu_set["$debian_pkg"]} ]]; then
    common_count=$((common_count + 1))
  fi

  if [[ $common_count -gt 0 ]]; then
    if [[ ${deepin_set["$debian_pkg"]} ]]; then
      common_deepin_count=$((common_deepin_count + 1))
    fi
    if [[ ${ubuntu_set["$debian_pkg"]} ]]; then
      common_ubuntu_count=$((common_ubuntu_count + 1))
    fi
  fi
done

if [[ $common_deepin_count -gt 0 ]]; then
  echo "Number of common packages between Debian and Deepin: $common_deepin_count"
else
  echo "No common packages found between Debian and Deepin."
fi

if [[ $common_ubuntu_count -gt 0 ]]; then
  echo "Number of common packages between Debian and Ubuntu: $common_ubuntu_count"
else
  echo "No common packages found between Debian and Ubuntu."
fi

unset PGPASSWORD
