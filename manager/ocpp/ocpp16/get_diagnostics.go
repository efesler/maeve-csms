// SPDX-License-Identifier: Apache-2.0

package ocpp16

type GetDiagnosticsJson struct {
	// Location corresponds to the JSON schema field "location".
	Location string `json:"location" yaml:"location" mapstructure:"location"`

	// Retries corresponds to the JSON schema field "retries".
	Retries int `json:"retries,omitempty" yaml:"retries,omitempty" mapstructure:"retries,omitempty"`

	// RetryInterval corresponds to the JSON schema field "retryInterval".
	RetryInterval int `json:"retryInterval,omitempty" yaml:"retryInterval,omitempty" mapstructure:"retryInterval,omitempty"`

	// StartTime corresponds to the JSON schema field "startTime".
	StartTime *string `json:"startTime,omitempty" yaml:"startTime,omitempty" mapstructure:"startTime,omitempty"`

	// StopTime corresponds to the JSON schema field "stopTime".
	StopTime *string `json:"stopTime,omitempty" yaml:"stopTime,omitempty" mapstructure:"stopTime,omitempty"`
}

func (*GetDiagnosticsJson) IsRequest() {}
