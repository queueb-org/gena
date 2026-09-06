package test

import (
	"path"
	"testing"

	"common.queueb.org/tests"
)

func TestWithTestModule(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		test, sub := tests.F(in).WithTestFatalF()
		WithTestModule(test)
		if sub.Failed() {
			in.Errorf("WithTestModule() expected not to fail")
		}
	})
}

func TestWithInvalidTestModule(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		test, sub := tests.F(in).WithTestFatalF()
		WithInvalidTestModule(test)
		if sub.Failed() {
			in.Errorf("WithInvalidTestModule() expected not to fail")
		}
	})
}

func TestWithPathJoin(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		test, sub := tests.F(in).WithTestFatalF()
		_, join := WithPathJoin(test)
		_ = join("test")

		if sub.Failed() {
			in.Errorf("WithPathJoin() expected not to fail")
		}
	})
}

func TestWithDirEntry(t *testing.T) {
	rootDir := WithTestModule(t)

	t.Run("ok", func(in *testing.T) {
		test, sub := tests.F(in).WithTestFatalF()
		WithDirEntry(test, path.Join(rootDir, "pkg"), "config")
		if sub.Failed() {
			in.Errorf("WithDirEntry() expected not to fail")
		}
	})

	t.Run("not-dir", func(in *testing.T) {
		test, sub := tests.F(in).WithTestFatalF()

		WithDirEntry(test, path.Join(rootDir, "go.mod"), "config")
		if !sub.Failed() {
			in.Errorf("WithDirEntry() expected to fail")
		}
	})

	t.Run("not-found", func(in *testing.T) {
		test, sub := tests.F(in).WithTestFatalF()

		WithDirEntry(test, path.Join(rootDir, "pkg"), "config-v1")
		if !sub.Failed() {
			in.Errorf("WithDirEntry() expected to fail")
		}
	})
}
