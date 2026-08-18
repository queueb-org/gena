package cli

import (
	"encoding/json"
	"testing"
)

func TestToes(t *testing.T) {
	for _, test := range []struct {
		name     string
		in       any
		call     func(any) string
		expected string
	}{
		{"yaml-ok", 1337, toYAML, "1337\n"},
		{"yaml-ok", 1337, toJSON, "1337"},
	} {
		t.Run(test.name, func(in *testing.T) {
			if result := test.call(test.in); result != test.expected {
				in.Errorf("\nexp:\n`%v`\ngot:\n`%v`\n", test.expected, result)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		if result := marshal(1337, json.Marshal); result != "1337" {
			in.Errorf("expected: 1337, got: %v", result)
		}
	})

	t.Run("err", func(in *testing.T) {
		ch := make(chan int, 2)
		expected := "json: unsupported type: chan int"

		if result := marshal(ch, json.Marshal); result != expected {
			in.Errorf("expected: %v, got: %v", expected, result)
		}
	})
}
