package gen

import (
	"testing"

	"queueb.org/gena/pkg/util/types"
)

var (
	testPackages = types.Packages{
		types.Package{
			Module:  "example.org/app/pkg/config",
			Related: "pkg/config",
			Path:    "/home/user/go/src/github.com/user/app/pkg/config",
		},
		types.Package{
			Module:  "example.org/app/pkg/config/v1alpha1",
			Related: "pkg/config/v1alpha1",
			Path:    "/home/user/go/src/github.com/user/app/pkg/config/v1alpha1",
		},
	}
)

func TestGenerator_Exec(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		gen := New("echo", testPackages, OptDefaults, OptExtendArgs([]string{"-V=4"}))
		if err := gen.Exec(); err != nil {
			in.Errorf("got error: %v", err)
		}
	})
}
