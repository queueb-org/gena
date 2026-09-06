package common

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalid = errors.New("invalid interface")
)

var _ io.WriteCloser = &nopCloser{w: nil}

type nopCloser struct {
	w io.Writer
}

// NopCloser creates No Operation closer proxy (adapter)
// for [io.Writer] object.
func NopCloser(w io.Writer) (io.WriteCloser, error) {
	if _, ok := w.(io.Closer); ok {
		return nil, fmt.Errorf("%w: writer implements io.Closer interface as well, "+
			"it's not allowed to use nopCloser", ErrInvalid)
	}
	return &nopCloser{w: w}, nil
}

// MustNopCloser returns [io.WriteCloser] interface base on input
// writer, otherwise panics.
func MustNopCloser(w io.Writer) io.WriteCloser {
	newW, err := NopCloser(w)
	if err != nil {
		panic(err)
	}

	return newW
}

// Write provides [io.Writer] interface.
func (w *nopCloser) Write(buf []byte) (n int, err error) {
	return w.w.Write(buf)
}

// Close provides [io.Closer] which does nothing.
func (w *nopCloser) Close() (err error) {
	return nil
}

// Close tries to Close all objects passed through objs.
func Close(objs ...any) {
	for _, o := range objs {
		if v, ok := o.(io.Closer); ok {
			_ = v.Close()
		}
	}
}

// Or returns first met non-blank element of slice.
func Or[E comparable](items ...E) E {
	var blank E
	for _, item := range items {
		if blank == item {
			continue
		}
		return item
	}

	return blank
}
