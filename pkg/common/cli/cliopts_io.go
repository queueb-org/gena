package cli

import (
	"errors"
	"io"
	"os"
	"sync"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"common.queueb.org/tools"

	"queueb.org/gena/pkg/common"
)

const (
	// defaultFilePerm is a default permission which is used once file created.
	defaultFilePerm = 0o640
	// defaultWriteMode is a mode to open a file for writing.
	defaultWriteMode = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	// DefaultFormat keeps default golang format, it represent input value
	// as json string.
	DefaultFormat = "{{ json . }}"
)

// IO provides Input & Output Options for Commander apps (and sub-apps).
type IO struct {
	Out    *string
	In     *string
	Format *string

	mu  sync.RWMutex
	out io.Writer
	in  io.Reader

	// close flags keeps if we need to Close out and in
	// with [IO.Close] call. Note, os.Stdin and os.Stdout
	// should not be closed.
	closeOutput bool
	closeInput  bool
}

// AddFlags initializes [IO] cli-arguments.
func (c *IO) AddFlags(flags *pflag.FlagSet, envParentPrefixes ...string) {
	env := tools.MakeEnv(envParentPrefixes...)

	if c.Out != nil {
		flags.StringVarP(c.Out, OutputFlag, OutputShortFlag,
			tools.EnvP(env, OutputFlag, OutputDefaultValue),
			tools.EnvUsageP(env, OutputFlag, "command output, use blank for stdout."),
		)
	}

	if c.In != nil {
		flags.StringVarP(c.In, InputFlag, InputShortFlag,
			tools.EnvP(env, InputFlag, InputDefaultValue),
			tools.EnvUsageP(env, InputFlag, "command input (location), blank is treated as stdin."),
		)
	}

	if c.Format != nil {
		flags.StringVarP(c.Format, FormatFlag, FormatShortFlag,
			tools.EnvP(env, FormatFlag, FormatDefaultValue),
			tools.EnvUsageP(env, FormatFlag, "go template format to output values"),
		)
	}
}

// PreRunE should be given through cobra's command [cobra.Command.PreRunE]
// PreRunE initializes and validates input and output arguments.
func (c *IO) PreRunE(_ *cobra.Command, args []string) (err error) {
	errs := &common.ErrStack{}
	if c.In != nil && *c.In != "" {
		c.closeInput = true
		if c.in, err = os.Open(*c.In); err != nil {
			// Do not store (*os.File)(nil) pointer to avoid
			// Close operations over interface conversions over [IO.Close].
			c.in = nil
		}
		errs.Add(err)
	}

	if c.Out != nil && *c.Out != "" {
		c.closeOutput = true
		if c.out, err = os.OpenFile(*c.Out, defaultWriteMode, defaultFilePerm); err != nil {
			// Same as for c.in
			c.out = nil
		}
		errs.Add(err)
	}

	return errs.Collect()
}

// Close resources.
func (c *IO) Close() error {
	var err error
	if c.closeInput {
		if v, ok := c.in.(io.Closer); ok {
			err = v.Close()
		}
	}

	if c.closeOutput {
		if v, ok := c.out.(io.Closer); ok {
			err = errors.Join(err, v.Close())
		}
	}

	return err
}

// Output returns configured [IO] out writer object.
func (c *IO) Output() io.Writer {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.out == nil {
		c.out = os.Stdout
		c.closeOutput = false
	}

	return c.out
}

// Input returns configured [IO] input reader object.
func (c *IO) Input() io.Reader {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.in == nil {
		c.in = io.NopCloser(os.Stdin)
		c.closeInput = false
	}

	return c.in
}

// SetOutput sets [IO] writer interface (i.e. [IO.Output]).

func (c *IO) SetOutput(w io.Writer, shouldClose bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.out = w
	c.closeOutput = shouldClose
}

// SetInput sets [IO] reader interface (i.e. [IO.Input]).
func (c *IO) SetInput(r io.Reader, shouldClose bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.in = r
	c.closeInput = shouldClose
}

// Print formats input value.
// If format is blank, uses [DefaultFormat] instead.
func (c *IO) Print(v any, opts ...*IOption) (err error) {
	opt := mergeIOptions(opts...)
	return c.render(v, opt)
}

func (c *IO) render(v any, opt *IOption) (err error) {
	if c.Format == nil {
		return
	}

	var tpl *template.Template
	if tpl, err = template.New("app-format").
		Funcs(FuncMap).     // default func maps
		Funcs(opt.FuncMap). // options
		Parse(common.Or(*c.Format, DefaultFormat)); err == nil {

		err = tpl.Execute(c.Output(), v)
	}

	return
}
