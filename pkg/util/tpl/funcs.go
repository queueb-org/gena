package tpl

import (
	"strings"
	"text/template"
)

var (
	// FuncMap keeps default template FuncMap to evaluate text go templates.
	FuncMap = template.FuncMap{
		"join": strings.Join,
	}
)
