package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

/*
 * Build script for open62541 single file C library.
 */

const LibraryDir = "./open62541/"

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func run() error {
	log.Printf("Starting build of open62541 static library...")

	src := filepath.Join(LibraryDir, "src", "open62541.c")
	obj := filepath.Join(LibraryDir, "lib", "open62541.o")
	ar := filepath.Join(LibraryDir, "lib", "libopen62541.a")
	includeDir := filepath.Join(LibraryDir, "include")

	log.Printf("Compiling %s to %s", src, obj)
	if err := compile(obj, src, includeDir); err != nil {
		return fmt.Errorf("compile: %w", err)
	}
	defer os.Remove(obj)

	log.Printf("Creating static library %s from %s", ar, obj)
	if err := archive(ar, obj); err != nil {
		return fmt.Errorf("archive: %w", err)
	}

	log.Printf("Finished building static library %s", ar)
	return nil
}

func compile(obj, src, includeDir string) error {
	objDir := filepath.Dir(obj)
	if err := os.MkdirAll(objDir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", objDir, err)
	}

	cmd := exec.Command("gcc", "-c",
		"-I", includeDir,
		"-o", obj,
		src,
		"-lmbedtls", "-lmbedx509", "-lmbedcrypto")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running gcc: %w", err)
	}
	log.Printf("Compiled %s to %s", src, obj)
	return nil
}

func archive(ar, obj string) error {
	cmd := exec.Command("ar", "rcs", ar, obj)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running ar: %w", err)
	}
	log.Printf("Archived %s to %s", obj, ar)
	return nil
}
