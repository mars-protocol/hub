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

tmp_dir="./tmp-swagger-gen"
docs_dir="./docs"
proto_dirs=$(find ./proto ./third_party/proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

mkdir -p $tmp_dir

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
