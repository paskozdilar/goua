package main

/*
 * Download script for open62541 single file C library.
 */

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	log.Printf("Starting download of latest open62541 release...")
	r, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("getting latest release: %w", err)
	}
	log.Printf("Latest release: %s", r.TagName)

	for _, a := range r.Assets {
		var dir string
		switch a.Name {
		case "open62541.c":
			dir = filepath.Join(OutputDir, "src")
		case "open62541.h":
			dir = filepath.Join(OutputDir, "include")
		default:
			log.Printf("Skipping asset: %s", a.Name)
			continue
		}

		log.Printf("Downloading asset: %s", a.Name)
		if err := downloadAsset(dir, a); err != nil {
			return fmt.Errorf("downloading asset %s: %v", a.Name, err)
		}
	}
	log.Printf("Finished downloading assets for open62541 release %s", r.TagName)

	return nil
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

func downloadAsset(dir string, a Asset) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create output directory %s: %w", OutputDir, err)
	}

	resp, err := http.Get(a.BrowserDownloadURL)
	if err != nil {
		return fmt.Errorf("HTTP GET %s: %w", a.BrowserDownloadURL, err)
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath.Join(dir, a.Name))
	if err != nil {
		return fmt.Errorf("create file %s: %w", a.Name, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("write to file %s: %w", a.Name, err)
	}
	return nil
}
