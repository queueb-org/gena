package common

import (
	"reflect"
	"testing"
)

func TestDeref(t *testing.T) {
	for _, test := range []struct {
		name     string
		in       *string
		fallback string
		expected string
	}{
		{"ok", new("test"), "fallback", "test"},
		{"nil", nil, "fallback", "fallback"},
	} {
		t.Run(test.name, func(in *testing.T) {
			if result := Deref(test.in, test.fallback); !reflect.DeepEqual(result, test.expected) {
				in.Errorf("expected: %#v, got: %#v", test.expected, result)
			}
		})
	}
}
