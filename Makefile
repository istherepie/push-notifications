# Get current directory
current_dir := $(shell pwd)

# Get current commit hash
# commit_hash := $(shell git rev-parse --short=7 HEAD)

# Targets
.PHONY: test

all: testing clean build

test:
	@echo "Running all tests"
	go clean -testcache
	go test -v -race github.com/istherepie/push-notifications/eventbroker
	go test -v github.com/istherepie/push-notifications/cmd/notification-server/webserver
	go test -v github.com/istherepie/push-notifications/metrics

build:
	@echo "Building binaries"

	mkdir $(current_dir)/build

clean:
	@echo "Cleaning up..."
	rm -rf $(current_dir)/build