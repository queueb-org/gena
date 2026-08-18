package commandrun

import (
	"common.queueb.org/tools"
	"github.com/spf13/pflag"
)

// Version command CLI arguments.
const (
	DiscoverFlag      = "discover"
	DiscoverShortFlag = "d"

	AllFlag         = "all"
	AllShortFlag    = "a"
	AllDefaultValue = false

	ChangeDirFlag         = "change-dir"
	ChangeDirShortFlag    = ""
	ChangeDirDefaultValue = true
)

// Defaults
var (
	DiscoverDefaultValue = []string{"."}
)

// AddFlags adds flags
func (c *Command) AddFlags(flags *pflag.FlagSet, envParentPrefixes ...string) {
	env := tools.MakeEnv(envParentPrefixes...)

	flags.StringSliceVarP(&c.discover, DiscoverFlag, DiscoverShortFlag,
		tools.EnvP(env, DiscoverFlag, DiscoverDefaultValue),
		tools.EnvUsageP(env, DiscoverFlag, "directories to discover suitable packages."),
	)

	flags.BoolVarP(&c.all, AllFlag, AllShortFlag,
		tools.EnvP(env, AllFlag, AllDefaultValue),
		tools.EnvUsageP(env, AllFlag, "use all generators instead specific ones given through arguments."),
	)

	flags.BoolVarP(&c.changeDir, ChangeDirFlag, ChangeDirShortFlag,
		tools.EnvP(env, ChangeDirFlag, ChangeDirDefaultValue),
		tools.EnvUsageP(env, ChangeDirFlag, "change where application is located before running generators."),
	)
}
