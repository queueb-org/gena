package config

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Foo struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Stub string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList represents [Foo] objects over API.
type FooList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []Foo
}
