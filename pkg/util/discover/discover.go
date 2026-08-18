package discover

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"

	"queueb.org/gena/pkg/util"
	"queueb.org/gena/pkg/util/types"
)

var (
	helper = util.NewOSHelper()
)

// Type aliases.
type (
	// GoModule = types.GoModule
	Package  = types.Package
	Packages = types.Packages
)

type rawPackage struct {
	Name    string
	Module  string
	Related string
	Path    string
	ModPath string
}

type rawPackages []rawPackage

// Collect returns packages in plain form.
// Applies sort by default.
func (p *rawPackages) Collect() (out []string) {
	for _, item := range *p {
		out = append(out, item.Module)
	}

	sort.Strings(out)
	return
}

func newPackage(p rawPackage) types.Package {
	return types.Package{
		Module:  p.Module,
		Related: p.Related,
		Path:    p.Path,
	}
}

// CollectApps traverses and searches applications with its sub-modules (packages).
func CollectApps(roots []string, annotation string) (apps []*types.App, err error) {
	var packages rawPackages
	if packages, err = findPackages(roots, annotation); err != nil {
		return
	}

	container := make(map[string]*types.App)
	for _, p := range packages {
		appName := p.Name
		if app, exists := container[appName]; !exists {
			container[appName] = &types.App{
				Name:     appName,
				Path:     p.ModPath,
				Packages: types.Packages{newPackage(p)},
			}
		} else {
			app.Packages = append(app.Packages, newPackage(p))
		}
	}

	return slices.Collect(maps.Values(container)), nil
}

// findPackages traverses and searches for an annotation in sources.
// Example:
//
//	findPackages([]string{"."}, "+k8s:deepcopy-gen=package")
//	// []string{"queueb.org/gena/pkg/config", "queueb.org/gena/pkg/config/v1alpha1"}
func findPackages(roots []string, annotation string) (packages rawPackages, err error) {
	annotation = normalizeAnnotation(annotation)

	if annotation == "" {
		return nil, errors.New("annotation must not be empty")
	}

	var discovered map[string]rawPackage
	if discovered, err = collectPackagesDiscovered(roots, annotation); err != nil {
		return
	}

	packages = slices.Collect(maps.Values(discovered))
	return
}

func collectPackagesDiscovered(roots []string, annotation string) (discovered map[string]rawPackage, err error) {
	discovered = make(map[string]rawPackage)

	for _, root := range roots {
		absoluteRoot, err := helper.Abs(root)
		if err != nil {
			return nil, fmt.Errorf("resolve directory %q: %w", root, err)
		}
		// WalkDir does not follow symbolic links, including when its root is a
		// link. Resolve only the explicitly supplied root; links encountered
		// while walking the module remain untouched.
		absoluteRoot, err = helper.EvalSymlinks(absoluteRoot)
		if err != nil {
			return nil, fmt.Errorf("resolve directory symlinks %q: %w", root, err)
		}
		info, err := helper.Stat(absoluteRoot)
		if err != nil {
			return nil, fmt.Errorf("inspect directory %q: %w", root, err)
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("%q is not a directory", root)
		}

		module, err := findGoModule(absoluteRoot)
		if err != nil {
			return nil, fmt.Errorf("discover module for %q: %w", root, err)
		}

		if err := walkPackages(absoluteRoot, module, annotation, discovered); err != nil {
			return nil, fmt.Errorf("discover packages in %q: %w", root, err)
		}
	}

	return
}

func findGoModule(directory string) (types.GoModule, error) {
	for current := directory; ; current = filepath.Dir(current) {
		goMod := filepath.Join(current, "go.mod")
		contents, err := helper.ReadFile(goMod)
		if err == nil {
			modulePath := modfile.ModulePath(contents)
			if modulePath == "" {
				return types.GoModule{}, fmt.Errorf("invalid go.mod: %s", goMod)
			}
			return types.GoModule{Path: current, Name: modulePath}, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return types.GoModule{}, fmt.Errorf("read %s: %w", goMod, err)
		}
		parent := filepath.Dir(current)
		if parent == current {
			return types.GoModule{}, errors.New("go.mod not found")
		}
	}
}

type walkOpt struct {
	root       string
	module     types.GoModule
	annotation string
	path       string
	entry      fs.DirEntry
	walkErr    error
}

// testWalk checks [walkOpt] for:
//   - if there's any walking directory error (permission denied or others).
//   - if directory should be skipped
//   - if directory contains "go.mod" file.
//   - if there's searched annotation.
func testWalk(opt walkOpt) (skip bool, err error) {
	if opt.walkErr != nil {
		return true, opt.walkErr
	}

	if !opt.entry.IsDir() {
		return true, nil
	}

	if opt.path != opt.root && shouldSkipDirectory(opt.entry.Name()) {
		return true, filepath.SkipDir
	}

	if opt.path != opt.module.Path {
		// skip sub-modules
		if _, err = helper.Stat(filepath.Join(opt.path, "go.mod")); err == nil {
			return true, filepath.SkipDir
		}

		// return err for any errors except ErrNoExist (file does not exist)
		if !errors.Is(err, os.ErrNotExist) {
			return true, err
		}
	}

	var matches bool
	if matches, err = directoryHasAnnotation(opt.path, opt.annotation); err != nil || !matches {
		return true, err
	}

	return false, nil
}

func walkPackages(root string, module types.GoModule, annotation string, discovered map[string]rawPackage) error {
	return filepath.WalkDir(root, func(location string, entry fs.DirEntry, walkErr error) (err error) {
		var skip bool
		if skip, err = testWalk(walkOpt{
			root:       root,
			module:     module,
			annotation: annotation,
			path:       location,
			entry:      entry,
			walkErr:    walkErr,
		}); err != nil || skip {
			return
		}

		var relative string
		if relative, err = filepath.Rel(module.Path, location); err == nil {
			packagePath := module.Path
			if relative != "." {
				packagePath += "/" + filepath.ToSlash(relative)
			}

			discovered[packagePath] = rawPackage{
				Name:    module.Name,
				Path:    location,
				ModPath: module.Path,
				Module:  path.Join(module.Name, relative),
				Related: relative,
			}
		}

		return
	})
}

// shouldSkipDirectory ignores directories designed not for go modules.
func shouldSkipDirectory(name string) bool {
	return name == "vendor" || name == "testdata" || strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_")
}

// shouldSkipSourceFile reports whether an entry should be excluded from package discovery.
func shouldSkipSourceFile(entry os.DirEntry) (should bool) {
	name := entry.Name()
	should =
		entry.IsDir() || // ignore directories
			// ignore files without .go extension
			!strings.HasSuffix(name, ".go") ||
			// ignore test files
			strings.HasSuffix(name, "_test.go") ||
			// ignore explicit ignored files, e.g. _ignored.go or .ignored.go
			// the same is for _doc.go and .doc.go
			strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".")

	return
}

func directoryHasAnnotation(directory, annotation string) (bool, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return false, err
	}

	hasPackage := false
	hasAnnotation := false
	for _, entry := range entries {
		if shouldSkipSourceFile(entry) {
			continue
		}
		filename := filepath.Join(directory, entry.Name())
		file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
		if err != nil {
			return false, fmt.Errorf("parse %s: %w", filename, err)
		}
		hasPackage = file.Name != nil
		if commentsContainAnnotation(file.Comments, annotation) {
			hasAnnotation = true
		}
	}
	return hasPackage && hasAnnotation, nil
}

func commentsContainAnnotation(groups []*ast.CommentGroup, annotation string) bool {
	for _, group := range groups {
		for _, comment := range group.List {
			for line := range strings.SplitSeq(comment.Text, "\n") {
				if strings.HasPrefix(normalizeAnnotation(line), annotation) {
					return true
				}
			}
		}
	}
	return false
}

func normalizeAnnotation(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "//")
	value = strings.TrimPrefix(value, "/*")
	value = strings.TrimSuffix(value, "*/")
	value = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(value), "*"))
	value = strings.TrimPrefix(value, "+")
	return strings.TrimSpace(value)
}
