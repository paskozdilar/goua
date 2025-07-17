package goua

/*
 * Run `go generate` to download the latest version of the open62541 library
 * and build the static library.
 */

//go:generate go run ./tools/download/main.go
//go:generate go run ./tools/build/main.go
