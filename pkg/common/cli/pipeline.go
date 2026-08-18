package cli

import "github.com/spf13/cobra"

// RunFuncE represents cobra RunE related functions, like
// PreRunE, PersistentPreRunE, RunE and so on.
type RunFuncE = func(cmd *cobra.Command, args []string) (err error)

// PipelineRunE builds a chain of [cobra.Command] RunE, PreRunE, etc functions
// into a single one run in a form of a pipeline.
// Example:
//
//	cmd := &cobra.Command{
//		PreRunE: PipelineRunE(func1, func2)
//	}
func PipelineRunE(funcs ...RunFuncE) RunFuncE {
	return func(cmd *cobra.Command, args []string) (err error) {
		for _, fnc := range funcs {
			// nil function is valid, however, such should be skipped.
			if fnc == nil {
				continue
			}

			if err = fnc(cmd, args); err != nil {
				return
			}
		}

		return nil
	}
}
