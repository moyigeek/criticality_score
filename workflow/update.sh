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
    ensure_exists github_updated.src
}

LOCK_FILE="/var/lock/csflow.lock"

run() {
    touch "$1_updated.src"

    # check if lock exists
    if [ -f "$LOCK_FILE" ]; then
        echo "Lock file exists, another process is running, or the previous process was killed."
        echo "Please remove the lock file: $LOCK_FILE"
        echo "* Update by triggering $1"
        echo "** Target: $1_updated"
        echo "** Time: $(date '+%Y-%m-%d %H:%M:%S (%Z)')"
        echo "** Status: Failed (Lock file exists)"
        exit 1
    fi

    # create lock file
    touch "$LOCK_FILE"

    # when the script is interrupted, remove the lock file
    trap 'rm -f "$LOCK_FILE"; exit 1' INT TERM EXIT

    echo "* Update by triggering $1"
    echo "** Target: $1_updated"
    echo "** Time: $(date '+%Y-%m-%d %H:%M:%S (%Z)')"
    echo 
    echo "=================START================="

    make -s -e -f "$SCRIPT_DIR/flow.mk"

    echo "==================END=================="
    echo "** Finish: $(date '+%Y-%m-%d %H:%M:%S (%Z)')"
    echo
}

# getopts -C : change dir

while getopts "C:" opt; do
    case $opt in
        C)
            echo "Changing directory to $OPTARG"
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

if [ "$1" = "package" ] || [ "$1" = "git" ] || 
    [ "$1" = "depsdev" ] || [ "$1" = "gitlink" ] ||
    [ "$1" = "github" ]; then
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



