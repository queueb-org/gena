package tools

import (
	"io"
	"os"
	"testing"
)

// TestOption keeps options to set for testing command objects.
type TestOption struct {
	Out io.Writer
}

func mergeTestOptions(_ testing.TB, opts ...*TestOption) *TestOption {
	def := &TestOption{
		Out: os.Stdout,
	}

	for _, opt := range opts {
		if opt.Out != nil {
			def.Out = opt.Out
		}
	}

	return def
}

// TOptWithOut sets [TestOption.Out] as [TestOption] object.
func TOptWithOut(out io.Writer) *TestOption {
	return &TestOption{Out: out}
}
