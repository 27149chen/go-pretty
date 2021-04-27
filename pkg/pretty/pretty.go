package pretty

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"golang.org/x/tools/godoc/util"
)

var (
	ReservedTag = "Reserved for enterprise only"
	startTag = fmt.Sprintf("=== Code %s. START ===", ReservedTag)
	endTag = fmt.Sprintf("=== Code %s. END ===", ReservedTag)
	fileTag = fmt.Sprintf("=== File %s. ===", ReservedTag)
)

var excludes []string
var ignores = []string{".git", "idea", ".DS_Store"}

type prettyFunc func(name, tmpDir string) error

func PopulateExcludes(prettyIgnore string) error {
	file, err := os.Open(prettyIgnore)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		if strings.HasPrefix(text, "#") || strings.HasPrefix(text, "//") {
			continue
		}
		excludes = append(excludes, text)
	}

	return nil
}

func Prettify(path string) error {
	return fileWalk(path, cleanCode)
}

func cleanCode(name, tmpDir string) error {
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

	tmpFileName := filepath.Join(tmpDir, filepath.Base(name))
	tmpFile, err := os.Create(tmpFileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = tmpFile.Close()
	}()

	w := bufio.NewWriter(tmpFile)

	scanner := bufio.NewScanner(file)
	var skipped bool
	var changed bool
	for scanner.Scan() {
		text := scanner.Text()
		if !util.IsText([]byte(text)) {
			_ = os.Remove(tmpFileName)
			return nil
		}
		if strings.Contains(text, fileTag) {
			_ = os.Remove(tmpFileName)
			return os.Remove(name)

		}

		if strings.Contains(text, startTag) {
			skipped = true
			changed = true
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
	if err := tmpFile.Close(); err != nil {
		return err
	}

	if changed {
		if err := os.Remove(name); err != nil {
			return err
		}

		return os.Rename(tmpFileName, name)
	}

	return os.Remove(tmpFileName)
}

func fileWalk(path string, action prettyFunc) error {
	tmp, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(tmp)
	}()

	return godirwalk.Walk(path, &godirwalk.Options{
		Callback: func(name string, de *godirwalk.Dirent) error {
			for _, ig := range ignores {
				if strings.Contains(name, ig) {
					return godirwalk.SkipThis
				}
			}

			if !de.IsRegular() {
				return nil
			}

			fmt.Printf("Processing %s\n", name)
			return action(name, tmp)
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
}
