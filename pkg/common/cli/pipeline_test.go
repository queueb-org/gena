package cli

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

var noOpRunE RunFuncE = func(cmd *cobra.Command, _ []string) (err error) {
	return nil
}

func TestPipelineRunE(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		var container string
		fnc := func(cmd *cobra.Command, _ []string) (err error) {
			container = "here some result"
			return nil
		}

		if err := PipelineRunE(nil, noOpRunE, noOpRunE, fnc)(nil, nil); err != nil {
			in.Fatalf("got error: %v", err)
		}

		expected := "here some result"
		if result := container; result != expected {
			in.Errorf("expected: %v, got: %v", expected, result)
		}
	})

	t.Run("on-err", func(in *testing.T) {
		expectedErr := errors.New("any issue")
		fnc := func(cmd *cobra.Command, _ []string) (err error) {
			return expectedErr
		}
		if err := PipelineRunE(noOpRunE, fnc)(nil, nil); !errors.Is(err, expectedErr) {
			in.Errorf("expected: %v, got: %v", expectedErr, err)
		}
	})
}
