package util

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"log"
)

func Unzip(zipFile, destDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer safeClose(r)

	os.MkdirAll(destDir, 0755) // 7=rwx 5=r-x 5=r-x

	for _, f := range r.File {
		err := extractFileFromZip(f, destDir)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFileFromZip(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer safeClose(rc)

	path := filepath.Join(dest, f.Name)

	if f.FileInfo().IsDir() {
		os.MkdirAll(path, f.Mode())
	} else {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer safeClose(f)

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}

func safeClose(f io.Closer) {
	if err := f.Close(); err != nil {
		log.Print(err)
	}
}