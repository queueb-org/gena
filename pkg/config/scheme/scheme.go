// Package scheme provides the runtime scheme and codecs for all supported
// versions of the gena configuration API.
package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	internal "queueb.org/gena/pkg/config"
	"queueb.org/gena/pkg/config/v1alpha1"
)

var (
	// Scheme contains the internal configuration API and every supported wire
	// version.
	Scheme = runtime.NewScheme()

	// Codecs performs strict decoding, defaulting, and conversion between the
	// registered configuration API versions.
	Codecs = serializer.NewCodecFactory(Scheme, serializer.EnableStrict)
)

func init() {
	utilruntime.Must(internal.Install(Scheme))
	utilruntime.Must(v1alpha1.Install(Scheme))

	utilruntime.Must(
		Scheme.SetVersionPriority(schema.GroupVersion(v1alpha1.GroupVersion)),
	)
}
