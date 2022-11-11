#!/bin/sh

rm cutron cutron.exe
rm *.zip

# macOS x86 binaries
env GOOS=darwin GOARCH=amd64 go build ./cmd/cutron/
zip calcutron-x86_64-macos.zip cutron

# macOS ARM binaries
env GOOS=darwin GOARCH=arm64  go build ./cmd/cutron/
zip calcutron-aarch64-macos.zip cutron

# Linux x86 binaries
env GOOS=linux GOARCH=amd64 go build ./cmd/cutron/
zip calcutron-x86_64-linux.zip cutron

# Linux ARM binaries
env GOOS=linux GOARCH=arm64 go build ./cmd/cutron/
zip calcutron-aarch64-linux.zip cutron

# Windows x86 binaries
env GOOS=windows GOARCH=amd64 go build ./cmd/cutron/
zip calcutron-x86_64-windows.zip cutron.exe

# Windows ARM binaries
env GOOS=windows GOARCH=arm64 go build ./cmd/cutron/
zip calcutron-aarch64-windows.zip cutron.exe
