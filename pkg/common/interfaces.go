package common

import "github.com/spf13/cobra"

// Apper is an interface for [cobra.Command] obects wrapped
// over superior struct.
type Apper interface {
	Register() *cobra.Command
}
