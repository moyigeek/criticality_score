#!/usr/bin/env bash


SCRIPT_DIR=$(realpath "${BASH_SOURCE%/*}")

ensure_exists() {
    if [ ! -f "$1" ]; then
        echo "! Creating $1"
        touch "$1"
    fi
}

prepare() {
    ensure_exists package_updated.src
    ensure_exists git_updated.src
    ensure_exists depsdev_updated.src
    ensure_exists gitlink_updated.src
}


run() {
    touch "$1_updated.src"
    make -s -e -f "$SCRIPT_DIR/flow.mk"
}

# getopts -C : change dir

while getopts "C:" opt; do
    case $opt in
        C)
            cd "$OPTARG" || exit 1
            ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            exit 1
            ;;
    esac
done

shift $((OPTIND-1))

# $1 must be 'package' or 'git' or 'depsdev' or 'gitlink'

if [ "$1" = "package" ] || [ "$1" = "git" ] || [ "$1" = "depsdev" ] || [ "$1" = "gitlink" ]; then
    prepare
    run "$1"
else
    if [ -z "$1" ]; then 
        echo "No argument supplied, please specify one of: package, git, depsdev, gitlink"
    else
        echo "Invalid argument: $1"
    fi
    exit 1
fi



