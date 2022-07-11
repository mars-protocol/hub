#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

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
  build_tags += gcc
else ifeq (rocksdb,$(findstring rocksdb,$(MARS_BUILD_OPTIONS)))
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# ********** process linker flags **********

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=osmosis \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=osmosisd \
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

BUILD_FLAGS := -ldflags '$(ldflags)'

all: proto-gen lint test install

###############################################################################
###                                  Build                                  ###
###############################################################################

install:
	@echo "🤖 Installing marsd..."
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/marsd
	@echo "✅ Completed installation!"

build:
	@echo "🤖 Building marsd..."
	go build $(BUILD_FLAGS) -o build/bin/marsd ./cmd/marsd
	@echo "✅ Completed build!"

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test:
	@echo "🤖 Running tests..."
	go test -mod=readonly ./x/...
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

###############################################################################
###                                Linting                                  ###
###############################################################################

golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

lint:
	@echo "🤖 Running linter..."
	go run $(golangci_lint_cmd) run --timeout=10m
	@echo "✅ Completed tests!"
