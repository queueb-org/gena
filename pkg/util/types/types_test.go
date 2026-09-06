package types

import (
	"reflect"
	"slices"
	"testing"
)

var (
	testPackage = Package{
		Module:  "example.org/app/pkg/config",
		Related: "pkg/config",
		Path:    "/tmp/directory/app/pkg/config",
	}
	testPackages = Packages{
		testPackage,
		Package{
			Module:  "example.org/app/pkg/config/v1alpha1",
			Related: "pkg/config/v1alpha1",
			Path:    "/tmp/directory/app/pkg/config/v1alpha1",
		},
	}
)

func TestApps_String(t *testing.T) {
	// code coverage test only.
	t.Run("ok", func(in *testing.T) {
		app := &App{
			Name:     "example.org/app",
			Path:     "/home/user/go/src/github.com/user/app",
			Packages: testPackages,
		}
		_ = Apps{app, app, app, app, app}.String()
	})
}

func TestPackages_Add(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		pkgs := Packages{}
		pkgs.Add(testPackages...)

		if !reflect.DeepEqual(pkgs, testPackages) {
			in.Errorf("\nexp:\n%#v\ngot:\n%#v\n", testPackages, pkgs)
		}
	})
}

func TestPackages_Collect(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		expected := []string{
			"example.org/app/pkg/config",
			"example.org/app/pkg/config/v1alpha1",
		}
		if result := testPackages.Collect(); !slices.Equal(result, expected) {
			in.Errorf("expected: %v, got: %v", expected, result)
		}
	})
}

// Related is used for client-gen package collecting.
func TestPackages_Related(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		expected := []string{"pkg/config", "pkg/config/v1alpha1"}
		result := testPackages.Related()
		if !slices.Equal(result, expected) {
			in.Errorf("expected: %#v, got: %#v", expected, result)
		}
	})
}
