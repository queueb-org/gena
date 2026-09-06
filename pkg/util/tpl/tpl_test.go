package tpl

import (
	"testing"

	"queueb.org/gena/pkg/util/discover"
	"queueb.org/gena/pkg/util/test"
)

const (
	defaultAnnotation = "+k8s:deepcopy-gen=package"
)

func TestRender(t *testing.T) {
	dir := test.WithTestModule(t)
	apps, err := discover.CollectApps([]string{dir}, defaultAnnotation)
	if err != nil {
		t.Fatalf("could not collect apps: %v", err)
	}
	ktx := CollectContext(apps[0])

	for _, test := range []struct {
		name     string
		setup    func() (string, *Context)
		expected string
	}{
		{
			name: "ok",
			setup: func() (string, *Context) {
				return "{{ .App }}", ktx
			},
			expected: "example.org/application",
		},
		{
			name: "invalid-template",
			setup: func() (string, *Context) {
				return "{{ printf ()}}", ktx
			},
			expected: "",
		},
		{
			name: "invalid-render",
			setup: func() (string, *Context) {
				return "{{ .Fatal }}", ktx
			},
			expected: "",
		},
	} {
		t.Run(test.name, func(in *testing.T) {
			tpl, v := test.setup()
			if result := Render(tpl, v); result != test.expected {
				in.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
