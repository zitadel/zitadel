package helpers

import (
	"os"
	"strings"
)

func PruneHome(pwd string) string {
	if strings.HasPrefix(pwd, "~") {
		userhome, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		pwd = userhome + pwd[1:]
	}
	return pwd
}
