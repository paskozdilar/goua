package main

/*
 * Download script for open62541 single file C library.
 */

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	ReleasesURL = "https://api.github.com/repos/open62541/open62541/releases/latest"
	OutputDir   = "./open62541/"
)

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	r, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("getting latest release: %w", err)
	}
	log.Printf("Latest release: %s", r.TagName)

	if alreadyDownloaded(r.TagName) {
		log.Printf("Release %s is already downloaded, skipping download.", r.TagName)
		return nil
	}

	log.Printf("Starting download of latest open62541 release...")
	for _, a := range r.Assets {
		if a.Name != "open62541.h" && a.Name != "open62541.c" {
			log.Printf("Skipping asset: %s", a.Name)
			continue
		}

		log.Printf("Downloading asset: %s", a.Name)
		if err := downloadAsset(a); err != nil {
			return fmt.Errorf("downloading asset %s: %v", a.Name, err)
		}
	}
	log.Printf("Finished downloading assets for open62541 release %s", r.TagName)

	return nil
}

func alreadyDownloaded(tag string) bool {
	for _, file := range []string{"open62541.h", "open62541.c"} {
		f, err := os.Open(file)
		if err != nil {
			return false // File does not exist, not downloaded
		}
		defer f.Close()

		s := bufio.NewScanner(f)
		prefix := " * Git-Revision:"
		for s.Scan() {
			line := s.Text()
			if strings.HasPrefix(line, prefix) {
				version := strings.TrimSpace(strings.TrimPrefix(line, prefix))
				if version != tag {
					return false // Installed version does not match the latest release
				}
				break
			}
		}
	}
	return true
}

func getLatestRelease() (*Release, error) {
	resp, err := http.Get(ReleasesURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET %s: %w", ReleasesURL, err)
	}
	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)
	r := &Release{}
	if err := d.Decode(r); err != nil {
		return nil, fmt.Errorf("decoding JSON response: %w", err)
	}
	if err := d.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("decoding JSON response: expected EOF")
	}
	return r, nil
}

func downloadAsset(a Asset) error {
	resp, err := http.Get(a.BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("HTTP GET %s: %w", a.BrowserDownloadURL, err)
	}
	defer resp.Body.Close()

	file, err := os.Create(a.Name)
	if err != nil {
		return fmt.Errorf("create file %s: %w", a.Name, err)
	}
	defer file.Close()

	// If .c file, add go:build tag to prevent compilation by gopls
	if strings.HasSuffix(a.Name, ".c") {
		if _, err := file.WriteString("//go:build ignore\n"); err != nil {
			return fmt.Errorf("write build tag to file %s: %w", a.Name, err)
		}
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("write to file %s: %w", a.Name, err)
	}
	return nil
}
