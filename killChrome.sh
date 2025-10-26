#!/bin/bash
if test -L "$(pwd)/profile/SingletonLock"; then
    SingletonLock="$(basename $(readlink -f "$(pwd)/profile/SingletonLock"))"
    pid="${SingletonLock/"$(hostname)-"/}"
    echo "Found ${SingletonLock} killing pid ${pid}"
    if kill -0 "${pid}" 2>/dev/null; then
        kill -TERM "${pid}"
    else
        echo "${pid} not found"
    fi
else
    echo "$(pwd)/profile/SingletonLock not found" 
fi