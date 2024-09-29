#!/usr/bin/env bash

WORKFLOW_DIR=$(realpath "${BASH_SOURCE%/*}")
OUTPUT_FILE=${WORKFLOW_DIR}/update.log

while getopts "w:o:" opt; do
    case $opt in
        w)
            WORKFLOW_DIR="$OPTARG"
            ;;
        o)
            OUTPUT_FILE="$OPTARG"
            ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            exit 1
            ;;
    esac
done

shift $((OPTIND-1))

sed -e "s|@WORKFLOW_DIR@|${WORKFLOW_DIR}|g; 
s|@OUTPUT_FILE@|${OUTPUT_FILE}|g" "${WORKFLOW_DIR}/update.crontab"
