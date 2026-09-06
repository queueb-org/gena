package tools

import (
	"github.com/spf13/cobra"

	"queueb.org/gena/pkg/common"
)

const (
	// EnvPrefix for tools sub-commands.
	EnvPrefix = "TOOLS"
)

// New creates sub-commands in a form of [cobra.Command].
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "tools <sub-command>",
		Aliases:       []string{"t", "utils"},
		Short:         "tooling sub-commands",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       exampleTools,
		GroupID:       common.GroupAux,
	}

	// sub-commands
	cmd.AddCommand(
		NewDiscoverCommand().Register(),
	)
	return cmd
}
