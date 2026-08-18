package common

import (
	"errors"
	"testing"
)

func TestErrStack(t *testing.T) {
	errExpected := errors.New("any error")

	for _, test := range []struct {
		name     string
		setup    func(t testing.TB) error
		expected error
	}{
		{
			name: "ok",
			setup: func(t testing.TB) error {
				errs := &ErrStack{}
				errs.Add(errExpected)
				return errs.Collect()
			},
			expected: errExpected,
		},
		{
			name: "nil",
			setup: func(t testing.TB) error {
				return (&ErrStack{}).Collect()
			},
			expected: nil,
		},
	} {
		t.Run(test.name, func(in *testing.T) {
			if err := test.setup(in); !errors.Is(err, test.expected) {
				in.Errorf("\nexp: %#v\ngot: %#v\n", test.expected, err)
			}
		})
	}

	t.Run("Error()", func(in *testing.T) {
		errs := &ErrStack{}
		errs.Add(errExpected, errExpected)
		expected := "[0] any error; [1] any error"
		if result := errs.Collect().Error(); result != expected {
			in.Errorf("\nexp: %v\ngot: %v\n", expected, result)
		}
	})
}

func TestErrorString(t *testing.T) {
	for _, test := range []struct {
		name     string
		in       error
		expected string
	}{
		{
			name:     "nil",
			in:       nil,
			expected: "",
		},
		{
			name:     "plain",
			in:       errors.New("test"),
			expected: "test",
		},
	} {
		t.Run(test.name, func(in *testing.T) {
			if result := ErrorString(test.in); result != test.expected {
				in.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
