#! /usr/bin/make
#
# Makefile for goa v3 docs plugin
#
# Targets:
# - "gen" generates the goa files for the example services

# include common Makefile content for plugins
include $(GOPATH)/src/goa.design/plugins/plugins.mk

gen:
	goa gen goa.design/plugins/v3/docs/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/docs/examples/calc" && \
	make example

example:
	@ rm -rf "$(GOPATH)/src/goa.design/plugins/docs/examples/calc/cmd" && \
	goa example goa.design/plugins/v3/docs/examples/calc/design -o "$(GOPATH)/src/goa.design/plugins/docs/examples/calc"

build-examples:
	@cd "$(GOPATH)/src/goa.design/plugins/docs/examples/calc" && \
		go build ./cmd/calc && go build ./cmd/calc-cli

clean:
	@cd "$(GOPATH)/src/goa.design/plugins/docs/examples/calc" && \
		rm -f calc calc-cli
