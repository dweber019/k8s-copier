/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CopyResourceSpec defines the desired state of CopyResource
type CopyResourceSpec struct {
	// The Kind of the Resource you like to copy
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Secret;ConfigMap
	Kind string `json:"kind"`

	// The MetaName of the Resource found in metadata.name
	// +kubebuilder:validation:Required
	MetaName string `json:"metaName"`

	// The TargetNamespace the Resource should be copied to
	// +kubebuilder:validation:Required
	TargetNamespace string `json:"targetNamespace"`
}

// CopyResourceStatus defines the observed state of CopyResource
type CopyResourceStatus struct {
	ResourceVersion string `json:"resourceVersion"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CopyResource is the Schema for the copyresources API
type CopyResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CopyResourceSpec   `json:"spec,omitempty"`
	Status CopyResourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CopyResourceList contains a list of CopyResource
type CopyResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CopyResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CopyResource{}, &CopyResourceList{})
}
