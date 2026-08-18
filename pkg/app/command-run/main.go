package commandrun

import (
	// embed is used to bind static resources into the app.
	_ "embed"

	"os"
	"strings"

	"common.queueb.org/alog"
	"github.com/spf13/cobra"

	"queueb.org/gena/pkg/common"
	"queueb.org/gena/pkg/util/discover"
	g "queueb.org/gena/pkg/util/gen"
	"queueb.org/gena/pkg/util/tpl"
	"queueb.org/gena/pkg/util/types"
)

var (
	logger = alog.Default()

	envPrefixes = []string{common.AppEnv, "RUN"}

	//go:embed "resources/command.run.txt"
	exampleRun string
)

// Command serves generators running purposes.
type Command struct {
	// gens keeps validated generators taken from user input
	// and configuration files.
	gens []string
	// all replaces gens with pre-set generators list.
	all bool
	// discover contains directories to discover required modules.
	discover []string
	// changeDir allows to change dir each found application before
	// running generators, true by default.
	changeDir bool

	// workingDir keeps application working directory to reverse processes
	// current directory to it if [Command.chDir] was invoked.
	workingDir string

	// makeGenFunc contains creating a new generator function.
	// by default it's set to [g.New]. The override is used mainly for tests.
	makeGenFunc func(string, types.Packages, ...g.GeneratorOption) *g.Generator

	cmd *cobra.Command
}

// Register returns configured [cobra.Command] of [Command].
func (c *Command) Register() *cobra.Command {
	return c.cmd
}

// New configures [Command].
func New() *Command {
	// defaults
	command := &Command{
		makeGenFunc: g.New,
	}

	cmd := &cobra.Command{
		Use:          "run [gen1 gen2 ...genN]",
		Aliases:      []string{"r"},
		SilenceUsage: true,
		Example:      exampleRun,
		GroupID:      common.GroupMain,
		Short:        "runs generators",
		Args:         cobra.ArbitraryArgs,
		PreRunE:      command.preRun,
		RunE:         command.Run,
	}
	command.AddFlags(cmd.Flags(), envPrefixes...)
	command.cmd = cmd
	return command
}

// lookupAnnotation looks for annotation suitable for generator.
// If nothing is found returns default annotation.
func (c *Command) lookupAnnotation(gen string) (annotation string) {
	var found bool
	if annotation, found = genToAnnotation[gen]; !found {
		annotation = genToAnnotation[""]
	}

	return
}

// renderOptions gets options per generator.
func (c *Command) renderOptions(gen string, ktx *tpl.Context) (options []string) {
	if container, found := genToOptions[gen]; found {
		for key, value := range container {
			options = append(options, key, tpl.Render(value, ktx))
		}
	}
	return
}

func (c *Command) fillGens(args []string) {
	if c.all {
		c.gens = allGenerators
		return
	}

	// TODO: replace validating generators from user input for something decent
	for _, arg := range args {
		if gen := strings.ToLower(arg); genRx.MatchString(gen) {
			c.gens = append(c.gens, gen)
		}
	}
}

// preRun validates and configures application input user data.
func (c *Command) preRun(_ *cobra.Command, args []string) (err error) {
	// set application working directory, if error sets blank.
	c.workingDir, _ = os.Getwd()
	c.fillGens(args)
	return
}

// chDir performs changing directory.
func (c *Command) chDir(dir string) {
	if c.changeDir {
		if err := os.Chdir(dir); err != nil {
			logger.With("err", err).Warn("could not perform chdir")
		}
	}
}

// Run main [Command] entrypoint.
func (c *Command) Run(_ *cobra.Command, _ []string) (err error) {
	logger.With("gens", c.gens).Info("running generators")

	var apps types.Apps

	for _, gen := range c.gens {
		log := logger.With("gen", gen)
		log.Notice("running")
		logger.With("dirs", c.discover).Debug("discovering dirs ..")

		apps, err = discover.CollectApps(c.discover, c.lookupAnnotation(gen))
		if err != nil {
			log.Error("could not retrieve packages for specified generator, please check if project is")
			return
		}

		if len(apps) == 0 {
			log.With("apps", apps.String()).Notice("not found, skipping ...")
			continue
		}

		// NOTE: clean go caches might cause problems.
		for _, app := range apps {
			// we require to change directory before running generators.
			c.chDir(app.Path)

			opts := c.renderOptions(gen, tpl.CollectContext(app))
			generator := c.makeGenFunc(gen, app.Packages,
				g.OptWithKubeVerbosity(),
				g.OptWithCopyrights(),
				g.OptWithFlagsValues(opts...),
			)
			if err = generator.Exec(); err != nil {
				return
			}

			// return return back
			c.chDir(c.workingDir)
		}
	}

	return
}
