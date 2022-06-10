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
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
# go get github.com/regen-network/cosmos-proto@latest # doesn't work in install mode

set -eo pipefail

protoc_install_gocosmos() {
  echo "Installing protobuf gocosmos plugin"
  # we should use go install, but regen-network/cosmos-proto contains
  # replace directives. It must not contain directives that would cause
  # it to be interpreted differently than if it were the main module.
  # So the command below issues a warning and we are muting it for now.
  #
  # Installing plugins must be done outside of the module
  (go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@v0.3.1 2> /dev/null)
}

# NOTE: you can comment this line out if gocosmos plugin is already installed
#protoc_install_gocosmos

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
if [ -d "./github.com/mars-protocol/mars" ]; then
  cp -r github.com/mars-protocol/mars/* ./
  rm -rf github.com
fi

go mod tidy
