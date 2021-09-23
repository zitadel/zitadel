package chore

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func downloadZitadelctl(zitadelctlPath, tag string) error {

	var err error
	zitadelctlPath, err = filepath.Abs(zitadelctlPath)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(zitadelctlPath); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(zitadelctlPath), os.ModePerm); err != nil {
		return err
	}

	orbctlBase := filepath.Base(zitadelctlPath)
	url := fmt.Sprintf("https://github.com/caos/zitadel/releases/latest/download/%s", orbctlBase)

	if tag != "" {
		/* TODO: Why are dev artifacts released with points??? */
		if !regexp.MustCompile("^v?[0-9]+.[0-9]+.[0-9]$").Match([]byte(tag)) {
			orbctlBase = strings.ReplaceAll(orbctlBase, "-", ".")
		}

		url = fmt.Sprintf("https://github.com/caos/zitadel/releases/download/%s/%s", tag, orbctlBase)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(zitadelctlPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}

	return os.Chmod(zitadelctlPath, os.ModePerm)
}
