/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package latency

import (
	"github.com/nalej/device-manager/internal/pkg/entities"
	"github.com/nalej/device-manager/internal/pkg/provider/latency"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-controller-go"
)

type Manager struct {
	pProvider latency.Provider

}

// NewManager creates a Manager using a set of clients.
func NewManager(provider latency.Provider) Manager {
	return Manager{
		pProvider: provider,
	}
}

func (m * Manager) RegisterLatency (request *grpc_device_controller_go.RegisterLatencyRequest) (*grpc_common_go.Success, error) {

	err := entities.ValidRegisterLatencyRequest(request)
	if err != nil {
		return nil, err
	}
	// Provider.AddLatency
	toAdd := entities.NewPingLatencyFromGRPC(request)
	err = m.pProvider.AddPingLatency(*toAdd)
	if err != nil {
		return nil, err
	}

	return &grpc_common_go.Success{}, nil
}