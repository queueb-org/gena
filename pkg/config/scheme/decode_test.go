package scheme

import (
	"errors"
	"os"
	"strings"
	"testing"

	"common.queueb.org/tests"
)

const (
	validConfig = `apiVersion: gena.apps.queueb.org/v1alpha1
kind: Config
`
	testValidConfigLocation = "../../../resources/config.example.yaml"
	testValidStatusLocation = "../../../resources/status.yaml"
)

func TestDecode(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		config, err := Decode(tests.WithOpenFile(in, testValidConfigLocation))
		if err != nil {
			in.Fatalf("Decode() error = %v", err)
		}
		if config == nil {
			in.Fatal("Decode() returned a nil config")
		}
	})

	t.Run("cant-read", func(in *testing.T) {
		fd, err := os.Open(testValidConfigLocation)
		if err != nil {
			in.Fatalf("could not read test configuration file: %v", err)
		}

		// explicitly closing
		_ = fd.Close()
		if _, err := Decode(fd); !errors.Is(err, ErrDecode) {
			in.Errorf("expected: `%v`, got: `%v`", ErrDecode, err)
		}
	})

	t.Run("unknown-field", func(in *testing.T) {
		raw := tests.WithReadFileString(in, testValidConfigLocation) + "\nunknown: value\n"

		_, err := Decode(strings.NewReader(raw))
		if err == nil {
			in.Fatal("Decode() error = nil, want a strict decoding error")
		}
		if !strings.Contains(err.Error(), "unknown field \"unknown\"") {
			in.Fatalf("Decode() error = %q, want an unknown field error", err)
		}
	})

	t.Run("unsupported-version", func(in *testing.T) {
		data := `apiVersion: gena.apps.queueb.org/v1beta1
kind: Config`

		_, err := Decode(strings.NewReader(data))
		if err == nil {
			in.Fatal("Decode() error = nil, want a version error")
		}

		if !strings.Contains(err.Error(),
			"no kind \"Config\" is registered for version") {
			in.Errorf("Decode() error = %q, want an unregistered version error", err)
		}
	})

	t.Run("unsupported-api-object", func(in *testing.T) {
		data := `apiVersion: v1
kind: Status
status: Failure`

		_, err := Decode(strings.NewReader(data))
		if err == nil {
			in.Fatal("Decode() error = nil, want a version error")
		}
		if !strings.Contains(err.Error(),
			"could not decode object: unexpected object /v1, Kind=Status") {
			in.Fatalf("Decode() error = %q, want an unregistered version error", err)
		}
	})

	t.Run("duplicate-field", func(in *testing.T) {
		data := validConfig + "kind: Config\n"

		_, err := Decode(strings.NewReader(data))
		if err == nil {
			in.Fatal("Decode() error = nil, want a duplicate field error")
		}
		if !strings.Contains(err.Error(), "key \"kind\" already set") {
			in.Fatalf("Decode() error = %q, want a duplicate field error", err)
		}
	})
}

func TestLoad(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		config, err := Load(testValidConfigLocation)
		if err != nil {
			in.Errorf("Load() error = %v", err)
		}
		if config == nil {
			in.Errorf("Load() returned a nil config")
		}
	})

	t.Run("no-file-exists", func(in *testing.T) {
		if _, err := Load("non existent"); !errors.Is(err, ErrLoad) {
			in.Errorf("expected: `%v`, got: `%v`", ErrLoad, err)
		}
	})

	t.Run("non-config-obj", func(in *testing.T) {
		if _, err := Load(testValidStatusLocation); !errors.Is(err, ErrLoad) {
			in.Errorf("expected: `%v`, got: `%v`", ErrLoad, err)
		}
	})
}
