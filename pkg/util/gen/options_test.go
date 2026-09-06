package gen

import (
	"os/exec"
	"runtime"
	"slices"
	"testing"
)

// NewTestCmd creates test cmd.
func NewTestCmd() *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd.exe", "/D", "/C", "echo", "test")
	default:
		return exec.Command("echo", "test")
	}

}

// WithContainsArgs checks if cmd has sub arguments
func WithContainsArgs(t testing.TB, cmd *exec.Cmd, sub []string) {
	var found int
	compacted := slices.Compact(cmd.Args)

	for _, arg := range compacted {
		for _, subArg := range sub {
			if arg == subArg {
				found++
			}
		}
	}

	if found < len(sub) {
		t.Errorf("no sub arguments found: %v", sub)
	}
}

// code coverage only.
func TestOpDefaults(t *testing.T) {
	OptDefaults(NewTestCmd())
}

func TestOptWithCopyrights(t *testing.T) {
	cmd := NewTestCmd()
	loc := "../../../hack/boilerplate.go.txt"
	t.Setenv(EnvGoHeaderFile, loc)
	OptWithCopyrights()(cmd)
	WithContainsArgs(t, cmd, []string{"--go-header-file=" + loc})
}

func TestOptExtendArgs(t *testing.T) {
	cmd := NewTestCmd()
	OptExtendArgs([]string{"--my-arg=1", "--my-arg=2", "--my-arg=3"})(cmd)
	WithContainsArgs(t, cmd, []string{"--my-arg=1", "--my-arg=2", "--my-arg=3"})
}

func TestOptWithFlagsValues(t *testing.T) {
	cmd := NewTestCmd()
	OptWithFlagsValues("-v", "4", "--some-option")(cmd)
	WithContainsArgs(t, cmd, []string{"-v=4", "--some-option="})
}

func TestOptWithKubeVerbosity(t *testing.T) {
	t.Run("set", func(in *testing.T) {
		cmd := NewTestCmd()
		in.Setenv(EnvKubeVerbose, "5")
		OptWithKubeVerbosity()(cmd)
		WithContainsArgs(in, cmd, []string{"-v=5"})
	})

	t.Run("not-set", func(in *testing.T) {
		cmd := NewTestCmd()
		OptWithKubeVerbosity()(cmd)
		WithContainsArgs(in, cmd, []string{"-v=0"})
	})
}
