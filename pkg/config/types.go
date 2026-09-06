package config

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Config struct {
	metav1.TypeMeta
	metav1.ObjectMeta
}

// // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// // ConfigList represents [Config] objects over API.
// type ConfigList struct {
// 	metav1.TypeMeta
// 	metav1.ListMeta
// 	Items []Config
// }
