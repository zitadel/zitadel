package orbctl

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func downloadOrbctl(orbctlPath, tag string) error {

	var err error
	orbctlPath, err = filepath.Abs(orbctlPath)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(orbctlPath); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(orbctlPath), os.ModePerm); err != nil {
		return err
	}

	orbctlBase := filepath.Base(orbctlPath)
	url := fmt.Sprintf("https://github.com/caos/orbos/releases/latest/download/%s", orbctlBase)

	if tag != "" {

		/* TODO: Why are dev artifacts released with points??? */
		if !regexp.MustCompile("^v?[0-9]+.[0-9]+.[0-9]$").Match([]byte(tag)) {
			orbctlBase = strings.ReplaceAll(orbctlBase, "-", ".")
		}

		url = fmt.Sprintf("https://github.com/caos/orbos/releases/download/%s/%s", tag, orbctlBase)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(orbctlPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return os.Chmod(orbctlPath, os.ModePerm)
}
