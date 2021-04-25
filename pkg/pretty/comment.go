package pretty

import (
	"bufio"
	"os"
	"strings"
)

const prefix = "// "
const packageLine = "package "

func Comment(path string) error {
	return fileWalk(path, commentCode)
}

func Uncomment(path string) error {
	return fileWalk(path, uncommentCode)
}

func commentCode(name, tmpDir string) error {
	for _, ex := range excludes {
		if strings.Contains(name, ex) {
			return commentFile(name)
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
	var commented bool
	var tagLine bool
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, fileTag) {
			_ = os.Remove(name+"_tmp")
			return commentFile(name)

		}

		if strings.Contains(text, startTag) {
			commented = true
			tagLine = true
		} else if strings.Contains(text, endTag) {
			commented = false
			tagLine = true
		} else {
			tagLine = false
		}

		if commented && !tagLine {
			text = prefix + text
		}
		_, err := w.WriteString(text + "\n")
		if err != nil {
			return err
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

func uncommentCode(name, tmpDir string) error {
	for _, ex := range excludes {
		if strings.Contains(name, ex) {
			return uncommentFile(name)
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
	var commented bool
	var tagLine bool
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, fileTag) {
			_ = os.Remove(name+"_tmp")
			return uncommentFile(name)

		}

		if strings.Contains(text, startTag) {
			commented = true
			tagLine = true
		} else if strings.Contains(text, endTag) {
			commented = false
			tagLine = true
		} else {
			tagLine = false
		}

		if commented && !tagLine {
			text = strings.Replace(text, prefix, "", 1)
		}
		_, err := w.WriteString(text + "\n")
		if err != nil {
			return err
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

func commentFile(name string) error {
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
	var commented bool
	for scanner.Scan() {
		text := scanner.Text()
		if !strings.HasPrefix(text, packageLine) {
			commented = true
		} else {
			commented = false
		}


		if commented {
			text = prefix + text
		}
		_, err := w.WriteString(text + "\n")
		if err != nil {
			return err
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

func uncommentFile(name string) error {
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
	var commented bool
	for scanner.Scan() {
		text := scanner.Text()
		if !strings.HasPrefix(text, packageLine) {
			commented = true
		} else {
			commented = false
		}


		if commented {
			text = strings.Replace(text, prefix, "", 1)
		}
		_, err := w.WriteString(text + "\n")
		if err != nil {
			return err
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
