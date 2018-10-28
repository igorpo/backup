package hasher

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// HashDir will take the directory and compute an MD5 hash based on its contents. We care about
// the directory modification time, file mode, name, size, and path. If any of these properties change,
// we want to update the hash.
func HashDir(path string) (string, error) {
	hash := md5.New()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		io.WriteString(hash, path)
		fmt.Fprintf(hash, "%v", info.IsDir())
		fmt.Fprintf(hash, "%v", info.ModTime())
		fmt.Fprintf(hash, "%v", info.Mode())
		fmt.Fprintf(hash, "%v", info.Name())
		fmt.Fprintf(hash, "%v", info.Size())
		return nil
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
