package tpl

import (
	"errors"
	"os"
	"testing"

	"common.queueb.org/tests"
	"queueb.org/gena/pkg/util"
	"queueb.org/gena/pkg/util/test"
)

// WithHelper sets custom helper for internal functions management.
func WithHelper(t testing.TB, newHelper *util.OSHelper) {
	orig := helper

	t.Cleanup(func() {
		helper = orig
	})

	helper = newHelper
}

func TestProjectDir(t *testing.T) {
	t.Run("from-env", func(in *testing.T) {
		dir := in.TempDir()
		in.Setenv(EnvProjectRoot, dir)
		if result := ProjectDir("/root"); result != dir {
			in.Errorf("expected: %v, got: %v", dir, result)
		}
	})

	t.Run("cant-get-workingdir", func(in *testing.T) {
		WithHelper(in, util.NewOSHelper().WithGetwd(func() (string, error) {
			return "", errors.New("any")
		}))
		expected := "/root"
		if result := ProjectDir("/root"); result != expected {
			in.Errorf("expected: %v, got: %v", expected, result)
		}
	})

	t.Run("ok", func(in *testing.T) {
		dir := test.WithTestModule(in)
		if err := os.Chdir(dir); err != nil {
			in.Fatalf("could not chdir: %v", err)
		}
		tests.WithWriteFileString(in, ".gena.yaml", "kind: Config\napiVersion: gena.apps.queueb.org/v1alpha1\n")
		if result := ProjectDir("/root"); result != dir {
			in.Errorf("expected: %v, got: %v", dir, result)
		}
	})

	t.Run("fallback", func(in *testing.T) {
		dir := in.TempDir()
		if err := os.Chdir(dir); err != nil {
			in.Fatalf("could not chdir: %v", err)
		}
		expected := "/root"
		if result := ProjectDir("/root"); result != expected {
			in.Errorf("expected: %v, got: %v", expected, result)
		}
	})
}
