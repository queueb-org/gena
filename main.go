package main

// Note, do not `go run main.go` since it would hung due to app version is being
// read automatically, which require additional options for building application.
// It's better go use `go install queueb.org/gena`.
// `go run .` would work

import (
	// embed is imported to bind static resources into application.
	_ "embed"

	"log/slog"
	"os"
	"runtime"

	"common.queueb.org/alog"

	"queueb.org/gena/pkg/app"
)

var (
	logger = alog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   false,
		Level:       alog.ParseVerbosity(),
		ReplaceAttr: alog.DefaultReplaceAttr,
	}))

	// exit function
	exit func(code int) = os.Exit
)

func init() {
	alog.ReplaceDefault(logger)
	logger.With("version", app.Version(), "arch",
		runtime.GOARCH,
		"os", runtime.GOOS).Debug("starting ...")

}

func main() {
	app := app.NewApp()
	if err := app.Execute(); err != nil {
		logger.With("err", err).Error("exited ..")
		exit(1)
	}
}
