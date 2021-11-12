package downloader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getReleaseHash(t *testing.T) {
	is := assert.New(t)

	tests := []struct {
		version string
		wantSHA string
	}{
		{
			version: "1.10.11",
			wantSHA: "7231b3efb8095d3dd18d7164c3fa84d7705759d3",
		},
		{
			version: "1.10.9",
			wantSHA: "eae3b1946a276ac099e0018fc792d9e8c3bfda6d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			gotSHA, err := getReleaseHash(tt.version)
			is.NoError(err)

			is.Equal(tt.wantSHA, gotSHA)
		})
	}
}

func Test_GethBinaryURL(t *testing.T) {
	is := assert.New(t)

	fmt.Println(runtime.GOARCH)

	tests := []struct {
		os, arch, version, hash string
		wantURL                 string
	}{
		{
			os:      "linux",
			arch:    "amd64",
			version: "1.10.11",
			hash:    "7231b3efb8095d3dd18d7164c3fa84d7705759d3",
			wantURL: "https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-1.10.11-7231b3ef.tar.gz",
		},
		{
			os:      "darwin",
			arch:    "amd64",
			version: "1.10.11",
			hash:    "7231b3efb8095d3dd18d7164c3fa84d7705759d3",
			wantURL: "https://gethstore.blob.core.windows.net/builds/geth-darwin-amd64-1.10.11-7231b3ef.tar.gz",
		},
		{
			os:      "darwin",
			arch:    "arm64",
			version: "1.10.11",
			hash:    "7231b3efb8095d3dd18d7164c3fa84d7705759d3",
			wantURL: "https://gethstore.blob.core.windows.net/builds/geth-darwin-amd64-1.10.11-7231b3ef.tar.gz",
		},
		{
			os:      "windows",
			arch:    "amd64",
			version: "1.10.11",
			hash:    "7231b3efb8095d3dd18d7164c3fa84d7705759d3",
			wantURL: "https://gethstore.blob.core.windows.net/builds/geth-windows-amd64-1.10.11-7231b3ef.zip",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.os+"-"+tt.arch+"-"+tt.version, func(t *testing.T) {
			t.Parallel()

			gotURL, exists, error := GethBinaryURL(tt.os, tt.arch, tt.version, tt.hash)
			is.NoError(error)
			is.True(exists)

			is.Equal(tt.wantURL, gotURL)
		})
	}
}

// Test_DownloadFile tests that the geth archive is downloaded to the correct location
// note: long-running test ~30s
func Test_DownloadFile(t *testing.T) {
	is := assert.New(t)

	tempDir, err := os.MkdirTemp("", "download-file")
	is.NoError(err)
	defer os.RemoveAll(tempDir)

	f, err := DownloadFile("https://gethstore.blob.core.windows.net/builds/geth-darwin-amd64-1.10.11-7231b3ef.tar.gz", tempDir)
	is.NoError(err)

	is.FileExists(f.Name())
}

// Test_ExtractGethArchive tests that the geth archive is extracted to the correct location and ready for use
// note: long-running test ~40s
func Test_ExtractGethArchive(t *testing.T) {
	is := assert.New(t)

	tests := []struct {
		url, os, version string
	}{
		{
			url:     "https://gethstore.blob.core.windows.net/builds/geth-darwin-amd64-1.10.11-7231b3ef.tar.gz",
			os:      "darwin",
			version: "1.10.11",
		},
		{
			url:     "https://gethstore.blob.core.windows.net/builds/geth-windows-amd64-1.10.11-7231b3ef.zip",
			os:      "windows",
			version: "1.10.11",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.os+"-"+tt.version, func(t *testing.T) {
			t.Parallel()

			tempDir, err := os.MkdirTemp("", "download-file")
			is.NoError(err)
			defer os.RemoveAll(tempDir)

			f, err := DownloadFile(tt.url, tempDir)
			is.NoError(err)

			err = ExtractGethArchive(tt.os, f.Name(), tempDir)
			is.NoError(err)

			gethBinary := filepath.Join(tempDir, GethBinaryFilename(tt.os, tt.version))
			is.FileExists(gethBinary)
		})
	}
}
