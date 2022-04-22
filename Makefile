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

DEPDIR ?= .deps
.PRECIOUS: %/.f $(DEPDIR)/%.d

%/.f:
	$(QUIET)mkdir -p $(dir $@)
	$(QUIET)touch $@

NOINC = clean, mrproper

SRC := $(shell find pkg -type f -name '*.go' -print) $(shell find cmd -type f -name '*.go' -print) go.mod
SRC += pkg/lua/luabslib/cppprofile_gen.go

ALLBINS := $(addprefix $(BINDIR)/, $(ALLBIN))

.PHONY: all
all: build

# Required for globs to work correctly
SHELL      = /usr/bin/env bash

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(ALLBINS)

pkg/lua/luabslib/cppprofile_gen.go: pkg/lua/luabslib/definitions/CPPProfile.xml $(wildcard pkg/lua/luabslib/gen/*.go)
	go generate ./pkg/lua/luabslib/

$(BINDIR)/%: $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./cmd/$(notdir $@)

.PHONY: test
test: build
	python3 tests/runtest.py
