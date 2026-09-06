package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Config represents code-Generators Assistant (Gen-A) main configuration file.
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// // +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// // ConfigList represents [Config] objects over API.
// type ConfigList struct {
// 	metav1.TypeMeta `json:",inline"`
// 	metav1.ListMeta `json:"metadata,omitempty"`
// 	Items           []Config `json:"items,omitempty"`
// }
