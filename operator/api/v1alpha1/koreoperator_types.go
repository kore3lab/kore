/*
Copyright 2022.

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

func init() {
	SchemeBuilder.Register(&KoreOperator{}, &KoreOperatorList{})
}

// type InstallStatus int32
type InstallStatus string

const (
	STATUS_NONE        InstallStatus = "NONE"
	STATUS_RECONCILING InstallStatus = "RECONCILING"
	STATUS_ERROR       InstallStatus = "ERROR"
	STATUS_UPDATING    InstallStatus = "UPDATING"
	STATUS_DELETING    InstallStatus = "DELETING"
	STATUS_DELETED     InstallStatus = "DELETED"
	STATUS_HEALTHY     InstallStatus = "HEALTHY"
	STATUS_COMPLETE    InstallStatus = "COMPLETE"
)

// KoreOperator is the Schema for the koreoperators API
type KoreOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              *KoreOperatorSpec   `json:"spec,omitempty"`
	Status            *KoreOperatorStatus `json:"status,omitempty"`
}

// KoreOperatorSpec defines the desired state of KoreOperator
type KoreOperatorSpec struct {
	Revision   string                           `json:"revision,omitempty"`
	Components map[string]KoreOperatorComponent `json:"components,omitempty"`
	Values     map[string]interface{}           `json:"values,omitempty"`
}

type KoreOperatorComponent struct {
	Enabled bool `json:"enabled,omitempty"`
}

type KoreOperatorStatus struct {
	Status  InstallStatus `json:"status,omitempty"`
	Message string        `json:"message,omitempty"`
}

type KoreOperatorComponentStatus struct {
	Version string        `json:"version,omitempty"`
	Status  InstallStatus `json:"status,omitempty"`
	Error   string        `json:"error,omitempty"`
}

//+kubebuilder:object:root=true

// KoreOperatorList contains a list of KoreOperator
type KoreOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KoreOperator `json:"items"`
}
