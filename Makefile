GOCACHE ?= $(PWD)/.gocache
GOMODCACHE ?= $(PWD)/.gomodcache
GO_ENV := GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE)

export GOCACHE
export GOMODCACHE

include ./makefiles/dev.mk
include ./makefiles/test.mk
include ./makefiles/lint.mk
include ./makefiles/client.mk
