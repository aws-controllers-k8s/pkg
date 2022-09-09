SHELL := /bin/bash # Use bash syntax

.PHONY: all test

all: test

test: 			## Run code tests
	go test ${GO_TAGS} ./...

help:           ## Show this help.
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep | sed -e 's/\\$$//' \
		| awk -F'[:#]' '{print $$1 = sprintf("%-30s", $$1), $$4}'
