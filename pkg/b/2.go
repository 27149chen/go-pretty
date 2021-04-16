package b

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	start1 = "// === start ==="
	end1 = "// === end ==="
)

func cleanCode1(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// === Code Reserved for enterprise only. START ===
	tmpFile, err := os.Create(name+"_tmp")
	if err != nil {
		return err
	}
	// === Code Reserved for enterprise only. END ===

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

	// === Code Reserved for enterprise only. START ===
	if err := os.Remove(name); err != nil {
		return err
	}

	return os.Rename(name+"_tmp", name)
	// === Code Reserved for enterprise only. END ===
}

func fileWalk1(path string) error {
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
