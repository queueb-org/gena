package common

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestGlobalFlags(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		flags := &pflag.FlagSet{}
		GlobalFlags(flags)

		if _, err := flags.GetInt(VerbosityFlag); err != nil {
			in.Errorf("got error: %v", err)
		}
	})
}

// code coverage test only, if theres's something to check => feel free to add it.
func TestVerbosity(t *testing.T) {
	Verbosity()
}
