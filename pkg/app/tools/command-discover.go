package tools

import (
	"errors"
	"fmt"
	"slices"

	"github.com/spf13/cobra"

	"queueb.org/gena/pkg/common"
	"queueb.org/gena/pkg/common/cli"
	"queueb.org/gena/pkg/util/discover"
	"queueb.org/gena/pkg/util/types"
)

var _ common.Apper = &DiscoverCommand{}

// DiscoverCommand provides a way to discover Go packages.
type DiscoverCommand struct {
	*cli.IO

	// flags
	annotation string
	output     string
	format     string

	// internals
	args []string
	cmd  *cobra.Command
}

// Register returns configured [cobra.Command] configured for [DiscoverCommand]
// for its further use.
func (c *DiscoverCommand) Register() *cobra.Command {
	return c.cmd
}

// NewDiscoverCommand returns configured [DiscoverCommand].
func NewDiscoverCommand() common.Apper {
	command := &DiscoverCommand{}
	IO := &cli.IO{
		Out:    &command.output,
		Format: &command.format,
	}
	command.IO = IO

	cmd := &cobra.Command{
		Use:           "discover [directory ...]",
		Short:         "discovers annotated Go packages",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       exampleDiscover,
		Args:          cobra.ArbitraryArgs,
		PreRunE:       cli.PipelineRunE(command.preRun),
		RunE:          command.Run,
	}
	command.cmd = cmd
	command.AddFlags(cmd.Flags(), common.AppEnv, EnvPrefix)
	return command
}

func (c *DiscoverCommand) print(apps []*types.App) (err error) {
	if common.Deref(c.Format, "") != "" {
		return c.IO.Print(apps)
	}

	for _, app := range apps {
		for _, p := range app.Packages {
			if _, err = fmt.Fprintln(c.Output(), p.Module); err != nil {
				return
			}
		}
	}

	return
}

// preRun sets (validates) configuration before main entrypoint is invoked.
func (c *DiscoverCommand) preRun(cmd *cobra.Command, args []string) (err error) {
	// Since we are opening resources, if some of them haven't been opened properly
	// we should try to close it and return an error back.
	if err = c.IO.PreRunE(cmd, args); err != nil {
		err = errors.Join(err, c.IO.Close())
		return
	}

	c.args = slices.Clone(args)
	if len(args) == 0 {
		c.args = []string{"."}
	}

	return
}

// Run is the main entrypoint for [DiscoverCommand].
func (c *DiscoverCommand) Run(_ *cobra.Command, args []string) (err error) {
	defer func() {
		err = errors.Join(err, c.IO.Close())
	}()

	var apps []*types.App

	if apps, err = discover.CollectApps(c.args, c.annotation); err == nil {
		if err = c.print(apps); err != nil {
			return
		}
	}

	return
}
