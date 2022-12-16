#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
BUILDDIR ?= $(CURDIR)/build
LEDGER_ENABLED ?= true

# ********** Golang configs **********

export GO111MODULE = on

GO_MAJOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)

# ********** process build tags **********

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq (cleveldb,$(findstring cleveldb,$(MARS_BUILD_OPTIONS)))
  build_tags += gcc cleveldb
else ifeq (rocksdb,$(findstring rocksdb,$(MARS_BUILD_OPTIONS)))
  build_tags += gcc rocksdb
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# ********** process linker flags **********

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=mars \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=marsd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq (cleveldb,$(findstring cleveldb,$(MARS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
else ifeq (rocksdb,$(findstring rocksdb,$(MARS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb
endif
ifeq (,$(findstring nostrip,$(MARS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ifeq ($(LINK_STATICALLY),true)
	ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags '$(build_tags)' -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(MARS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

all: proto-gen lint test install

###############################################################################
###                                  Build                                  ###
###############################################################################

install: enforce-go-version
	@echo "🤖 Installing marsd..."
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/marsd
	@echo "✅ Completed installation!"

build: enforce-go-version
	@echo "🤖 Building marsd..."
	@go build $(BUILD_FLAGS) -o $(BUILDDIR)/ ./cmd/marsd
	@echo "✅ Completed build!"

# Make sure that Go version is 1.19
#
# From Osmosis discord:
# https://discord.com/channels/798583171548840026/837144686387920936/1049449765240315925
#
# > Valardragon - 12/05/2022 10:18 PM
# > It was just pointed out from `@jhernandez | stargaze.zone`, that the choice
#   of golang version between go 1.18 and go 1.19 is consensus critical.
#   With insufficient info, this preliminarily seems due to go 1.19 changing the
#   memory model format, and something state-affecting in cosmwasm getting altered.
#   https://github.com/persistenceOne/incident-reports/blob/main/06-nov-2022_V4_upgrade_halt.md
enforce-go-version:
	@echo "🤖 Go version: $(GO_MAJOR_VERSION).$(GO_MINOR_VERSION)"
ifneq ($(GO_MINOR_VERSION),19)
	@echo "❌ ERROR: Go version 1.19 is required"
	@exit 1
endif

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test:
	@echo "🤖 Running tests..."
	@go test -mod=readonly ./x/...
	@echo "✅ Completed tests!"

###############################################################################
###                                Protobuf                                 ###
###############################################################################

# We use osmolabs' docker image instead of tendermintdev/sdk-proto-gen.
# The only difference is that the Osmosis version uses Go 1.19 while the
# tendermintdev one uses 1.18.
protoVer=v0.8
protoImageName=osmolabs/osmo-proto-gen:$(protoVer)
containerProtoGenGo=mars-proto-gen-go-$(protoVer)
containerProtoGenSwagger=mars-proto-gen-swagger-$(protoVer)

proto-gen: proto-go-gen proto-swagger-gen

proto-go-gen:
	@echo "🤖 Generating Go code from protobuf..."
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGenGo}$$"; then docker start -a $(containerProtoGenGo); else docker run --name $(containerProtoGenGo) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./scripts/protocgen.sh; fi
	@echo "✅ Completed Go code generation!"

proto-swagger-gen:
	@echo "🤖 Generating Swagger code from protobuf..."
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGenSwagger}$$"; then docker start -a $(containerProtoGenSwagger); else docker run --name $(containerProtoGenSwagger) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./scripts/protoc-swagger-gen.sh; fi
	@echo "✅ Completed Swagger code generation!"

###############################################################################
###                                Linting                                  ###
###############################################################################

golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	@echo "🤖 Running linter..."
	@go run $(golangci_lint_cmd) run --timeout=10m
	@echo "✅ Completed tests!"
