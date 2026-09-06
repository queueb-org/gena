package common

import "github.com/spf13/cobra"

// Exportable constants.
const (
	AppName = "gena"
	AppEnv  = "GENA"

	// Cobra related constants.

	GroupMain = "Main"
	GroupAux  = "Aux"
)

var Groups = []*cobra.Group{
	{
		ID:    GroupMain,
		Title: "Main:",
	},
	{
		ID:    GroupAux,
		Title: "Aux:",
	},
}
