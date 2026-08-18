package discover

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"common.queueb.org/tests"
	"queueb.org/gena/pkg/util"
	"queueb.org/gena/pkg/util/test"
	"queueb.org/gena/pkg/util/types"
)

const (
	defaultAnnotation = "+k8s:deepcopy-gen=package"
)

// Shortcuts
var (
	WithTestModule        = test.WithTestModule
	WithInvalidTestModule = test.WithInvalidTestModule
	WithPathJoin          = test.WithPathJoin
	WithDirEntry          = test.WithDirEntry
	writeTestFile         = tests.WithWriteFileString
)

// WithHelper sets custom helper for internal functions management.
func WithHelper(t testing.TB, newHelper *util.OSHelper) {
	orig := helper

	t.Cleanup(func() {
		helper = orig
	})

	helper = newHelper
}

// Tests

func TestCollectApps(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		dir := WithTestModule(in)
		apps, err := CollectApps([]string{dir}, defaultAnnotation)
		if err != nil {
			in.Fatalf("got error: %v", err)
		}

		expectedAppsLen := 1
		if result := len(apps); result != expectedAppsLen {
			in.Errorf("expected: %v, got: %d", expectedAppsLen, result)
		}
	})

	t.Run("on-err", func(in *testing.T) {
		dir := WithInvalidTestModule(in)
		_, err := CollectApps([]string{dir}, defaultAnnotation)
		if err == nil {
			in.Errorf("expected error got nil instead")
		}
	})
}

func TestFindPackages(t *testing.T) {
	root := WithTestModule(t)

	tests := []struct {
		name       string
		roots      []string
		annotation string
		want       []string
		wantErr    string
	}{
		{
			name:       "recursive, sorted and filtered",
			roots:      []string{root},
			annotation: "//k8s:deepcopy-gen=package",
			want: []string{
				"example.org/application/pkg/config",
				"example.org/application/pkg/config/v1alpha1",
			},
		},
		{
			name:       "multiple roots are deduplicated",
			roots:      []string{filepath.Join(root, "pkg", "config"), filepath.Join(root, "pkg", "config", "v1alpha1")},
			annotation: defaultAnnotation,
			want: []string{
				"example.org/application/pkg/config",
				"example.org/application/pkg/config/v1alpha1",
			},
		},
		{
			name:       "custom block comment annotation",
			roots:      []string{root},
			annotation: "+example:generate=true",
			want:       []string{"example.org/application/internal/custom"},
		},
		{
			name:       "empty annotation",
			roots:      []string{root},
			annotation: " ",
			wantErr:    "annotation must not be empty",
		},
		{
			name:       "not a directory",
			roots:      []string{filepath.Join(root, "go.mod")},
			annotation: defaultAnnotation,
			wantErr:    "is not a directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(in *testing.T) {
			got, err := findPackages(test.roots, test.annotation)
			if test.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					in.Fatalf("got error %v, want one containing %q", err, test.wantErr)
				}
				return
			}
			if err != nil {
				in.Fatalf("got error: %v", err)
			}

			packages := got.Collect()
			if !reflect.DeepEqual(packages, test.want) {
				in.Errorf("got packages %v, want %v", packages, test.want)
			}
		})
	}
}

func TestFindPackages_SymlinkRoot(t *testing.T) {
	root := WithTestModule(t)
	link := filepath.Join(t.TempDir(), "application")
	if err := os.Symlink(root, link); err != nil {
		t.Skipf("cannot create directory symlink: %v", err)
	}

	got, err := findPackages([]string{link}, defaultAnnotation)
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	want := []string{
		"example.org/application/pkg/config",
		"example.org/application/pkg/config/v1alpha1",
	}

	packages := got.Collect()
	if !reflect.DeepEqual(packages, want) {
		t.Errorf("got packages %v, want %v", packages, want)
	}
}

func TestCollectPackagesDiscovered(t *testing.T) {
	rootDir := WithTestModule(t)

	for _, test := range []struct {
		name    string
		setup   func(t testing.TB) ([]string, string)
		wantErr bool
	}{
		{
			name: "ok",
			setup: func(t testing.TB) ([]string, string) {
				return []string{rootDir}, defaultAnnotation
			},
			wantErr: false,
		},
		{
			name: "abs-err",
			setup: func(t testing.TB) ([]string, string) {

				WithHelper(t, util.NewOSHelper().WithABS(func(s string) (string, error) {
					return "", errors.New("issue")
				}))
				return []string{rootDir}, defaultAnnotation
			},
			wantErr: true,
		},
		{
			name: "eval-symlinks-err",
			setup: func(t testing.TB) ([]string, string) {
				WithHelper(t, util.NewOSHelper().WithEvalSymlinks(func(string) (string, error) {
					return "", errors.New("issue")
				}))
				return []string{rootDir}, defaultAnnotation
			},
			wantErr: true,
		},
		{
			name: "no-dir",
			setup: func(t testing.TB) ([]string, string) {
				WithHelper(t, util.NewOSHelper().WithStat(func(name string) (os.FileInfo, error) {
					return nil, errors.New("error")
				}))
				return []string{rootDir}, defaultAnnotation
			},
			wantErr: true,
		},
		{
			name: "is-file",
			setup: func(t testing.TB) ([]string, string) {
				return []string{path.Join(rootDir, "go.mod")}, defaultAnnotation
			},
			wantErr: true,
		},
		{
			name: "no-go-mod",
			setup: func(t testing.TB) ([]string, string) {
				return []string{path.Join(rootDir, "../")}, defaultAnnotation
			},
			wantErr: true,
		},
		{
			name: "walk-packages-err",
			setup: func(t testing.TB) ([]string, string) {
				dir := WithInvalidTestModule(t)
				return []string{path.Join(dir)}, defaultAnnotation
			},
			wantErr: true,
		},
	} {
		t.Run(test.name, func(in *testing.T) {
			roots, annotation := test.setup(in)
			_, err := collectPackagesDiscovered(roots, annotation)
			if !test.wantErr && err != nil {
				in.Fatalf("got error: %v", err)
			}
			if test.wantErr && err == nil {
				in.Errorf("expected error but got nil instead")
			}
		})
	}
}

func TestDirectoryHasAnnotation(t *testing.T) {
	annotation := normalizeAnnotation(defaultAnnotation)

	t.Run("ok", func(in *testing.T) {
		dir := WithTestModule(in)
		has, err := directoryHasAnnotation(path.Join(dir, "pkg/config"), annotation)
		if !has || err != nil {
			in.Errorf("want: (true, nil), got: (%#v, %#v)", has, err)
		}
	})

	t.Run("non-existent-dir", func(in *testing.T) {
		_, err := directoryHasAnnotation("non existent dir", annotation)
		if err == nil {
			in.Errorf("want error but got nil instead")
		}
	})

	t.Run("broken-go-file", func(in *testing.T) {
		dir, p := WithPathJoin(in)
		writeTestFile(in, p("go.mod"), "module example.fqdn\n\ngo 1.26\n")
		writeTestFile(in, p("main.go"), "invalid go file\n")

		_, err := directoryHasAnnotation(dir, annotation)
		if err == nil {
			in.Errorf("want error but got nil instead")
		}
	})
}

func TestFindGoModule(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		dir := WithTestModule(in)
		module, err := findGoModule(dir)
		if err != nil {
			in.Errorf("got error: %v", err)
		}
		want := types.GoModule{
			Path: dir,
			Name: "example.org/application",
		}
		if !reflect.DeepEqual(want, module) {
			in.Errorf("\nexp: %#v\ngot: %#v\n", want, module)
		}
	})

	t.Run("invalid-go-mod", func(in *testing.T) {
		dir, p := WithPathJoin(in)
		writeTestFile(in, p("go.mod"), "invalid module\n")
		if _, err := findGoModule(dir); err == nil {
			in.Errorf("want error but got nil instead")
		}
	})

	t.Run("no-go-mod", func(in *testing.T) {
		dir, p := WithPathJoin(in)
		writeTestFile(in, p("main.go"), "package main\n\n")
		if _, err := findGoModule(dir); err == nil {
			in.Errorf("want error but got nil instead")
		}
	})

	t.Run("custom-read-err", func(in *testing.T) {
		WithHelper(in, util.NewOSHelper().WithReadFile(func(s string) ([]byte, error) {
			return nil, errors.New("a custom error")
		}))

		dir, p := WithPathJoin(in)
		writeTestFile(in, p("main.go"), "package main\n\n")
		if _, err := findGoModule(dir); err == nil {
			in.Errorf("want error but got nil instead")
		}
	})
}

func TestWalkPackages(t *testing.T) {
	// setup
	dir := WithTestModule(t)
	module, err := findGoModule(dir)
	if err != nil {
		t.Fatalf("no suitable go module found")
	}
	annotation := normalizeAnnotation(defaultAnnotation)

	t.Run("ok", func(in *testing.T) {
		discovered := make(map[string]rawPackage)
		err = walkPackages(dir, module, annotation, discovered)
		if err != nil {
			in.Errorf("got error: %v", err)
		}
	})
}

func Test_testWalk(t *testing.T) {
	t.Run("cannot-walk", func(in *testing.T) {
		if _, err := testWalk(walkOpt{
			walkErr: errors.New("any issue"),
		}); err == nil {
			in.Errorf("got error: %v", err)
		}
	})

	t.Run("no-go-mod", func(in *testing.T) {
		annotation := normalizeAnnotation(defaultAnnotation)

		WithHelper(in, util.NewOSHelper().WithStat(func(name string) (os.FileInfo, error) {
			return nil, errors.New("any issue")
		}))

		dir, p := WithPathJoin(in)
		writeTestFile(in, p("go.mod"), "module example.fqdn\n\ngo 1.26")
		writeTestFile(in, p("sub/main.go"), "package main\n\n")

		module := types.GoModule{Path: dir, Name: "example.fqdn"}
		entry := WithDirEntry(in, dir, "sub")

		_, err := testWalk(walkOpt{
			root:       dir,
			module:     module,
			annotation: annotation,
			path:       path.Join(dir, "sub"),
			entry:      entry,
		})

		if err == nil {
			in.Errorf("got error: %v", err)
		}
	})
}

func TestTestWalk_Suite(t *testing.T) {
	rootDir := WithTestModule(t)
	annotation := normalizeAnnotation(defaultAnnotation)

	for _, test := range []struct {
		name  string
		setup func(t testing.TB) walkOpt
		skip  bool
	}{
		{
			name: "ok",
			setup: func(t testing.TB) walkOpt {
				target := path.Join(rootDir, "pkg/config")
				return walkOpt{
					root:       rootDir,
					path:       target,
					annotation: annotation,
					entry:      WithDirEntry(t, path.Join(rootDir, "pkg"), "config"),
				}
			},
			skip: false,
		},
		{
			name: "walk-error",
			setup: func(t testing.TB) walkOpt {
				return walkOpt{walkErr: errors.New("issue")}
			},
			skip: true,
		},
		{
			name: "not-dir",
			setup: func(t testing.TB) walkOpt {
				return walkOpt{
					entry: WithDirEntry(t, rootDir, "go.mod"),
				}
			},
			skip: true,
		},
		{
			name: "skip-dir",
			setup: func(t testing.TB) walkOpt {
				return walkOpt{
					path:  "example.fqdn",
					root:  "testdata",
					entry: WithDirEntry(t, rootDir, "testdata"),
				}
			},
			skip: true,
		},
		{
			name: "no-annotation",
			setup: func(t testing.TB) walkOpt {
				target := path.Join(rootDir, "internal/custom")
				return walkOpt{
					root:       rootDir,
					path:       target,
					annotation: annotation,
					entry:      WithDirEntry(t, path.Join(rootDir, "internal"), "custom"),
				}
			},
			skip: true,
		},
		{
			name: "sub-module",
			setup: func(t testing.TB) walkOpt {
				target := path.Join(rootDir, "nested")
				return walkOpt{
					root:       rootDir,
					path:       target,
					annotation: annotation,
					entry:      WithDirEntry(t, rootDir, "nested"),
				}
			},
			skip: true,
		},
		{
			name: "custom-err",
			setup: func(t testing.TB) walkOpt {
				WithHelper(t, util.NewOSHelper().WithStat(func(name string) (os.FileInfo, error) {
					return nil, errors.New("custom error")
				}))

				target := path.Join(rootDir, "nested")
				return walkOpt{
					root:       rootDir,
					path:       target,
					annotation: annotation,
					entry:      WithDirEntry(t, rootDir, "nested"),
				}
			},
			skip: true,
		},
	} {
		t.Run(test.name, func(in *testing.T) {
			if want, _ := testWalk(test.setup(in)); want != test.skip {
				in.Errorf("want: %v, got: %v", want, test.skip)
			}
		})
	}
}
