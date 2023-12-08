/*
Copyright 2023.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NfsShareSpec defines the desired state of NfsShare
type NfsShareSpec struct {
	// +kubebuilder:validation:Required
	Kyma string `json:"kyma"`

	// +kubebuilder:validation:Required
	Provider string `json:"provider"`

	// +kubebuilder:validation:Required
	Capacity string `json:"capacity"`
}

// NfsShareStatus defines the observed state of NfsShare
type NfsShareStatus struct {
	State    StatusState `json:"state,omitempty"`
	Capacity string      `json:"capacity,omitempty"`

	// +optional
	Scope *Scope `json:"scope,omitempty"`

	// List of status conditions to indicate the status of a Peering.
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NfsShare is the Schema for the nfsshares API
type NfsShare struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NfsShareSpec   `json:"spec,omitempty"`
	Status NfsShareStatus `json:"status,omitempty"`
}

func (in *NfsShare) Kyma() string {
	return in.Spec.Kyma
}

//+kubebuilder:object:root=true

// NfsShareList contains a list of NfsShare
type NfsShareList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NfsShare `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NfsShare{}, &NfsShareList{})
}
