package analyzer

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const mkDirRights = 0755

// downloading errors
var (
	FetchFileError = errors.New("Failed to download file from url")
	CreateFileError = errors.New("Failed to create output file")
	CopyDataError = errors.New("Failed to copy data from response")
)

// Downloads archive
func downloadPackageArchive(url, outFolder string) (string, error) {
	if err := os.MkdirAll(outFolder, mkDirRights); err != nil {
		return "", fmt.Errorf("create cache directory: %w", err)
	}

	path := filepath.Join(outFolder, "Packages.gz")

	if info, err := os.Stat(path); err == nil && info.Size() > 0 {
		return path, nil // используем кэш
	} else if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("check cache: %w", err)
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", FetchFileError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: HTTP %s", FetchFileError, resp.Status)
	}


	out, err := os.Create(path)
	if err != nil {
		return "", CreateFileError
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", CopyDataError
	}

	if err := out.Close(); err != nil {
		return "", CopyDataError
	}

	return path, nil
}