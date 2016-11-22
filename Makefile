VERSION = 0.3
VERSION_INFO_FILE = version-info.go

REPO_BASE = github.com/tintoy/docker-machine-driver-terraform

default: fmt build test

fmt:
	go fmt github.com/tintoy/docker-machine-driver-terraform/...

# Peform a development (current-platform-only) build.
dev: version fmt
	go build -o _bin/docker-machine-driver-terraform

install: dev
	go install

# Perform a full (all-platforms) build.
build: version build-windows64 build-linux64 build-mac64

build-windows64:
	GOOS=windows GOARCH=amd64 go build -o _bin/windows-amd64/docker-machine-driver-terraform.exe

build-linux64:
	GOOS=linux GOARCH=amd64 go build -o _bin/linux-amd64/docker-machine-driver-terraform

build-mac64:
	GOOS=darwin GOARCH=amd64 go build -o _bin/darwin-amd64/docker-machine-driver-terraform

# Produce archives for a GitHub release.
dist: build
	zip -9 _bin/windows-amd64.zip _bin/windows-amd64/docker-machine-driver-terraform.exe
	zip -9 _bin/linux-amd64.zip _bin/linux-amd64/docker-machine-driver-terraform
	zip -9 _bin/darwin-amd64.zip _bin/darwin-amd64/docker-machine-driver-terraform

test: fmt
	go test -v $(REPO_BASE) $(REPO_BASE)/fetch $(REPO_BASE)/terraform

version: $(VERSION_INFO_FILE)

$(VERSION_INFO_FILE): Makefile
	@echo "Update version info: v$(VERSION)."
	@echo "package main\n\n// DriverVersion is the current version of the Terraform driver for Docker Machine.\nconst DriverVersion = \"v$(VERSION) (`git rev-parse HEAD`)\"" > $(VERSION_INFO_FILE)
