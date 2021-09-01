package rsync

import (
	"os/exec"
	"path/filepath"
)

func GetCommand(
	configPath string,
	sourceName string,
	sourceFilePath string,
	destinationName string,
	destinationFilePath string,
) *exec.Cmd {
	return exec.Command(
		"rclone",
		"--no-check-certificate",
		"--config="+configPath,
		"sync",
		sourceName+":"+sourceFilePath,
		destinationName+":"+filepath.Join(destinationFilePath, sourceFilePath),
	)
}
