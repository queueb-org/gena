package tools

import "testing"

func TestNew(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		if cmd := New(); cmd == nil {
			in.Errorf("New() expected return not a nil pointer.")
		}
	})
}
