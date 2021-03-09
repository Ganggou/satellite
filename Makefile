export GOSUMDB=off
VERSION := 1.0.0
PROJECTNAME := satellite
LIBNAME=libsp
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
LDFLAGS=-ldflags "$(RAW_LDFLAGS)"
GOBUILD=go build $(LDFLAGS)
GOOS ?= $(shell go env GOOS)

linux:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(GOBIN)/$(PROJECTNAME)-$(VERSION)-$(GOOS)

darwin:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(GOBIN)/$(PROJECTNAME)-$(VERSION)-$(GOOS)
