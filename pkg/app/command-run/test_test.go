package commandrun

import "testing"

// TestOption contains different options for testing [Command].
type TestOption struct {
	Dir  string
	Args []string
}

func mergeTestOptions(t testing.TB, opts ...*TestOption) *TestOption {
	def := &TestOption{
		Dir:  "",
		Args: []string{},
	}

	for _, o := range opts {
		if o.Dir != "" {
			def.Dir = o.Dir
		}

		// note, there's no uniqueness check ;)
		if o.Args != nil {
			def.Args = append(def.Args, o.Args...)
		}
	}

	return def
}

// TOptDir sets discover directory.
func TOptDir(dir string) *TestOption {
	return &TestOption{Dir: dir}
}

// TOptDeepcopyGen sets deepcopy-gen as an argument.
func TOptDeepcopyGen() *TestOption {
	return &TestOption{Args: []string{"deepcopy-gen"}}
}

// TOptDefaulterGen sets defaulter-gen as an argument.
func TOptDefaulterGen() *TestOption {
	return &TestOption{Args: []string{"defaulter-gen"}}
}
