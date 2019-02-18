/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package latency

import (
	"github.com/nalej/derrors"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"github.com/nalej/device-manager/internal/pkg/provider/latency"
	"github.com/nalej/grpc-device-controller-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
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

func (m * Manager) RegisterLatency (request *grpc_device_controller_go.RegisterLatencyRequest) (derrors.Error) {

	// Provider.AddLatency
	toAdd := entities.NewPingLatencyFromGRPC(request)
	err := m.pProvider.AddPingLatency(*toAdd)
	if err != nil {
		return err
	}

	return  nil
}

func (m * Manager) GetLatency(device *grpc_device_go.DeviceId) (*grpc_device_manager_go.LatencyMeasure, derrors.Error) {
	return nil, derrors.NewUnimplementedError("not implemented yet")
}
func (m * Manager) GetLatencyList(group *grpc_device_go.DeviceGroupId) (*grpc_device_manager_go.LatencyMeasureList, derrors.Error){
	return nil, derrors.NewUnimplementedError("not implemented yet")
}