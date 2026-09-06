package app

import (
	"github.com/spf13/pflag"

	"common.queueb.org/tools"
)

// Version command CLI arguments.
const (
	ShortFlag         = "short"
	ShortDefaultValue = false
)

func addVersionFlags(flags *pflag.FlagSet, parentEnvPrefixes ...string) {
	env := tools.MakeEnv(parentEnvPrefixes...)

	flags.BoolVarP(&short, ShortFlag, "",
		tools.Env(env(ShortFlag), ShortDefaultValue),
		tools.EnvUsage(env(ShortFlag), "show short version."),
	)
}
