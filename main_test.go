package main

import (
	"testing"

	"common.queueb.org/tests"
)

// helpers

func WithExit(t testing.TB, exitCodes *[]int) {
	orig := exit

	t.Cleanup(func() {
		exit = orig
	})

	exit = func(code int) {
		*exitCodes = append(*exitCodes, code)
	}
}

// tests

func TestMain(t *testing.T) {
	t.Run("ok", func(in *testing.T) {
		var codes []int
		WithExit(t, &codes)
		main()

		if len(codes) != 0 {
			in.Errorf("abnormal exit codes expected to be blank")
		}
	})

	t.Run("on-err", func(in *testing.T) {
		var codes []int
		WithExit(t, &codes)

		tests.WithArgs(in, []string{"gena", "__test"})
		main()

		if len(codes) == 0 {
			in.Errorf("exit codes expected not to be blank")
		}
	})
}
