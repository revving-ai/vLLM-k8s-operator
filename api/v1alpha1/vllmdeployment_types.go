/*
Copyright 2024.

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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VllmDeploymentSpec defines the desired state of VllmDeployment.
type VllmDeploymentSpec struct {
	Replicas    int32           `json:"replicas"`
	Model       *ModelConfig    `json:"model"`
	VLLMConfig  *VLLMConfig     `json:"vLLMConfig"`
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	Container   v1.Container    `json:"containers,omitempty"`
	// TODO (similar to prometheus): VolumeClaimTemplate EmbeddedPersistentVolumeClaim `json:"volumeClaimTemplate,omitempty"`
}

type ModelConfig struct {
	Name  string `json:"name"`
	HfURL string `json:"hf_url"`
}

type VLLMConfig struct {
	Port                 int    `json:"port"`
	GpuMemoryUtilization string `json:"gpu-memory-utilization"`
	LogLevel             string `json:"log-level"`
	BlockSize            int    `json:"block-size"`
	MaxModelLen          int    `json:"max-model-len"`
}

// VllmDeploymentStatus defines the observed state of VllmDeployment.
type VllmDeploymentStatus struct {
	// The current state of the Prometheus deployment.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []Condition `json:"conditions,omitempty"`
}

type Condition struct {
	// Type of the condition being reported.
	// +required
	Type ConditionType `json:"type"`
}

// +kubebuilder:validation:MinLength=1
type ConditionType string

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// VllmDeployment is the Schema for the vllmdeployments API.
type VllmDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VllmDeploymentSpec   `json:"spec,omitempty"`
	Status VllmDeploymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VllmDeploymentList contains a list of VllmDeployment.
type VllmDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VllmDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VllmDeployment{}, &VllmDeploymentList{})
}
