package util

import (
	"os"
	"path/filepath"
)

// OSHelper is used mostly for OS related operations, but with possibilities to override
// system functions with testing assets
type OSHelper struct {
	ReadFile     func(string) ([]byte, error)
	Stat         func(name string) (os.FileInfo, error)
	Abs          func(string) (string, error)
	EvalSymlinks func(string) (string, error)
	Getwd        func() (string, error)
}

// NewOSHelper creates configured [OSHelper].
func NewOSHelper() *OSHelper {
	return &OSHelper{
		ReadFile:     os.ReadFile,
		Stat:         os.Stat,
		Abs:          filepath.Abs,
		EvalSymlinks: filepath.EvalSymlinks,
		Getwd:        os.Getwd,
	}
}

// WithReadFile sets custom ReadFile compatible function.
func (h *OSHelper) WithReadFile(fnc func(string) ([]byte, error)) *OSHelper {
	h.ReadFile = fnc
	return h
}

// WithStat sets compatible [os.Stat] function.
func (h *OSHelper) WithStat(fnc func(name string) (os.FileInfo, error)) *OSHelper {
	h.Stat = fnc
	return h
}

// WithABS sets compatible [filepath.Abs] function.
func (h *OSHelper) WithABS(fnc func(string) (string, error)) *OSHelper {
	h.Abs = fnc
	return h
}

// WithEvalSymlinks sets compatible [filepath.EvalSymlinks] function.
func (h *OSHelper) WithEvalSymlinks(fnc func(string) (string, error)) *OSHelper {
	h.EvalSymlinks = fnc
	return h
}

// WithGetwd sets compatible [os.Getwd] function.
func (h *OSHelper) WithGetwd(fnc func() (string, error)) *OSHelper {
	h.Getwd = fnc
	return h
}
