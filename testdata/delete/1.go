package delete

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	start = "// === start ==="
	end = "// === end ==="
)

func cleanCode(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	tmpFile, err := os.Create(name+"_tmp")
	if err != nil {
		return err
	}

	w := bufio.NewWriter(tmpFile)

	scanner := bufio.NewScanner(file)
	var skipped bool
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, start) {
			skipped = true
		} else if strings.Contains(text, end) {
			skipped = false
			continue
		}

		if !skipped {
			_, err := w.WriteString(text + "\n")
			if err != nil {
				return err
			}
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	if err := os.Remove(name); err != nil {
		return err
	}

	return os.Rename(name+"_tmp", name)
}

func fileWalk(path string) error {
	return filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}

			if !info.Mode().IsRegular() || !strings.HasSuffix(info.Name(), ".go") {
				return nil
			}
			fmt.Printf("Processing %s\n", path)
			return cleanCode(path)
		})
}
