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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=irqtuningconfigurations,shortName=irtc
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

// IRQTuningConfiguration is the Schema for the configuration API used by IRQ Tuning
type IRQTuningConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IRQTuningConfigurationSpec `json:"spec,omitempty"`
	Status GenericConfigStatus        `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// IRQTuningConfigurationList contains a list of IRQTuningConfiguration
type IRQTuningConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IRQTuningConfiguration `json:"items"`
}

type IRQTuningConfigurationSpec struct {
	GenericConfigSpec `json:",inline"`

	// Config is custom field for irq tuning configuration
	Config IRQTuningConfig `json:"config"`
}

type IRQTuningConfig struct {
	// +kubebuilder:default=false
	EnableIRQTuner *bool `json:"enableIRQTuner"`
	// +kubebuilder:default=false
	EnableIRQCoresExclusion *bool `json:"enableIRQCoresExclusion"`
	// irq tuning periodic interval
	// +kubebuilder:default=5
	Interval *int `json:"interval"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=50
	IRQCoresExpectedCPUUtil  *int                          `json:"irqCoresExpectedCPUUtil"`
	IRQCoreNetOverLoadThresh *IRQCoreNetOverloadThresholds `json:"irqCoreNetOverLoadThresh"`
	IRQLoadBalanceConf       *IRQLoadBalanceConfig         `json:"irqLoadBalanceConf"`
	IRQCoresAdjustConf       *IRQCoresAdjustConfig         `json:"irqCoresAdjustConf"`
	IRQCoresExclusionConf    *IRQCoresExclusionConfig      `json:"irqCoresExclusionConf"`
}

type IRQCoreNetOverloadThresholds struct {
	// ratio of softnet_stat 3rd col time_squeeze packets / softnet_stat 1st col processed packets
	// +kubebuilder:default=0.1
	IRQCoreSoftNetTimeSqueezeRatio *float64 `json:"irqCoreSoftNetTimeSqueezeRatio"`
}
type IRQLoadBalanceConfig struct {
	// interval of two successive irq load balance MUST greater-equal this interval
	// +kubebuilder:default=10
	// +kubebuilder:validation:Minimum=1
	SuccessiveTuningInterval *int                            `json:"successiveTuningInterval"`
	Thresholds               *IRQLoadBalanceTuningThresholds `json:"thresholds"`
	// two successive tunes whose interval is less-equal this threshold will be considered as pingpong tunings
	// +kubebuilder:default=180
	PingPongIntervalThresh *int `json:"pingPongIntervalThresh"`
	// ping pong count greater-equal this threshold will trigger increasing irq cores
	// +kubebuilder:default=1
	PingPongCountThresh *int `json:"pingPongCountThresh"`
	// max number of irqs are permitted to be tuned from some irq cores to other cores in each time, allowed value {1, 2}
	// +kubebuilder:validation:Enum={1,2}
	// +kubebuilder:default=2
	IRQsTunedNumMaxEachTime *int `json:"irqsTunedNumMaxEachTime"`
	// max number of irq cores whose affinitied irqs are permitted to tuned to other cores in each time, allowed value {1,2}
	// +kubebuilder:validation:Enum={1,2}
	// +kubebuilder:default=1
	IRQCoresTunedNumMaxEachTime *int       `json:"irqCoresTunedNumMaxEachTime"`
	RPSConf                     *RPSConfig `json:"rpsConf"`
}
type IRQLoadBalanceTuningThresholds struct {
	// irq core cpu util threshold, which will trigger irq cores load balance, generally this value should greater-equal IrqCoresExpectedCpuUtil
	// +kubebuilder:validation:default=65
	IRQCoreCPUUtilThresh *int `json:"irqCoreCPUUtilThresh"`
	// threshold of cpu util gap between source core and dest core of irq affinity changing
	// +kubebuilder:validation:default=20
	IRQCoreCPUUtilGapThresh *int `json:"irqCoreCPUUtilGapThresh"`
}
type RPSConfig struct {
	// if enable RPS when some single irq's cpu demand exceeds one whole CPU
	// +kubebuilder:default=true
	EnableRPS *bool `json:"enableRPS"`
	// if reset RPS dest cpus when irq cores adjusted, be careful, reset dest cpus may result in out-of-order packets, so limit the reset rate.
	// +kubebuilder:default=true
	ResetDestCPUsWhenIrqCoresAdjusted *bool `json:"resetDestCPUsWhenIrqCoresAdjusted"`

	EnterThresh *RPSEnterThresholds `json:"enterThresh"`

	ExitThresh *RPSExitThresholds `json:"exitThresh"`
}
type RPSEnterThresholds struct {
	// generally this value greater-than IrqLoadBalanceTuningThresholds.IrqCoreCpuUtilThresh
	// +kubebuilder:default=80
	SingleIrqCPUUtil *int `json:"singleIrqCPUUtil"`
}
type RPSExitThresholds struct {
	// Eligible PPS means (rps enabled queue's pps) - (non-rps queues's max PPS)
	// +kubebuilder:default=120
	EligiblePPSDuration *int `json:"eligiblePPSDuration"`
}
type IRQCoresAdjustConfig struct {
	// minimum percent of (100 * irq cores/total(or socket) cores), valid value [0,100], default 2
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=2
	IRQCoresPercentMin *int `json:"irqCoresPercentMin"`

	// maximum percent of (100 * irq cores/total(or socket) cores), valid value [0,100], default 30
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=30
	IRQCoresPercentMax *int `json:"irqCoresPercentMax"`

	IRQCoresIncConf *IRQCoresIncConfig `json:"irqCoresIncConf"`
	IRQCoresDecConf *IRQCoresDecConfig `json:"irqCoresDecConf"`
}

type IRQCoresIncConfig struct {
	// interval of two successive irq cores increase MUST greater-equal this interval
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=5
	SuccessiveIncInterval *int `json:"successiveIncInterval"`
	// scale factor of irq cores cpu usage to calculate expected irq cores number when irq cores cpu util nearly full, default 1.5
	// +kubebuilder:default=1.5
	IRQCoresFullScaleFactor *float64 `json:"irqCoresFullScaleFactor"`

	Thresholds *IRQCoresIncThresholds `json:"thresholds"`
}

type IRQCoresIncThresholds struct {
	// threshold of increasing irq cores, generally this thresh equal to or a litter greater-than IrqCoresExpectedCpuUtil
	// +kubebuilder:default=60
	IRQCoresAvgCPUUtilThresh *int `json:"irqCoresAvgCPUUtilThresh"`
}

type IRQCoresDecConfig struct {
	// interval of two successive irq cores decrease MUST greater-equal this interval
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=120
	SuccessiveDecInterval *int                   `json:"successiveDecInterval"`
	Thresholds            *IRQCoresDecThresholds `json:"thresholds"`
	// max cores to decrease each time, deault 1
	// +kubebuilder:default=1
	DecCoresMaxEachTime *int `json:"decCoresMaxEachTime"`
}

type IRQCoresDecThresholds struct {
	// threshold of decreasing irq cores, generally this thresh should be less-than IrqCoresExpectedCpuUtil
	// +kubebuilder:default=35
	IRQCoresAvgCPUUtilThresh *int `json:"irqCoresAvgCPUUtilThresh"`
}

type IRQCoresExclusionConfig struct {
	Thresholds *IRQCoresExclusionThresholds `json:"thresholds"`
	// interval of successive enable/disable irq cores exclusion MUST >= SuccessiveSwitchInterval
	// +kubebuilder:default=600
	SuccessiveSwitchInterval *float64 `json:"successiveSwitchInterval"`
}

type IRQCoresExclusionThresholds struct {
	EnableThresholds *EnableIRQCoresExclusionThresholds `json:"enableThresholds"`

	DisableThresholds *DisableIRQCoresExclusionThresholds `json:"disableThresholds"`
}

type EnableIRQCoresExclusionThresholds struct {
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Format="uint64"
	// +kubebuilder:default=60000
	RxPPSThresh *uint64 `json:"rxPPSThresh"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=30
	SuccessiveCount *int `json:"successiveCount"`
}

type DisableIRQCoresExclusionThresholds struct {
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Format="uint64"
	// +kubebuilder:default=30000
	RxPPSThresh *uint64 `json:"rxPPSThresh"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=30
	SuccessiveCount *int `json:"successiveCount"`
}
