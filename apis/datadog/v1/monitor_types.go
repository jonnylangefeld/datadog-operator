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

package v1

import (
	"encoding/json"

	datadog "github.com/zorkian/go-datadog-api"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MonitorSpec defines the desired state of Monitor
type MonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Name of the monitor
	Name string `json:"name"`
	// Type of the monitor
	Type string `json:"type"`
	// Message of the monitor
	Message string `json:"message"`
	// Query of the monitor
	Query string `json:"query"`
	// Tags of the monitor
	Tags []string `json:"tags,omitempty"`
	// Options of the monitor
	// +kubebuilder:pruning:PreserveUnknownFields
	Options *runtime.RawExtension `json:"options,omitempty"`
}

// MonitorStatus defines the observed state of Monitor
type MonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// ID of the monitor
	ID int `json:"id,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Monitor is the Schema for the monitors API
type Monitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MonitorSpec   `json:"spec,omitempty"`
	Status MonitorStatus `json:"status,omitempty"`
}

// Transform transforms a kubernetes custom resource object into a datadog api object
func (in *Monitor) Transform() (*datadog.Monitor, error) {
	options := &datadog.Options{}
	if in.Spec.Options != nil {
		err := json.Unmarshal(in.Spec.Options.Raw, options)
		if err != nil {
			return nil, err
		}
	}

	return &datadog.Monitor{
		Id:      &in.Status.ID,
		Type:    &in.Spec.Type,
		Query:   &in.Spec.Query,
		Name:    &in.Spec.Name,
		Message: &in.Spec.Message,
		Tags:    in.Spec.Tags,
		Options: options,
	}, nil
}

// +kubebuilder:object:root=true

// MonitorList contains a list of Monitor
type MonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Monitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Monitor{}, &MonitorList{})
}
