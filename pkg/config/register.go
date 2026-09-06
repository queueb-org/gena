package config

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"queueb.org/apimachinery/pkg/api"
)

// GroupName is a name for APIs united under single logical group.
const GroupName = "gena.apps.queueb.org"

// GroupVersion specifies the group and the version used to register the objects.
var GroupVersion = v1.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

// SchemeGroupVersion is group version used to register these objects
//
// Deprecated: Use GroupVersion instead.
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

// Exported variables
var (
	// SchemeBuilder is a scheme builder asset for the package
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &SchemeBuilder

	Resource = api.Resource(SchemeGroupVersion)
	Kind     = api.Kind(SchemeGroupVersion)

	Install = localSchemeBuilder.AddToScheme
)

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Config{},
	)
	return nil
}
