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
MINOR=11
BUILD=0

GO_FILES=$(shell find . -type f -name '*.go')

# React app source files and dependencies
WEBAPP_DIR=cmd/mdl/webapp
WEBAPP_SRC_FILES=$(shell find $(WEBAPP_DIR)/src -type f \( -name '*.tsx' -o -name '*.ts' -o -name '*.css' -o -name '*.html' \) 2>/dev/null || true)
WEBAPP_CONFIG_FILES=$(WEBAPP_DIR)/package.json $(WEBAPP_DIR)/tsconfig.json $(WEBAPP_DIR)/webpack.config.js $(WEBAPP_DIR)/webpack.config.base.js $(WEBAPP_DIR)/.babelrc.js
WEBAPP_BUILD_OUTPUT=$(WEBAPP_DIR)/dist/main.js

# Only list test and build dependencies
# Standard dependencies are installed via go get
DEPEND=\
	golang.org/x/tools/cmd/goimports@latest \
	github.com/golangci/golangci-lint/cmd/golangci-lint@latest \
	github.com/mjibson/esc@latest

all: lint test build

ci: depend all

depend:
	@go mod download
	@for package in $(DEPEND); do go install $$package; done

lint:
ifneq ($(GOOS),windows)
	@if [ "`goimports -l $(GO_FILES) | tee /dev/stderr`" ]; then \
		echo "^ - Repo contains improperly formatted go files" && echo && exit 1; \
	fi
	@output=$$(golangci-lint run ./... | grep -v "^0 issues\\.$$"); \
	if [ -n "$$output" ]; then \
		echo "$$output" && echo "^ - golangci-lint errors!" && echo && exit 1; \
	fi
endif

test:
	go test ./... --coverprofile=cover.out

# Ensure package-lock.json exists or is updated if package.json changes
$(WEBAPP_DIR)/package-lock.json: $(WEBAPP_DIR)/package.json
	@echo "Generating/updating package-lock.json..."
	@cd $(WEBAPP_DIR) && npm install --package-lock-only

# Install npm dependencies if package.json or package-lock.json changed
$(WEBAPP_DIR)/node_modules/.install-timestamp: $(WEBAPP_DIR)/package.json $(WEBAPP_DIR)/package-lock.json
	@echo "Installing npm dependencies..."
	@cd $(WEBAPP_DIR) && npm install
	@touch $(WEBAPP_DIR)/node_modules/.install-timestamp

# Build React app only if source files or config changed
# package-lock.json is not a direct dependency here because its generation is handled by the rule above,
# and npm install (triggered by .install-timestamp) will use it.
$(WEBAPP_BUILD_OUTPUT): $(WEBAPP_SRC_FILES) $(WEBAPP_CONFIG_FILES) $(WEBAPP_DIR)/node_modules/.install-timestamp
	@echo "Building React app..."
	@cd $(WEBAPP_DIR) && npm run build

# Phony target that depends on the actual build output
build-ui: $(WEBAPP_BUILD_OUTPUT)

# Force rebuild of UI (useful for development)
build-ui-force:
	@echo "Force building React app..."
	@cd $(WEBAPP_DIR) && npm install
	@cd $(WEBAPP_DIR) && npm run build

build: build-ui
	@cd cmd/mdl && go install
	@cd cmd/stz && go install

serve: build
	@cmd/mdl/mdl serve

# Clean build artifacts
clean-ui:
	@echo "Cleaning React build artifacts..."
	@rm -rf $(WEBAPP_DIR)/dist/*
	@rm -f $(WEBAPP_DIR)/node_modules/.install-timestamp

clean: clean-ui

release: build
# First make sure all is clean
	@git diff-index --quiet HEAD
	@go mod tidy

# Bump version number
	@sed 's/Major = .*/Major = $(MAJOR)/' pkg/version.go > _tmp && mv _tmp pkg/version.go
	@sed 's/Minor = .*/Minor = $(MINOR)/' pkg/version.go > _tmp && mv _tmp pkg/version.go
	@sed 's/Build = .*/Build = $(BUILD)/' pkg/version.go > _tmp && mv _tmp pkg/version.go
	@sed 's/badge\/Version-.*/badge\/Version-v$(MAJOR).$(MINOR).$(BUILD)-blue.svg)/' README.md > _tmp && mv _tmp README.md
	@sed 's/model[ @]v.*\/\(.*\)tab=doc/model@v$(MAJOR).$(MINOR).$(BUILD)\/\1tab=doc/' README.md > _tmp && mv _tmp README.md
	@sed 's/mdl v[0-9]*\.[0-9]*\.[0-9]*, editor started\./mdl v$(MAJOR).$(MINOR).$(BUILD), editor started./' README.md > _tmp && mv _tmp README.md
	@sed 's/model@v.*\/\(.*\)tab=doc/model@v$(MAJOR).$(MINOR).$(BUILD)\/\1tab=doc/' DSL.md > _tmp && mv _tmp DSL.md

# Commit and push
	@git add .
	@git commit -m "Release v$(MAJOR).$(MINOR).$(BUILD)"
	@git tag v$(MAJOR).$(MINOR).$(BUILD)
	@git push origin main
	@git push origin v$(MAJOR).$(MINOR).$(BUILD)

.PHONY: all ci depend lint test build-ui build-ui-force build serve release clean-ui clean
