package archiver

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Archiver is an interface for which implementors will archive src files to a dest
type Archiver interface {
	Archive(src, dest string) error
}

type zipper struct{}

// Zip is an Archiver that zips/unzips files
var Zip Archiver = (*zipper)(nil)

// Archive copies the contents of the src directory into the destination directory
func (z *zipper) Archive(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return err
	}

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	w := zip.NewWriter(destination)
	defer w.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		// skip directories because paths of files encode that information for us already
		if info.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		curr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer curr.Close()

		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, curr)
		if err != nil {
			return err
		}

		return nil
	})
}
