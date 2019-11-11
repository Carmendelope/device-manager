/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package entities

import (
	"github.com/nalej/grpc-device-controller-go"
	"time"
)

type Latency struct {
	// organization identifier
	OrganizationId string `json:"organization_id,omitempty"`
	// device_group identifier
	DeviceGroupId string `json:"device_group_id,omitempty"`
	// device identifier
	DeviceId string `json:"device_id,omitempty"`
	// Latency to register
	Latency int `json:"latency,omitempty"`
	// timestamp
	Inserted int64 `json:"inserted, omitempty"`
}

func NewEmptyLatency() *Latency {
	return &Latency{
		OrganizationId: "",
		DeviceGroupId:  "",
		DeviceId:       "",
		Latency:        -1,
		Inserted:       0,
	}
}

func NewPingLatencyFromGRPC(request *grpc_device_controller_go.RegisterLatencyRequest) *Latency {
	return &Latency{
		OrganizationId: request.OrganizationId,
		DeviceGroupId:  request.DeviceGroupId,
		DeviceId:       request.DeviceId,
		Latency:        int(request.Latency),
		Inserted:       time.Now().Unix(),
	}
}
