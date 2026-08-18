package app

import (
	// embed is imported to bind static resources into application.
	_ "embed"
	"errors"

	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	run "queueb.org/gena/pkg/app/command-run"
	"queueb.org/gena/pkg/app/tools"
	"queueb.org/gena/pkg/common"
)

var (
	//go:embed "resources/command.main.txt"
	exampleApp string

	version = Version()

	// flags
	short bool

	// commands
	versionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   "prints application complex version",
		GroupID: common.GroupAux,
		Run: func(_ *cobra.Command, _ []string) {

			if short {
				fmt.Println(version)
				return
			}

			version := fmt.Sprintf(
				"%[1]s version: %[2]s, %[3]s/%[4]s %[5]s",
				common.AppName, version, runtime.GOOS, runtime.GOARCH, runtime.Version())
			fmt.Println(version)
		},
	}

	testCmd = &cobra.Command{
		Use:    "__test",
		Short:  "testing command",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("any issue")
		},
	}
)

func init() {
	// version command
	addVersionFlags(versionCmd.Flags(), common.AppEnv, "VERSION")
}

// NewApp creates gena application instance
func NewApp() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "gena",
		Aliases:       []string{"na"},
		Short:         "code-Generators Assistant",
		Example:       exampleApp,
		SilenceErrors: true,
	}
	common.GlobalFlags(cmd.PersistentFlags())
	cmd.AddGroup(common.Groups...)
	cmd.AddCommand(versionCmd, testCmd,
		run.New().Register(),
		// app sub-groups
		tools.New(),
	)
	return cmd
}
