package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"queueb.org/gena/pkg/common"
)

// Helpers

func NewTestIO(t testing.TB) *IO {
	return &IO{
		Out:    new(""),
		In:     new(""),
		Format: new(""),

		out: &bytes.Buffer{},
		in:  &bytes.Buffer{},
	}
}

func TestIO_AddFlags(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		flags := &pflag.FlagSet{}
		NewTestIO(in).AddFlags(flags, "APP")

		_, err := flags.GetString(OutputFlag)
		if err != nil {
			in.Errorf("got error: %v", err)
		}
	})
}

func TestIO_PreRunE(t *testing.T) {
	for _, test := range []struct {
		name     string
		setup    func(t testing.TB) *IO
		expected error
	}{
		{
			name: "ok",
			setup: func(t testing.TB) *IO {
				return NewTestIO(t)
			},
			expected: nil,
		},
		{
			name: "in-out-not-found",
			setup: func(t testing.TB) *IO {
				dir := t.TempDir()
				cmd := NewTestIO(t)
				cmd.In = new(path.Join(dir, "non-existent-in"))
				cmd.Out = new(path.Join(dir, "non-existent/non-existent-out"))
				return cmd
			},
			expected: fmt.Errorf("no such file or directory"),
		},
	} {
		t.Run(test.name, func(in *testing.T) {
			cmd := test.setup(in)
			err := cmd.PreRunE(nil, nil)
			if !errors.Is(err, test.expected) {
				// check if string matches.
				left, right := common.ErrorString(err), common.ErrorString(test.expected)
				if !strings.Contains(left, right) {
					in.Errorf("\nexp: %v\ngot: %v\n", test.expected, err)
				}
			}
		})
	}
}

func TestIO_Close(t *testing.T) {
	for _, test := range []struct {
		name     string
		setup    func(t testing.TB) *IO
		expected error
	}{
		{
			name: "ok",
			setup: func(t testing.TB) *IO {
				return NewTestIO(t)
			},
			expected: nil,
		},
		{
			name: "close",
			setup: func(t testing.TB) *IO {
				cmd := NewTestIO(t)
				cmd.SetInput(io.NopCloser(&bytes.Buffer{}), true)
				cmd.SetOutput(common.MustNopCloser(&bytes.Buffer{}), true)
				return cmd
			},
			expected: nil,
		},
	} {
		t.Run(test.name, func(in *testing.T) {
			if err := test.setup(in).Close(); !errors.Is(err, test.expected) {
				in.Errorf("\nexp: %v\ngot: %v\n", test.expected, err)
			}
		})
	}
}

func TestIO_Output(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		cmd := NewTestIO(in)
		cmd.SetOutput(nil, false) // reset
		out := cmd.Output()
		if !reflect.DeepEqual(out, os.Stdout) {
			in.Errorf("stdout expected, got: %#v", out)
		}
	})
}

func TestIO_Input(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		cmd := NewTestIO(in)
		cmd.SetInput(nil, false) // reset
		_in := cmd.Input()
		if !reflect.DeepEqual(_in, io.NopCloser(os.Stdin)) {
			in.Errorf("stdout expected, got: %#v", _in)
		}
	})
}

type TestObj struct {
	Name  string
	Value *TestObj
}

var (
	testObj = &TestObj{
		Name: "root",
		Value: &TestObj{
			Name:  "sub",
			Value: nil,
		},
	}
)

func TestIO_Print(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		cmd := NewTestIO(in)
		if err := cmd.Print(testObj); err != nil {
			in.Errorf("got error: %v", err)
		}
	})
}

func TestIO_render(t *testing.T) {
	for _, test := range []struct {
		name  string
		setup func(t testing.TB) (*IO, any)
		err   error
	}{
		{
			name: "format-nil",
			setup: func(t testing.TB) (*IO, any) {
				cmd := NewTestIO(t)
				cmd.Format = nil
				return cmd, testObj
			},
			err: nil,
		},
		{
			name: "format-not-nil",
			setup: func(t testing.TB) (*IO, any) {
				cmd := NewTestIO(t)
				cmd.Format = new("{{ .Name }}")
				return cmd, testObj
			},
			err: nil,
		},
	} {
		t.Run(test.name, func(in *testing.T) {
			cmd, v := test.setup(in)
			opt := mergeIOptions()
			if err := cmd.render(v, opt); !errors.Is(err, test.err) {
				in.Errorf("expected: %v, got: %v", test.err, err)
			}
		})
	}
}
