/*
Copyright 2019 Thomas Liang.

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
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ShardingRedisSpec defines the desired state of ShardingRedis
type ShardingRedisSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Replicas  int32                   `json:"replicas,omitempty"`
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	Image     string                  `json:"image,omitempty"`
}

// ShardingRedisStatus defines the observed state of ShardingRedis
type ShardingRedisStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase      Phase       `json:"phase"`
	Conditions []Condition `json:"conditions"`
}

type Phase string

type Condition struct {
	Type           ConditionType `json:"type"`
	Reason         string        `json:"reason"`
	TransitionTime string        `json:"transitionTime"`
}

type ConditionType string

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ShardingRedis is the Schema for the shardingredis API
// +k8s:openapi-gen=true
type ShardingRedis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ShardingRedisSpec   `json:"spec,omitempty"`
	Status ShardingRedisStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ShardingRedisList contains a list of ShardingRedis
type ShardingRedisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ShardingRedis `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ShardingRedis{}, &ShardingRedisList{})
}
