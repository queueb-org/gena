package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSHelper_Suite(t *testing.T) {
	// note, this test case is not testing anything except code coverage.
	// if you find something useful to check, please do.
	t.Run("ok", func(in *testing.T) {
		helper := NewOSHelper().
			WithABS(filepath.Abs).
			WithEvalSymlinks(filepath.EvalSymlinks).
			WithReadFile(os.ReadFile).
			WithStat(os.Stat).
			WithGetwd(os.Getwd)

		if helper == nil {
			in.Errorf("expected not to be nil")
		}
	})
}
