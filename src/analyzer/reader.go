package analyzer

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

// Unpacks archive into files
func openPackageArchive(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open archive: %w", err)
	}

	gz, err := gzip.NewReader(file)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("read gzip archive: %w", err)
	}
	return struct {
		io.Reader
		io.Closer
	}{
		Reader: gz,
		Closer: closeBoth{gz, file},
	}, nil
}

type closeBoth struct {
	io.Closer
	file *os.File
}


// closes file and reader both
func (c closeBoth) Close() error {
	err := c.Closer.Close()
	fileErr := c.file.Close()
	if err != nil {
		return err
	}
	return fileErr
}
