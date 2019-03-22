GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test

DEPCMD=dep

all: build test
build:
			$(GOBUILD)
test:
			$(GOTEST) -v ./... -cover	
deps:
			$(DEPCMD) ensure

