package tools

import (
	"bytes"
	"errors"
	"path"
	"slices"
	"testing"

	"queueb.org/gena/pkg/util/test"
	"queueb.org/gena/pkg/util/types"
)

var (
	testBlankApps = types.Apps{
		&types.App{},
	}
)

// helpers
type fakeWriter struct {
	err error
}

// Write provides [io.Writer] interface.
func (f *fakeWriter) Write(buf []byte) (n int, err error) {
	return 0, f.err
}

// NewTestDiscoverCommand creates [DiscoverCommand] applicable for testing.
func NewTestDiscoverCommand(t testing.TB, opts ...*TestOption) *DiscoverCommand {
	opt := mergeTestOptions(t, opts...)
	cmd := NewDiscoverCommand().(*DiscoverCommand)
	cmd.IO.SetOutput(opt.Out, false) // stdout
	return cmd
}

// tests

func TestDiscoverCommand_Run(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		root := test.WithTestModule(in)
		output := &bytes.Buffer{}
		command := NewTestDiscoverCommand(in, TOptWithOut(output))
		command.annotation = AnnotationDefaultValue
		command.args = []string{root}

		if err := command.Run(command.cmd, []string{root}); err != nil {
			in.Errorf("got error: %v", err)
		}
		got, want := output.String(), "example.org/application/pkg/config\nexample.org/application/pkg/config/v1alpha1\n"
		if got != want {
			in.Errorf("got output %q, want %q", got, want)
		}
	})

	t.Run("fail-to-write", func(in *testing.T) {
		root := test.WithTestModule(in)
		expectedErr := errors.New("any write issue")
		output := &fakeWriter{err: expectedErr}
		command := NewTestDiscoverCommand(in, TOptWithOut(output))
		command.annotation = AnnotationDefaultValue
		command.args = []string{root}

		err := command.Run(command.cmd, []string{root})
		if !errors.Is(err, expectedErr) {
			in.Errorf("expected: %v, got: %v", expectedErr, err)
		}
	})
}

func TestDiscoverCommand_preRun(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		command := NewTestDiscoverCommand(t)
		if err := command.preRun(nil, nil); err != nil {
			in.Errorf("got error: %v", err)
		}

		want := []string{"."}
		if !slices.Equal(want, command.args) {
			in.Errorf("want: %#v, got: %#v", want, command.args)
		}
	})

	t.Run("io-init-err", func(in *testing.T) {
		command := NewTestDiscoverCommand(in)
		dir := in.TempDir()

		command.IO.In = new(path.Join(dir, "non-existent-in"))
		if err := command.preRun(nil, nil); err == nil {
			in.Errorf("expected error but got nil instead")
		}
	})
}

func TestDiscoverCommand_print(t *testing.T) {
	t.Run("set-format", func(in *testing.T) {
		cmd := NewTestDiscoverCommand(in)
		cmd.Format = new("json .")
		if err := cmd.print(testBlankApps); err != nil {
			in.Errorf("got error: %v", err)
		}
	})
}
