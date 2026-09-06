package gen

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
)

// GeneratorOption configures execution command options
type GeneratorOption func(cmd *exec.Cmd)

var defaultOptions = []GeneratorOption{
	OptDefaults,
}

// Options

// OptDefaults sets default generator executable options.
func OptDefaults(cmd *exec.Cmd) {
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
}

// EnvGoHeaderFile is an environment variable which contains
// a path where boilerplate.txt file is placed.
const EnvGoHeaderFile = "GENA_GO_HEADER_FILE"

// OptWithCopyrights sets --go-header-file=<copyrights.txt> option from environment
// Requires GENA_GO_HEADER_FILE environment variables to be set.
func OptWithCopyrights() GeneratorOption {
	return func(cmd *exec.Cmd) {
		if loc, found := os.LookupEnv(EnvGoHeaderFile); found {
			cmd.Args = append(cmd.Args, fmt.Sprintf("--go-header-file=%s", loc))
		}
	}
}

// OptExtendedArgs extends generator arguments with extended.
func OptExtendArgs(extended []string) func(cmd *exec.Cmd) {
	return func(cmd *exec.Cmd) {
		cmd.Args = append(cmd.Args, extended...)
	}
}

// OptWithFlagsValues sets a value over cli-arguments (i.e. flag)
// Note, keyValues contains pairs for flag and its value.
// It can be unbalanced, in this case the latest value is considered as blank
//
//	OptWithFlagsValues("-v", "4", "--my-flag", "")        // valid => -v=4 --my-flag=
//	OptWithFlagsValues("-v", "4", "--use-local-storage")  // valid => -v=4 --use-local-storage=
func OptWithFlagsValues(keyValues ...string) GeneratorOption {
	return func(cmd *exec.Cmd) {
		kv := slices.Clone(keyValues)
		outOfBounds := len(kv)
		for i := 0; i < outOfBounds; i += 2 {
			key := kv[i]
			value := ""
			if i < outOfBounds-1 {
				value = kv[i+1]
			}
			cmd.Args = append(cmd.Args, fmt.Sprintf("%s=%s", key, value))
		}
	}
}

// EnvKubeVerbose contains verbosity flag level.
const EnvKubeVerbose = "KUBE_VERBOSE"

// OptWithKubeVerbosity configures kubernetes
func OptWithKubeVerbosity() GeneratorOption {
	var (
		value string
		found bool
	)
	if value, found = os.LookupEnv(EnvKubeVerbose); !found {
		value = "0"
	}
	return OptWithFlagsValues("-v", value)
}
