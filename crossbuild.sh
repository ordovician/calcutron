#!/bin/sh

rm assemble simulate disassemble debugger
rm *.zip
rm *.exe

# macOS x86 binaries
env GOOS=darwin GOARCH=amd64  go build ./cmd/disassemble
env GOOS=darwin GOARCH=amd64  go build ./cmd/assemble/
env GOOS=darwin GOARCH=amd64  go build ./cmd/simulate/
env GOOS=darwin GOARCH=amd64  go build ./cmd/debugger
zip calcutron-x86_64-macos.zip assemble simulate disassemble debugger

# macOS ARM binaries
env GOOS=darwin GOARCH=arm64  go build ./cmd/disassemble
env GOOS=darwin GOARCH=arm64  go build ./cmd/assemble/
env GOOS=darwin GOARCH=arm64  go build ./cmd/simulate/
env GOOS=darwin GOARCH=arm64  go build ./cmd/debugger
zip calcutron-aarch64-macos.zip assemble simulate disassemble debugger

# Linux x86 binaries
env GOOS=linux GOARCH=amd64 go build ./cmd/disassemble
env GOOS=linux GOARCH=amd64 go build ./cmd/assemble/
env GOOS=linux GOARCH=amd64 go build ./cmd/simulate/
env GOOS=linux GOARCH=amd64 go build ./cmd/debugger
zip calcutron-x86_64-linux.zip assemble simulate disassemble debugger

# Linux ARM binaries
env GOOS=linux GOARCH=arm64 go build ./cmd/disassemble
env GOOS=linux GOARCH=arm64 go build ./cmd/assemble/
env GOOS=linux GOARCH=arm64 go build ./cmd/simulate/
env GOOS=linux GOARCH=arm64 go build ./cmd/debugger
zip calcutron-aarch64-linux.zip assemble simulate disassemble debugger

# Windows x86 binaries
env GOOS=windows GOARCH=amd64 go build ./cmd/disassemble
env GOOS=windows GOARCH=amd64 go build ./cmd/assemble/
env GOOS=windows GOARCH=amd64 go build ./cmd/simulate/
env GOOS=windows GOARCH=amd64 go build ./cmd/debugger
zip calcutron-x86_64-windows.zip *.exe

# Windows ARM binaries
env GOOS=windows GOARCH=arm64 go build ./cmd/disassemble
env GOOS=windows GOARCH=arm64 go build ./cmd/assemble/
env GOOS=windows GOARCH=arm64 go build ./cmd/simulate/
env GOOS=windows GOARCH=arm64 go build ./cmd/debugger
zip calcutron-aarch64-windows.zip *.exe