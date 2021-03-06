-include config.mk
.SECONDEXPANSION:

BUILDDIR ?= $(CURDIR)

BINDIR      ?= $(BUILDDIR)/bin
ALLBIN ?= $(notdir $(shell find cmd -mindepth 1 -type d -print))

# go option
PKG        := ./...
TAGS       :=
TESTS      := .
TESTFLAGS  :=
LDFLAGS    := -w -s
GOFLAGS    := -gcflags=all='-N'

TAGS := $(shell git describe --tags 2>/dev/null)
ifneq ($(TAGS),)
VERSION_OPTIONS := -X "github.com/gueckmooh/bs/pkg/version.version_hash=$(TAGS)"
endif
COMMIT_OPTIONS := -X "github.com/gueckmooh/bs/pkg/version.commit_hash=$(shell git describe --always)"
BUILD_TIME_OPTION := -X "github.com/gueckmooh/bs/pkg/version.build_time=$(shell date)"

LDFLAGS += $(VERSION_OPTIONS) $(COMMIT_OPTIONS) $(BUILD_TIME_OPTION)

DEPDIR ?= .deps
.PRECIOUS: %/.f $(DEPDIR)/%.d

%/.f:
	$(QUIET)mkdir -p $(dir $@)
	$(QUIET)touch $@

NOINC = clean, mrproper

SRC := $(shell find pkg -type f -name '*.go' -print) $(shell find cmd -type f -name '*.go' -print) go.mod
GENERATED_SRC := pkg/lua/luabslib/cppprofile_gen.go \
                 pkg/lua/luabslib/profile_gen.go \
                 pkg/lua/luabslib/component_gen.go \
                 pkg/lua/luabslib/components_gen.go \
                 pkg/lua/luabslib/project_gen.go \
                 pkg/lua/luabslib/git_repository_gen.go
SRC += $(GENERATED_SRC)
.PRECIOUS: $(GENERATED_SRC)

ALLBINS := $(addprefix $(BINDIR)/, $(ALLBIN))

.PHONY: all
all: build

# Required for globs to work correctly
SHELL      = /usr/bin/env bash

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(ALLBINS)

pkg/lua/luabslib/%_gen.go: pkg/lua/luabslib/%.go $(wildcard pkg/lua/luabslib/gen/*.go)
	go generate $<

$(BINDIR)/%: $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./cmd/$(notdir $@)

.PHONY: test
test: unit-test integ-test

TEST_SRC := $(shell find ./pkg -type f -name '*_test.go' -print | sed 's|/[^/]*$$||g' | uniq)
.PHONY: unit-test
unit-test:
	go test -v $(TEST_SRC)

.PHONY: integ-test
integ-test: build
	python3 tests/runtest.py
