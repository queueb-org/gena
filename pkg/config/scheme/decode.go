package scheme

import (
	"errors"
	"fmt"
	"io"
	"os"

	internal "queueb.org/gena/pkg/config"
)

var (
	// ErrLoad is returned if any error occurred during loading configuration file.
	ErrLoad = errors.New("could not load configuration file")
	// ErrDecode is returned if there was an error during object decoding operations.
	ErrDecode = errors.New("could not decode object")
)

// TODO: make decoder more flexible => enable strict on a user side as an option.

// Decode decodes one YAML configuration document, applies defaults for its
// declared API version, and converts it to the internal Config type.
//
// The document must contain a registered apiVersion and kind. Unknown and
// duplicate fields are rejected.
func Decode(r io.Reader) (*internal.Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%w: could not read object: %w", ErrDecode, err)
	}

	// UniversalDecoder reads versioned API and converts it to internal
	// (conversion and defaulter will be applied automatically).
	obj, gvk, err := Codecs.UniversalDecoder().Decode(data, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: could not read configuration: %w", ErrDecode, err)
	}

	config, ok := obj.(*internal.Config)
	if !ok {
		return nil, fmt.Errorf("%w: unexpected object %s, expected %T",
			ErrDecode, gvk, &internal.Config{},
		)
	}

	return config, nil
}

// Load reads and decodes one YAML configuration document from path.
func Load(path string) (config *internal.Config, err error) {
	var fd *os.File
	if fd, err = os.Open(path); err != nil {
		return nil, fmt.Errorf("%w: read configuration %w", ErrLoad, err)
	}
	defer fd.Close()

	if config, err = Decode(fd); err != nil {
		err = fmt.Errorf("%w: decode configuration failed: %w", ErrLoad, err)
	}

	return
}
