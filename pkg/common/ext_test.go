package common

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"common.queueb.org/tests"
)

func TestNopCloser(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		w, err := NopCloser(&bytes.Buffer{})
		if err != nil {
			in.Fatalf("got error: %v", err)
		}

		// check [nopCloser.Close]
		if err = w.Close(); err != nil {
			in.Fatalf("got error: %v", err)
		}

		// check [nopCloser.Write]
		_, err = io.Copy(w, bytes.NewReader([]byte("this is a string")))
		if err != nil {
			in.Errorf("got error: %v", err)
		}
	})

	t.Run("should-not-be-closer", func(in *testing.T) {
		w, err := NopCloser(&bytes.Buffer{})
		if err != nil {
			in.Fatalf("got error: %v", err)
		}
		if newW, err := NopCloser(w); !errors.Is(err, ErrInvalid) || newW != nil {
			in.Errorf("\nexp: %#v\ngot: %#v\n", ErrInvalid, err)
		}
	})
}

func TestMustNopCloser(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		defer tests.PanicNotExpected(in)()
		MustNopCloser(&bytes.Buffer{})
	})

	t.Run("panic", func(in *testing.T) {
		defer tests.PanicExpected(in)()
		MustNopCloser(MustNopCloser(&bytes.Buffer{}))
	})
}

func TestClose(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		Close(1337, &bytes.Buffer{}, io.NopCloser(&bytes.Buffer{}))
	})
}

func TestOr_string(t *testing.T) {
	for _, test := range []struct {
		name     string
		in       []string
		expected string
	}{
		{"ok", []string{"", "", "one", "", "two", "", "three"}, "one"},
		{"nil", nil, ""},
	} {
		t.Run(test.name, func(in *testing.T) {
			if result := Or(test.in...); result != test.expected {
				in.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
