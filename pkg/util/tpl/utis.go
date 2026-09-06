package tpl

import (
	"os"
	"path"
	"path/filepath"

	"queueb.org/gena/pkg/util"
)

const (
	EnvProjectRoot = "GENA_PROJECT_ROOT"
)

var (
	helper = util.NewOSHelper()
)

// ProjectDir discovers project root or reads it from environment.
func ProjectDir(fallback string) (root string) {
	// set fallback as defaults.
	root = fallback

	// Use environment variable
	if dir, found := os.LookupEnv(EnvProjectRoot); found {
		return dir
	}

	pwd, err := helper.Getwd()
	if err != nil {
		return
	}

	for current := pwd; ; current = filepath.Dir(current) {
		stat, err := helper.Stat(path.Join(current, ".gena.yaml"))
		if err == nil && !stat.IsDir() {
			return current
		}

		if current == "/" {
			break
		}
	}

	return
}
