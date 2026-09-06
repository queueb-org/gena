package config

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
)

func TestAddKnownTypes(t *testing.T) {
	if err := addKnownTypes(runtime.NewScheme()); err != nil {
		t.Errorf("got error: %v", err)
	}
}
