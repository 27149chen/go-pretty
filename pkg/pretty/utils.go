package pretty

import (
	"os"
	"path/filepath"
)

func RemoveEmptyDir(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		files, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		if len(files) != 0 {
			return nil
		}

		return os.Remove(path)
	})
}
