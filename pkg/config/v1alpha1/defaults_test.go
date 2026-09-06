package v1alpha1

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

// WithNewLocalScheme creates local [runtime.Scheme] and installs
// there current package types.
func WithNewLocalScheme(t testing.TB) *runtime.Scheme {
	scheme := runtime.NewScheme()
	if err := Install(scheme); err != nil {
		t.Fatalf("could not install scheme: %v", err)
	}
	return scheme
}

func TestSetDefaults_Config(t *testing.T) {
	scheme := WithNewLocalScheme(t)

	for _, entry := range []struct {
		name     string
		in       *Config
		expected *Config
	}{
		{
			name:     "defaults",
			in:       &Config{},
			expected: &Config{},
		},
	} {
		t.Run(entry.name, func(in *testing.T) {
			// runs default operation.
			scheme.Default(entry.in)

			if !reflect.DeepEqual(entry.in, entry.expected) {
				in.Errorf("\nexp:\n%#v\ngot:\n%#v\n", entry.expected, entry.in)
			}
		})
	}
}
