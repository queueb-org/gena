package common

import (
	"github.com/spf13/pflag"
)

var (
	verbosity int
)

// Exportable CLI constants.
const (
	VerbosityFlag      = "verbosity"
	VerbosityShortFlag = "v"
)

// GlobalFlags registers global application flags.
func GlobalFlags(flags *pflag.FlagSet) {
	flags.IntVarP(&verbosity, VerbosityFlag, VerbosityShortFlag, verbosity,
		"sets application verbosity (e.g. -v=4)")
}

// Verbosity returns set verbosity level.
func Verbosity() int {
	return verbosity
}
