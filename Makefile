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
.PRECIOUS: %/.f %.c $(DEPDIR)/%.d

%/.f:
	$(QUIET)mkdir -p $(dir $@)
	$(QUIET)touch $@

NOINC = clean, mrproper

SRC := $(shell find . -type f -name '*.go' -print) go.mod

ALLBINS := $(addprefix $(BINDIR)/, $(ALLBIN))

.PHONY: all
all: build

# Required for globs to work correctly
SHELL      = /usr/bin/env bash

# ------------------------------------------------------------------------------
#  build

.PHONY: build
build: $(ALLBINS)

$(BINDIR)/%: $(SRC)
	GO111MODULE=on go build $(GOFLAGS) -trimpath -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./cmd/$(notdir $@)
