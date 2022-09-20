package actions

import (
	"strings"

	"github.com/zitadel/logging"
)

type parameter map[string]interface{}

func (param parameter) set(name string, value interface{}) {
	param[name] = value
}

func (param parameter) setPath(path []string, value interface{}) {
	parent := param
	var ok bool
	for _, p := range path[:len(path)-1] {
		if _, ok := parent[p]; !ok {
			parent[p] = parameter{}
		}
		if parent, ok = parent[p].(parameter); !ok {
			logging.WithFields("path", strings.Join(path, "/")).Warn("overwritten path")
			panic("non parameter type overwritten")
		}
	}
	parent[path[len(path)-1]] = value
}
