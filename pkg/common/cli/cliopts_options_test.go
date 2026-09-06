package cli

import (
	"reflect"
	"testing"
)

func TestMergeIOptions(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		opt := mergeIOptions(&IOption{
			FuncMap: FuncMap,
		})

		expected := &IOption{
			FuncMap: FuncMap,
		}

		if !reflect.DeepEqual(opt, expected) {
			in.Errorf("\nexp:\n%#v\ngot:\n%#v\n", expected, opt)
		}
	})
}
