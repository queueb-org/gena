package tpl

import (
	"bytes"
	"text/template"

	"common.queueb.org/alog"
	"queueb.org/gena/pkg/util/types"
)

var (
	// default package logger.
	logger = alog.Default()
)

// Context is a template render Context object.
type Context struct {
	// Packages keeps collected packages.
	Packages types.Packages
	// ProjectDir contains project directory full path.
	// [Context.ProjectDir] is discovered by [ProjectDir] helper.
	ProjectDir string
	// App represents application
	// e.g. example.fqdn/my-app
	App string
	// AppDir contains base application directory
	// e.g. /home/user/go/src/github.com/user/my-app
	AppDir string
}

// CollectContext returns configured [context].
func CollectContext(app *types.App) *Context {
	projectDir := ProjectDir(app.Path)

	return &Context{
		App:        app.Name,
		AppDir:     app.Path,
		Packages:   app.Packages,
		ProjectDir: projectDir,
	}
}

// Render renders input with given context.
func Render(in string, ktx *Context) string {
	tpl, err := template.New("render-value").Funcs(FuncMap).Parse(in)
	if err != nil {
		logger.With("template", in, "err", err).Warn("could not parse template, using blank value")
		return ""
	}
	buf := &bytes.Buffer{}
	if err = tpl.Execute(buf, ktx); err != nil {
		logger.With("template", in, "err", err).Warn("could execute template")
		return ""
	}

	return buf.String()
}
