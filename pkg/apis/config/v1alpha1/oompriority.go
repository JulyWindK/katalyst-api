/*
Copyright 2022 The Katalyst Authors.

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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=oompriorityconfigurations,shortName=oomp
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type=date,JSONPath=.metadata.creationTimestamp
// +kubebuilder:printcolumn:name="PAUSED",type=boolean,JSONPath=".spec.paused"
// +kubebuilder:printcolumn:name="SELECTOR",type=string,JSONPath=".spec.nodeLabelSelector"
// +kubebuilder:printcolumn:name="PRIORITY",type=string,JSONPath=".spec.priority"
// +kubebuilder:printcolumn:name="NODES",type=string,JSONPath=".spec.ephemeralSelector.nodeNames"
// +kubebuilder:printcolumn:name="DURATION",type=string,JSONPath=".spec.ephemeralSelector.lastDuration"
// +kubebuilder:printcolumn:name="TARGET",type=integer,JSONPath=".status.targetNodes"
// +kubebuilder:printcolumn:name="CANARY",type=integer,JSONPath=".status.canaryNodes"
// +kubebuilder:printcolumn:name="UPDATED-TARGET",type=integer,JSONPath=".status.updatedTargetNodes"
// +kubebuilder:printcolumn:name="HASH",type=string,JSONPath=".status.currentHash"
// +kubebuilder:printcolumn:name="VALID",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].status"
// +kubebuilder:printcolumn:name="REASON",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].reason"
// +kubebuilder:printcolumn:name="MESSAGE",type=string,JSONPath=".status.conditions[?(@.type==\"Valid\")].message"

// OOMPriorityConfiguration is the Schema for the configuration API used by OOM priority policy
type OOMPriorityConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OOMPriorityConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus          `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// OOMPriorityConfigurationList contains a list of OOMPriorityConfiguration
type OOMPriorityConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OOMPriorityConfiguration `json:"items"`
}

// OOMPriorityConfigurationSpec defines the desired state of OOMPriorityConfiguration
type OOMPriorityConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	// Config is custom field for oom priority configuration
	Config OOMPriorityQoSConfig `json:"config"`
}

type OOMPriorityQoSConfig struct {
	// EnableOOMPriority indicates whether to enable oom priority policy
	// +optional
	EnableOOMPriority *bool `json:"enableOOMPriority,omitempty"`
	// Interval is the interval to check oom priority
	// +optional
	Interval *int `json:"interval,omitempty"`
	// LowPriorityWartermarkRatioWithSwapOn is the watermark ratio to trigger low priority oom
	// +optional
	LowPriorityWartermarkRatioWithSwapOn *int `json:"lowPriorityWartermarkRatioWithSwapOn,omitempty"`
	// LowPriorityWartermarkRatioWithSwapOff is the watermark ratio to trigger low priority oom
	// +optional
	LowPriorityWartermarkRatioWithSwapOff *int `json:"lowPriorityWartermarkRatioWithSwapOff,omitempty"`
	// HighPrioKatalystQoS is the list of qos types to be treated as high priority
	// +optional
	HighPriorityKatalystQoS []string `json:"highPriorityKatalystQoS,omitempty"`
	// LowPriorityKatalystQoS is the list of qos types to be treated as low priority
	// +optional
	LowPriorityKatalystQoS []string `json:"lowPriorityKatalystQoS,omitempty"`
	// HighPriorityCgroupPaths is the list of cgroup paths to be treated as high priority
	// +optional
	HighPriorityCgroupPaths []string `json:"highPriorityCgroupPaths,omitempty"`
	// LowPriorityCgroupPaths is the list of cgroup paths to be treated as low priority
	// +optional
	LowPriorityCgroupPaths []string `json:"lowPriorityCgroupPaths,omitempty"`
}
