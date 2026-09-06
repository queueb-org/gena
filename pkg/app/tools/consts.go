package tools

import (
	// embed is imported to bind static resources into the app.
	_ "embed"
)

var (
	//go:embed "resources/command.example.txt"
	exampleTools string
	//go:embed "resources/discover.example.txt"
	exampleDiscover string
)
