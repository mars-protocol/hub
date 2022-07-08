#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=mars \
	-X github.com/cosmos/cosmos-sdk/version.AppName=marsd \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

all: proto-gen test install

###############################################################################
###                                  Build                                  ###
###############################################################################

install:
	@echo "🤖 Installing marsd..."
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/marsd
	@echo "✅ Completed installation!"

build:
	@echo "🤖 Building marsd..."
	go build $(BUILD_FLAGS) -o bin/marsd ./cmd/marsd
	@echo "✅ Completed build!"

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test:
	@echo "🤖 Running tests..."
	go test -mod=readonly ./...
	@echo "✅ Completed tests!"

###############################################################################
###                                Protobuf                                 ###
###############################################################################

proto-gen: proto-go-gen proto-swagger-gen

proto-go-gen:
	@echo "🤖 Generating Go code from protobuf..."
	bash ./scripts/protocgen.sh
	@echo "✅ Completed Go code generation!"

proto-swagger-gen:
	@echo "🤖 Generating Swagger code from protobuf..."
	bash ./scripts/protoc-swagger-gen.sh
	@echo "✅ Completed Swagger code generation!"
