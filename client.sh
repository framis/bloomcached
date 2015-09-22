#!/bin/bash

while IFS='' read -r line || [[ -n "$line" ]]; do
    echo -n "$line" | nc localhost 3333
    echo ""
done < "$1"