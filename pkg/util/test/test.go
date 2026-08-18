package test

import (
	"os"
	"path"
	"testing"

	"common.queueb.org/tests"
)

// shortcuts
var (
	writeTestFile = tests.WithWriteFileString
)

// WithTestModule creates simple testing module and returns its
// directory location. Temp directory would be removed after test is done.
func WithTestModule(t testing.TB) (rootDir string) {
	t.Helper()
	root := t.TempDir()
	p := func(in string) string {
		return path.Join(root, in)
	}

	writeTestFile(t, p("go.mod"), "module example.org/application\n\ngo 1.26\n")
	writeTestFile(t, p("pkg/config/doc.go"), "// +k8s:deepcopy-gen=package\npackage config\n")
	writeTestFile(t, p("pkg/config/v1alpha1/doc.go"), "// +k8s:deepcopy-gen=package\npackage v1alpha1\n")
	writeTestFile(t, p("pkg/config/schema/schema.go"), "package schema\n")
	writeTestFile(t, p("internal/custom/doc.go"), "/*\n * +example:generate=true\n */\npackage custom\n")
	writeTestFile(t, p("testdata/ignored/doc.go"), "// +k8s:deepcopy-gen=package\npackage ignored\n")
	writeTestFile(t, p("nested/go.mod"), "module example.org/nested\n")
	writeTestFile(t, p("nested/doc.go"), "// +k8s:deepcopy-gen=package\npackage nested\n")
	writeTestFile(t, p("ignored/one/_doc.go"), "// +k8s:deepcopy-gen=package\npackage one\n")
	writeTestFile(t, p("ignored/two/.doc.go"), "// +k8s:deepcopy-gen=package\npackage two\n")

	return root
}

// WithInvalidTestModule sets invalid module for testing.
func WithInvalidTestModule(t testing.TB) (rootDir string) {
	t.Helper()
	root := t.TempDir()
	p := func(in string) string {
		return path.Join(root, in)
	}

	writeTestFile(t, p("go.mod"), "module example.org/application\n\ngo 1.26\n")
	writeTestFile(t, p("pkg/config/doc.go"), "// +k8s:deepcopy-gen=package\npackage config\n")
	writeTestFile(t, p("pkg/config/main.go"), "invalid golang file syntax\n")

	return root
}

// WithPathJoin configures directory path joining function.
func WithPathJoin(t testing.TB) (dir string, fnc func(in string) string) {
	dir = t.TempDir()
	fnc = func(in string) string {
		return path.Join(dir, in)
	}

	return
}

// WithDirEntry searches 'search' entry for target dir.
func WithDirEntry(t testing.TB, dir string, search string) os.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("could not read directory %s: %v", dir, err)
		return nil
	}

	for _, entry := range entries {
		if entry.Name() == search {
			return entry
		}
	}

	t.Fatalf("`%s` entry not found", search)
	return nil
}
