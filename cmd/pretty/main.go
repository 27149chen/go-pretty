package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ReservedTag = "Reserved for enterprise only"
	startTag = fmt.Sprintf("=== Code %s. START ===", ReservedTag)
	endTag = fmt.Sprintf("=== Code %s. END ===", ReservedTag)
	fileTag = fmt.Sprintf("=== File %s. ===", ReservedTag)
)

var excludedPaths = []string{"pkg/delete/"}

func cleanCode(name string) error {
	for _, p := range excludedPaths {
		if strings.Contains(name, p) {
			return os.Remove(name)
		}
	}

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
		if strings.Contains(text, fileTag) {
			_ = os.Remove(name+"_tmp")
			return os.Remove(name)

		}

		if strings.Contains(text, startTag) {
			skipped = true
		} else if strings.Contains(text, endTag) {
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

func main() {
	err := fileWalk(".")
	if err != nil {
		panic(err)
	}
}
