#!/bin/sh
#
# This script is intended to be run inside the osmolabs/osmo-proto-gen:v0.8
# docker container: https://hub.docker.com/r/osmolabs/osmo-proto-gen

set -eo pipefail

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep "option go_package" $file &> /dev/null ; then
      echo $file
      buf generate --template ./proto/buf.gen.go.yaml $file
    fi
  done
done

# move proto files to the right places
if [ -d "./github.com/mars-protocol/hub" ]; then
  cp -r github.com/mars-protocol/hub/* ./
  rm -rf github.com
fi

go mod tidy
