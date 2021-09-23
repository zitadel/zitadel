package chore_test

import "os"

func prefixedEnv(env string) string {
	return os.Getenv("ORBOS_E2E_" + env)
}
