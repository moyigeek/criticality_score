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

SQL_QUERY="SELECT git_link FROM git_metrics WHERE created_since = '0001-01-01' or contributor_count = 0"

export PGPASSWORD="$DB_PASSWORD"

echo "Fetching links from the database..."
links=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "$SQL_QUERY")

if [ -z "$links" ]; then
  echo "No links found with created_since = '0001-01-01'. Exiting."
  exit 1
fi

MAX_CONCURRENT=50
SEMAPHORE=0

function process_link() {
  link=$1
  echo "Processing link: $link"

  repo_path=$(echo "$link" | sed 's|^[a-zA-Z]*://||')
  local_repo_path="/home/ruijie/Storage/$repo_path"
  file_link="$local_repo_path"
  if [ -d "$file_link" ] || [ -f "$file_link" ]; then
    ./bin/cli --config ../config.json --update-db "$file_link"
  else
    ./bin/cli --config ../config.json --update-db "$link"
  fi

  if [ $? -ne 0 ]; then
    echo "Failed to process $link. Skipping."
  else
    echo "Successfully processed $link."
  fi

  SEMAPHORE=$((SEMAPHORE - 1))  # Decrease semaphore after the process finishes
}

for link in $links; do
  link=$(echo "$link" | xargs)

  # Wait if maximum concurrency is reached
  while [ $SEMAPHORE -ge $MAX_CONCURRENT ]; do
    sleep 1
  done

  SEMAPHORE=$((SEMAPHORE + 1))
  process_link "$link" &  # Launch in background immediately without waiting

done

wait  # Ensure all background tasks complete

unset PGPASSWORD

echo "Script completed."
