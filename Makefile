EXECUTABLE=bin/mavadsb
VERSION=$(shell git rev-parse --abbrev-ref HEAD)-$(shell git describe --tags --always --long --dirty)
HOSTNAME=$(shell hostname)
TIMESTAMP=$(shell date)
USER=$(shell id -u -n)
LDFLAGS=-ldflags='-s -w -X main.COMPILE_VERSION=$(VERSION) -X main.COMPILE_HOSTNAME=$(HOSTNAME) -X "main.COMPILE_TIMESTAMP=$(TIMESTAMP)" -X "main.COMPILE_USER=$(USER)"'
INPUT_FILES=cmd/*.go

# Architectures
## Linux
LINUX_AMD64=$(EXECUTABLE)_linux_amd64-$(VERSION)
LINUX_ARM64=$(EXECUTABLE)_linux_arm64-$(VERSION)
LINUX_ARMV7=$(EXECUTABLE)_linux_armv7-$(VERSION)

## Darwin
DARWIN_AMD64=$(EXECUTABLE)_darwin_amd64-$(VERSION)
DARWIN_ARM64=$(EXECUTABLE)_darwin_arm64-$(VERSION)

## BSDs
FREEBSD_AMD64=$(EXECUTABLE)_freebsd_amd64-$(VERSION)
OPENBSD_AMD64=$(EXECUTABLE)_openbsd_amd64-$(VERSION)

# Windows
WINDOWS_AMD64=$(EXECUTABLE)_windows_amd64-$(VERSION).exe

all: clean linux darwin freebsd windows
.PHONY: all clean upload

clean:
	rm bin/*; true

deps:
	go get ./...

ci: clean deps

upload:
	mc cp bin/* minio/private/mavadsb/$(VERSION)/
	mc share download --expire=72h minio/private/mavadsb/$(VERSION)/

linux: linux-amd64 linux-arm64 linux-armv7
linux-amd64: $(LINUX_AMD64)
linux-arm64: $(LINUX_ARM64)
linux-armv7: $(LINUX_ARMV7)

darwin: darwin-amd64 darwin-arm64
darwin-amd64: $(DARWIN_AMD64)
darwin-arm64: $(DARWIN_ARM64)

freebsd: $(FREEBSD_AMD64) 
openbsd: $(OPENBSD_AMD64) 

windows: $(WINDOWS_AMD64)

$(LINUX_ARMV7): $(INPUT_FILES)
	env GOOS=linux GOARCH=arm GOARM=7 go build -o $@ $(LDFLAGS) $^
$(LINUX_ARM64): $(INPUT_FILES)
	env GOOS=linux GOARCH=arm64 go build -o $@ $(LDFLAGS) $^
$(LINUX_AMD64): $(INPUT_FILES)
	env GOOS=linux GOARCH=amd64 go build -o $@ $(LDFLAGS) $^

$(DARWIN_AMD64): $(INPUT_FILES)
	env GOOS=darwin GOARCH=amd64 go build -o $@ $(LDFLAGS) $^
$(DARWIN_ARM64): $(INPUT_FILES)
	env GOOS=darwin GOARCH=arm64 go build -o $@ $(LDFLAGS) $^


$(FREEBSD_AMD64): $(INPUT_FILES)
	env GOOS=freebsd GOARCH=amd64 go build -o $@ $(LDFLAGS) $^
$(OPENBSD_AMD64): $(INPUT_FILES)
	env GOOS=openbsd GOARCH=amd64 go build -o $@ $(LDFLAGS) $^

$(WINDOWS_AMD64): $(INPUT_FILES)
	env GOOS=windows GOARCH=amd64 go build -o $@ $(LDFLAGS) $^
