#!/usr/bin/env make
#
# Makefile to script build and deployment automated targets
#

# External tools - defined as variables so we can optionally
# overwrite them if they are not in the $PATH.

GO = go 

all: clean fmt test 

install: install-libs

install-libs:
	$(GO) install ./...

clean:
	$(GO) clean ./...

fmt:
	$(GO) list ./... | while read pkg ; do $(GO) fmt $$pkg || exit 1 ; done

test:
	$(GO) list ./... | while read pkg ; do $(GO) test $$pkg ; done
