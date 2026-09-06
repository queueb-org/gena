package types

import (
	"fmt"
	"sort"
	"strings"
)

const (
	maxLineSize = 75
)

type Apps []*App

// String represents applications for simple user output.
func (a Apps) String() string {
	b := &strings.Builder{}
	size := len(a)
	for idx, app := range a {
		b.WriteString(fmt.Sprintf("[%d] %s", len(app.Packages), app.Name))

		if b.Len() > maxLineSize {
			_, _ = b.WriteString(" ...")
			break
		}

		if idx < size-1 {
			b.WriteString("; ")
		}
	}
	return b.String()
}

// App represents an application or a library with
// packages where desired packages were located.
type App struct {
	// Name keeps application's name
	// e.g. example.fqdn/my-app
	Name string
	// Path is where application sources located at filesystem,
	// e.g. /home/user/go/src/github.com/user/my-app
	Path string
	// Packages keeps application packages.
	Packages Packages
}

// Packages represents a slice of collected [Package] entries.
type Packages []Package

// Add adds [Package].
func (p *Packages) Add(src ...Package) {
	*p = append(*p, src...)
}

// Collect returns packages in plain form.
// Applies sort by default.
func (p *Packages) Collect() (out []string) {
	for _, item := range *p {
		out = append(out, item.Module)
	}

	sort.Strings(out)
	return
}

// Related prints related packages locations.
func (p *Packages) Related() (out []string) {
	for _, item := range *p {
		out = append(out, item.Related)
	}
	return
}

// GoModule is a simple go module meta representation.
type GoModule struct {
	Name string
	Path string
}

// Package contains
type Package struct {
	Module  string
	Related string
	Path    string
}
