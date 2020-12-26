# Get current directory
current_dir := $(shell pwd)

# Get current commit hash
# commit_hash := $(shell git rev-parse --short=7 HEAD)

# Targets
.PHONY: test

all: build

install-deps:
	npm --prefix $(current_dir)/ui install

test-frontend: install-deps
	npm --prefix $(current_dir)/ui test

test-backend:		
	go clean -testcache
	go test -v -race github.com/istherepie/push-notifications/eventbroker
	go test -v github.com/istherepie/push-notifications/cmd/notification-server/webserver
	go test -v github.com/istherepie/push-notifications/metrics

test: test-frontend test-backend
	@echo "Running all tests"

build:
	@echo "Building binaries"

	mkdir $(current_dir)/build
	go build -o $(current_dir)/build/notification-server $(current_dir)/cmd/notification-server/notification-server.go

clean:
	@echo "Cleaning up..."
	rm -rf $(current_dir)/build