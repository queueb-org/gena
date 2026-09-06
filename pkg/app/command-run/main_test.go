package commandrun

import (
	"slices"
	"testing"

	"queueb.org/gena/pkg/util/gen"
	"queueb.org/gena/pkg/util/test"
	"queueb.org/gena/pkg/util/types"
)

func TestCommand_Register(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		cmd := NewTestCommand(in)
		if c := cmd.Register(); c == nil {
			in.Errorf("cobra command expected not to be nil")
		}
	})
}

// NewTestCommand configures [Command] applicable for tests.
func NewTestCommand(t testing.TB, opts ...*TestOption) *Command {
	opt := mergeTestOptions(t, opts...)
	cmd := New()
	cmd.discover = []string{opt.Dir}

	if err := cmd.preRun(nil, opt.Args); err != nil {
		t.Fatalf("could not run preRun: %v", err)
	}

	return cmd
}

func TestCommand_Run(t *testing.T) {
	dir := test.WithTestModule(t)

	t.Run("ok", func(in *testing.T) {
		cmd := NewTestCommand(in, TOptDir(dir), TOptDeepcopyGen())
		if err := cmd.Run(nil, nil); err != nil {
			in.Errorf("got error: %v", err)
		}
	})

	// there's no data related selected generator.
	t.Run("no-gen-set", func(in *testing.T) {
		cmd := NewTestCommand(in, TOptDir(dir), TOptDefaulterGen())
		if err := cmd.Run(nil, nil); err != nil {
			in.Errorf("got error: %v", err)
		}
	})

	t.Run("cant-discover", func(in *testing.T) {
		dir := test.WithInvalidTestModule(in)
		cmd := NewTestCommand(in, TOptDir(dir), TOptDeepcopyGen())

		if err := cmd.Run(nil, nil); err == nil {
			in.Errorf("expected error but got nil instead")
		}
	})

	t.Run("cant-run-generator", func(in *testing.T) {
		cmd := NewTestCommand(in, TOptDir(dir), TOptDeepcopyGen())
		cmd.makeGenFunc = func(s string, p types.Packages, opts ...gen.GeneratorOption) *gen.Generator {
			return &gen.Generator{
				Name: "fake",
				Path: dir,
			}
		}
		if err := cmd.Run(nil, nil); err == nil {
			in.Errorf("error expected but got nil instead")
		}
	})
}

func TestCommand_chDir(t *testing.T) {
	// code coverage test only
	t.Run("no-ok", func(in *testing.T) {
		cmd := NewTestCommand(in)
		cmd.chDir("non existent")
	})
}

func TestCommand_lookupAnnotation(t *testing.T) {
	cmd := NewTestCommand(t)

	for _, test := range []struct {
		name     string
		in       string
		expected string
	}{
		{"ok", "deepcopy-gen", genToAnnotation["deepcopy-gen"]},
		{"not-found", "custom-deepcopy-gen", genToAnnotation[""]},
	} {
		t.Run(test.name, func(in *testing.T) {
			if result := cmd.lookupAnnotation(test.in); result != test.expected {
				in.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}

func TestCommand_fillGens(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		cmd := NewTestCommand(in)
		cmd.all = true
		cmd.fillGens(nil)
		expected := []string{
			"deepcopy-gen",
			"defaulter-gen",
			"conversion-gen",
			"validation-gen",
			"register-gen",
		}
		if result := cmd.gens; !slices.Equal(expected, result) {
			in.Errorf("\nexp: %#v\ngot: %#v\n", expected, result)
		}
	})
}
