
.DEFAULT_GOAL := help

## gen-file-from-template: Build a file from a template file
.PHONY: test
test: create-bin
	go test ./... -coverprofile=bin/coverage.txt -covermode=atomic

## copy-scripts: copies all scripts from the scripts folder to bin
.PHONY: copy-scripts
copy-scripts: 
	cp scripts/* bin

## create-bin: create the bin directory if it doesnt exist
.PHONY: create-bin
create-bin:
	if [ ! -d bin ]; then mkdir bin; fi

## all: 
.PHONY: all
all: create-bin

## coverage: Reports on the test coverage
.PHONY: coverage
coverage: test
	go tool cover -html=bin/coverage.txt

## help: type for getting this help
.PHONY: help
help: Makefile
	@echo 
	@echo " Choose a command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
