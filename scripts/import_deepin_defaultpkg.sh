#!/bin/bash

if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <pkglist> <config.json>"
    exit 1
fi

PKGLIST_FILE="$1"

if [ ! -f "$PKGLIST_FILE" ]; then
    echo "pkglist file not found!"
    exit 1
fi

CONFIG_PATH=$2
if [ ! -f "$CONFIG_PATH" ]; then
    echo "config.json file not found!"
    exit 1
fi

HOST=$(jq -r '.host' "$CONFIG_PATH")
PORT=$(jq -r '.port' "$CONFIG_PATH")
DATABASE=$(jq -r '.database' "$CONFIG_PATH")
USER=$(jq -r '.user' "$CONFIG_PATH")
PASSWORD=$(jq -r '.password' "$CONFIG_PATH")

if [ -z "$DATABASE" ] || [ -z "$USER" ] || [ -z "$PASSWORD" ] || [ -z "$HOST" ] || [ -z "$PORT" ]; then
    echo "Failed to retrieve database connection details from config.json!"
    exit 1
fi

export PGPASSWORD="$PASSWORD"

while IFS= read -r pkg_name; do

    if [ "$pkg_info" != "null" ]; then
        echo "Package found: $pkg_name"

        psql -h "$HOST" -p "$PORT" -U "$USER" -d "$DATABASE" -c "
            UPDATE deepin_packages
            SET default_install = 1
            WHERE package = '$pkg_name';
        "

        echo "Updated default_install to 1 for $pkg_name"
    else
        echo "Package $pkg_name not found in config.json"
    fi
done < "$PKGLIST_FILE"
unset PGPASSWORD
