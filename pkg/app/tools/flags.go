package tools

import (
	"github.com/spf13/pflag"

	"common.queueb.org/tools"
)

const (
	AnnotationFlag         = "annotation"
	AnnotationShortFlag    = "a"
	AnnotationDefaultValue = "+k8s:deepcopy-gen=package"
)

// AddFlags adds CLI flag to the application.
func (c *DiscoverCommand) AddFlags(flags *pflag.FlagSet, envParentPrefixes ...string) {
	env := tools.MakeEnv(envParentPrefixes...)

	c.IO.AddFlags(flags, envParentPrefixes...)

	flags.StringVarP(&c.annotation, AnnotationFlag, AnnotationShortFlag,
		tools.EnvP(env, AnnotationFlag, AnnotationDefaultValue),
		tools.EnvUsageP(env, AnnotationFlag, "package annotation to discover"),
	)
}
