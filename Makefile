#! /usr/bin/make
#
# Makefile for Model
#
# Targets:
# - "depend" retrieves the Go packages needed to run the linter and tests
# - "lint" runs the linter and checks the code format using goimports
# - "test" runs the tests
# - "release" creates a new release commit, tags the commit and pushes the tag to GitHub.
#
# Meta targets:
# - "all" is the default target, it runs "lint" and "test"
# - "ci" runs "depend" and "all"
#
MAJOR=1
MINOR=7
BUILD=8

GO_FILES=$(shell find . -type f -name '*.go')

# Only list test and build dependencies
# Standard dependencies are installed via go get
DEPEND=\
	golang.org/x/tools/cmd/goimports \
	honnef.co/go/tools/cmd/staticcheck \
	github.com/mjibson/esc

all: lint check-generated test

ci: depend all

depend:
	@go mod download
	@go get -v $(DEPEND)

generate:
	go generate ./cmd/mdl/

lint:
ifneq ($(GOOS),windows)
	@if [ "`goimports -l $(GO_FILES) | tee /dev/stderr`" ]; then \
		echo "^ - Repo contains improperly formatted go files" && echo && exit 1; \
	fi
	@if [ "`staticcheck ./... | tee /dev/stderr`" ]; then \
		echo "^ - staticcheck errors!" && echo && exit 1; \
	fi
endif

check-generated: generate
	@if ! git diff -s --exit-code cmd/mdl/webapp.go; then \
  		echo 'cmd/mdl/webapp.go is different, run `make generate` before commit!'; \
  	fi

test:
	env GO111MODULE=on go test ./...

release:
# First make sure all is clean
	@git diff-index --quiet HEAD
	@go mod tidy --compat=1.17

# Bump version number
	@sed 's/Major = .*/Major = $(MAJOR)/' pkg/version.go > _tmp && mv _tmp pkg/version.go
	@sed 's/Minor = .*/Minor = $(MINOR)/' pkg/version.go > _tmp && mv _tmp pkg/version.go
	@sed 's/Build = .*/Build = $(BUILD)/' pkg/version.go > _tmp && mv _tmp pkg/version.go
	@sed 's/badge\/Version-.*/badge\/Version-v$(MAJOR).$(MINOR).$(BUILD)-blue.svg)/' README.md > _tmp && mv _tmp README.md
	@sed 's/model@v.*\/\(.*\)tab=doc/model@v$(MAJOR).$(MINOR).$(BUILD)\/\1tab=doc/' README.md > _tmp && mv _tmp README.md
	@sed 's/model@v.*\/\(.*\)tab=doc/model@v$(MAJOR).$(MINOR).$(BUILD)\/\1tab=doc/' DSL.md > _tmp && mv _tmp DSL.md

# Make sure mdl and stz build
	@cd cmd/mdl && go install
	@cd cmd/stz && go install

# Commit and push
	@git add .
	@git commit -m "Release v$(MAJOR).$(MINOR).$(BUILD)"
	@git tag v$(MAJOR).$(MINOR).$(BUILD)
	@git push origin master
	@git push origin v$(MAJOR).$(MINOR).$(BUILD)
