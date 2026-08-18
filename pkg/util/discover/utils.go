package discover

import (
	"strings"
)

func makePath(fullPath string, appName string) string {
	if idx := strings.Index(fullPath, appName); idx > -1 {
		return fullPath[:idx+len(appName)]
	}

	return fullPath
}
