
/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
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
	Latency  int `json:"latency,omitempty"`
	// timestamp
	Inserted int64   `json:"inserted, omitempty"`
}

func NewPingLatencyFromGRPC(request * grpc_device_controller_go.RegisterLatencyRequest) *Latency{
	return &Latency{
		OrganizationId: request.OrganizationId,
		DeviceGroupId: 	request.DeviceGroupId,
		DeviceId: 		request.DeviceId,
		Latency: 		int(request.Latency),
		Inserted: 		time.Now().Unix(),
	}
}
