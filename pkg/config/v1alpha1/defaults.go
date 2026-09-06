package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

func init() {
	localSchemeBuilder.Register(addDefaultingFuncs)
}

// addDefaultingFuncs registers generated defaulting functions with the scheme.
func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_Config sets defaults for [Config] object.
func SetDefaults_Config(in *Config) {
	//lint:ignore S1023, intended.
	return
}
