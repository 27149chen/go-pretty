package pkg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const PrettyFile = ".prettyfile"

var (
	ReservedTag = "Reserved for enterprise only"
	startTag = fmt.Sprintf("=== Code %s. START ===", ReservedTag)
	endTag = fmt.Sprintf("=== Code %s. END ===", ReservedTag)
	fileTag = fmt.Sprintf("=== File %s. ===", ReservedTag)
)

var excludes []string
var ignores = []string{".git", "idea"}

type prettyFunc func(name string) error

func PopulateExcludes(prettyIgnore string) error {
	file, err := os.Open(prettyIgnore)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "#") || strings.HasPrefix(text, "/") {
			continue
		}
		excludes = append(excludes, text)
	}

	return nil
}

func Prettify(path string) error {
	return prettify(path, cleanCode)
}

func cleanCode(name string) error {
	for _, ex := range excludes {
		if strings.Contains(name, ex) {
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

func prettify(path string, action prettyFunc) error {
	return filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				panic(err)
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			for _, ig := range ignores {
				if strings.Contains(path, ig) {
					return nil
				}
			}

			fmt.Printf("Processing %s\n", path)
			return action(path)
		})
}
