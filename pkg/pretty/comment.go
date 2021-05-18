package pretty

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/godoc/util"

	"github.com/27149chen/go-pretty/pkg/util/sets"
)

const prefix = "// "
const packageLine = "package "

var supportedExtensions = sets.NewString(".go")

func Comment(path string) error {
	return fileWalk(path, commentCode)
}

func Uncomment(path string) error {
	return fileWalk(path, uncommentCode)
}

func commentCode(name, tmpDir string) error {
	return commentOrUncommentCode(name, tmpDir, true)
}

func uncommentCode(name, tmpDir string) error {
	return commentOrUncommentCode(name, tmpDir, false)
}

//func commentFile(name, tmpDir string) error {
//	return commentOrUncommentFile(name, tmpDir, true)
//}
//
//func uncommentFile(name, tmpDir string) error {
//	return commentOrUncommentFile(name, tmpDir, false)
//}

func commentOrUncommentCode(name, tmpDir string, comment bool) error {
	ext := filepath.Ext(name)
	if !supportedExtensions.Has(ext) {
		return nil
	}

	for _, ex := range excludes {
		if strings.Contains(name, ex) {
			return commentOrUncommentFile(name, tmpDir, comment)
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
	var commented bool
	var tagLine bool
	var changed bool
	for scanner.Scan() {
		text := scanner.Text()
		if !util.IsText([]byte(text)) {
			_ = os.Remove(tmpFileName)
			return nil
		}
		if strings.Contains(text, fileTag) {
			_ = os.Remove(tmpFileName)
			return commentOrUncommentFile(name, tmpDir, comment)
		}

		if strings.Contains(text, startTag) {
			commented = true
			tagLine = true
			changed = true
		} else if strings.Contains(text, endTag) {
			commented = false
			tagLine = true
		} else {
			tagLine = false
		}

		if commented && !tagLine {
			if comment {
				text = prefix + text
			} else {
				text = strings.Replace(text, prefix, "", 1)
			}

		}
		_, err := w.WriteString(text + "\n")
		if err != nil {
			return err
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

func commentOrUncommentFile(name, tmpDir string, comment bool) error {
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
	var commented bool
	for scanner.Scan() {
		text := scanner.Text()
		if !util.IsText([]byte(text)) {
			_ = os.Remove(tmpFileName)
			return nil
		}
		if !strings.HasPrefix(text, packageLine) {
			commented = true
		} else {
			commented = false
		}

		if commented {
			if comment {
				text = prefix + text
			} else {
				text = strings.Replace(text, prefix, "", 1)
			}
		}
		_, err := w.WriteString(text + "\n")
		if err != nil {
			return err
		}

	}

	if err := w.Flush(); err != nil {
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Remove(name); err != nil {
		return err
	}

	return os.Rename(tmpFileName, name)
}
