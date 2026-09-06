package gen

import (
	"os"
	"os/exec"

	"common.queueb.org/alog"
	"queueb.org/gena/pkg/util/types"
)

var (
	logger = alog.Default()
)

// Generator keeps a simple interface to run code-generator with
// different specified options collected from environment automatically
// and overridden by a user.
type Generator struct {
	// Name is generator name, e.g. deepcopy-gen
	Name string
	// Path where generator binary is located.
	Path string
	// Packages contains collected packages info.
	Packages types.Packages
	// Options accessible by a user.
	Options []GeneratorOption
}

// New creates generator.
func New(executable string, packages types.Packages, opts ...GeneratorOption) *Generator {
	gen := &Generator{
		Name:     executable,
		Path:     executable,
		Packages: packages,
	}

	var options []GeneratorOption
	options = append(options, defaultOptions...)
	options = append(options, opts...)
	gen.Options = options
	return gen
}

// Exec runs generator, returns error if any met.
func (g *Generator) Exec() (err error) {
	args := g.Packages.Collect()
	cmd := exec.Command(g.Path, args...)
	for _, o := range g.Options {
		o(cmd)
	}

	wd, _ := os.Getwd()
	logger.With("exec", cmd.Args, "working-dir", wd).Debug("running generator")
	// execute
	return cmd.Run()
}
