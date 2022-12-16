#!/usr/bin/env bash

#== Requirements ==
#
## make sure your `go env GOPATH` is in the `$PATH`
## Install:
## + latest buf (v1.0.0-rc11 or later)
## + protobuf v3
#
## All protoc dependencies must be installed not in the module scope
## currently we must use grpc-gateway v1
# cd ~
# go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.16.0
# go get github.com/regen-network/cosmos-proto@latest # doesn't work in install mode

set -eo pipefail

# The directory where the final output are to be stored
docs_dir="./docs"

# The directory where temporary swagger files are to be stored before they are
# combined. Will be deleted in the end
tmp_dir="./tmp-swagger-gen"
if [ -d $tmp_dir ]; then
  rm -rf $tmp_dir
fi
mkdir -p $tmp_dir

# Directories that contain protobuf files that are to be transpiled into swagger
# These include Mars modules and third party modules and services
cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk)
ibc_go_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/ibc-go/v4)
wasmd_dir=$(go list -f '{{ .Dir }}' -m github.com/CosmWasm/wasmd)
proto_dirs=$(find ./proto ${cosmos_sdk_dir}/proto ${ibc_go_dir}/proto ${wasmd_dir}/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

# Generate swagger files for `query.proto` and `service.proto`
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \)); do
    echo $file
    buf generate --template ./proto/buf.gen.swagger.yaml $file
  done
done

# Combine swagger files
# Uses nodejs package `swagger-combine`.
# All the individual swagger files need to be configured in `config.json` for merging
echo "Combining swagger files..."
swagger-combine ${docs_dir}/config.json \
  -o ${docs_dir}/swagger.yml \
  -f yaml \
  --continueOnConflictingPaths true \
  --includeDefinitions true

# Clean swagger files
rm -rf $tmp_dir
