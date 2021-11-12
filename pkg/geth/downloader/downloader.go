package downloader

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// DefaultGethBinaryDir is a subdir of the app dir - which contains the geth binaries
const DefaultGethBinaryDir = "bin"

var ErrUnsupportedArchive = errors.New("unsupported archive format")

// GethBinaryFilename returns the path to the geth binary for the specified version
func GethBinaryFilename(os, version string) string {
	var binFileExt string
	if os == "windows" {
		binFileExt = ".exe"
	}
	return fmt.Sprintf("geth-v%s%s", version, binFileExt)
}

// DownloadGethBinary downloads the geth binary for the specified version and OS
func DownloadGethBinary(downloadDir string) error {
	version := params.Version
	hash, err := getReleaseHash(version)
	if err != nil {
		return err
	}

	url, exists, err := GethBinaryURL(runtime.GOOS, runtime.GOARCH, version, hash)
	if err != nil {
		return errors.Wrap(err, "failed to get geth binary url")
	} else if !exists {
		return fmt.Errorf("geth binary for version %s not found", version)
	}

	archive, err := DownloadFile(url, downloadDir)
	if err != nil {
		return errors.Wrap(err, "failed to download geth binary")
	}
	defer os.RemoveAll(archive.Name())

	err = ExtractGethArchive(runtime.GOOS, archive.Name(), downloadDir)
	if err != nil {
		return errors.Wrap(err, "failed to extract geth binary")
	}

	return nil
}

// GethBinaryURL returns the URL of the geth binary for the specified version if the binary exists in remote blob storage
// example call: GethBinaryURL(runtime.GOOS, runtime.GOARCH, "1.10.11", "7231b3efb8095d3")
func GethBinaryURL(os, arch, version, hash string) (string, bool, error) {
	fileExt := "tar.gz"
	if os == "windows" {
		fileExt = "zip"
	}

	if os == "darwin" {
		// no download binaries for M1 macs since runtime.GOARCH == "arm64"
		arch = "amd64"
	}

	gethURL := fmt.Sprintf(
		"https://gethstore.blob.core.windows.net/builds/geth-%s-%s-%s-%.8s.%s",
		os,
		arch,
		version,
		hash,
		fileExt,
	)

	resp, err := http.Head(gethURL)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to get geth binary url")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", false, nil
	}

	return gethURL, true, nil
}

// getReleaseHash retrieves the latest stable geth release data from official go-ethereum github repo.
// version is the semantic version of the geth release
func getReleaseHash(semanticVersion string) (string, error) {
	var response map[string]interface{}

	resp, err := http.Get("https://api.github.com/repos/ethereum/go-ethereum/commits/v" + semanticVersion)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", response["sha"]), nil
}

// DownloadFile downloads a file from url and saves it to the specified filepath.
func DownloadFile(url string, dest string) (*os.File, error) {
	file := path.Base(url)

	log.Info().Msgf("Downloading file %s from %s", filepath.Join(dest, file), url)

	var path bytes.Buffer
	path.WriteString(dest + string(os.PathSeparator) + file)

	start := time.Now()

	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return nil, errors.Wrap(err, "failed to create directory")
	}

	outFile, err := os.Create(path.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file")
	}
	defer outFile.Close()

	headResp, err := http.Head(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get file size")
	}
	defer headResp.Body.Close()

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file")
	}
	defer resp.Body.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to persist file")
	}

	elapsed := time.Since(start)
	log.Info().Msgf("Download completed in %s", elapsed)

	return outFile, nil
}

// ExtractGethArchive extracts the geth binary from .tar.gz or .zip to the specified geth directory.
// The following operations are performed in this function:
// - extract the geth binary from the .tar.gz or .zip (windows) archive
// - move the geth binary to the specified geth directory and rename it to the geth binary name (geth-v<version>(.exe))
// - removes the .tar.gz or .zip (windows)archive
// - change mode of the geth binary to executable
func ExtractGethArchive(goOS, filePath, gethDir string) error {
	log.Info().Msgf("Extracting %s to %s\n", filePath, gethDir)

	start := time.Now()

	var gethBinaryDir string

	switch filepath.Ext(filePath) {
	case ".gz":
		err := ExtractTarGz(filePath, gethDir)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				log.Info().Msg("File already extracted")
			} else {
				return errors.Wrap(err, "failed to extract file")
			}
		}
		gethBinaryDir = strings.TrimSuffix(filePath, ".tar.gz")
	case ".zip":
		err := ExtractZip(filePath, gethDir)
		if err != nil {
			return errors.Wrap(err, "failed to extract file")
		}
		gethBinaryDir = strings.TrimSuffix(filePath, ".zip")
	default:
		return ErrUnsupportedArchive
	}

	log.Info().Msgf("Extraction completed in %s", time.Since(start))

	// move and rename extracted geth binary
	newGethBinaryPath := filepath.Join(gethDir, GethBinaryFilename(goOS, params.Version))
	gethBinaryPath := filepath.Join(gethBinaryDir, "geth")
	if goOS == "windows" {
		gethBinaryPath += ".exe"
	}

	err := os.Rename(gethBinaryPath, newGethBinaryPath)
	if err != nil {
		return errors.Wrap(err, "failed to rename geth binary")
	}

	// remove extracted directory
	err = os.RemoveAll(gethBinaryDir)
	if err != nil {
		return errors.Wrap(err, "failed to remove extracted directory")
	}

	// chmod geth binary to make it executable
	err = os.Chmod(newGethBinaryPath, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to chmod geth binary")
	}

	return nil
}

// ExtractTarGz extracts a .tar.gz file in a specified parent directory.
// The .tar.gz file is extracted to a directory with the same name as the tar.gz file.
func ExtractTarGz(filePath string, outDir string) error {
	gzipStream, err := os.Open(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer gzipStream.Close()

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return errors.Wrap(err, "failed to create gzip reader")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "failed to read tar header")
		}

		outPath := outDir + string(os.PathSeparator) + header.Name

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(outPath, 0755); err != nil {
				return errors.Wrap(err, "Mkdir() failed")
			}
		case tar.TypeReg:
			outFile, err := os.Create(outPath)
			if err != nil {
				return errors.Wrap(err, "Create() failed")
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return errors.Wrap(err, "Copy() failed")
			}
			outFile.Close()
		default:
			return errors.New(fmt.Sprintf(
				"ExtractTarGz: unknown type: %v in %s",
				header.Typeflag,
				header.Name,
			))
		}
	}
	return nil
}

// ExtractZip extracts a .zip file to the specified directory and returns the extracted directory
func ExtractZip(archive, outDir string) error {
	r, err := zip.OpenReader(archive)
	if err != nil {
		return errors.Wrap(err, "failed to open zip archive")
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(outDir, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(outDir)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return errors.Wrap(err, "failed to create file")
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return errors.Wrap(err, "failed to open file")
		}
		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return errors.Wrap(err, "failed to open file")
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return errors.Wrap(err, "failed to copy file")
		}
	}
	return nil
}
