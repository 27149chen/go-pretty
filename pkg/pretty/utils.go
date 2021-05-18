package pretty

import (
	"fmt"
	"os"
	"path/filepath"
)

func RemoveEmptyDir(root string) error {
	_, err := removeEmptyDir(root)

	return err
}

func removeEmptyDir(root string) (bool, error) {
	files, err := os.ReadDir(root)
	if err != nil {
		return false, err
	}

	if len(files) == 0 {
		fmt.Printf("Removing empty dir: %s\n", root)
		return true, os.Remove(root)
	}

	deletedDirs := 0
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		deleted, err := removeEmptyDir(filepath.Join(root, file.Name()))
		if err != nil {
			return false, err
		}
		if deleted {
			deletedDirs += 1
		}
	}

	if deletedDirs == len(files) {
		fmt.Printf("Removing empty dir: %s\n", root)
		return true, os.Remove(root)
	}

	return false, nil
}
