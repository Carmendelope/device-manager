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

package latency

import (
	"github.com/nalej/derrors"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"github.com/nalej/device-manager/internal/pkg/provider/latency"
	"github.com/nalej/grpc-device-controller-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/rs/zerolog/log"
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

func (m *Manager) RegisterLatency(request *grpc_device_controller_go.RegisterLatencyRequest) derrors.Error {

	// AddLatency
	toAdd := entities.NewPingLatencyFromGRPC(request)
	err := m.pProvider.AddPingLatency(*toAdd)
	if err != nil {
		return err
	}

	// Update lastLatency info
	err = m.pProvider.AddLastLatency(*toAdd)
	if err != nil {
		log.Warn().Interface("latency", toAdd).Msg("unable to update last latency")
		return err
	}

	return nil
}

func (m *Manager) GetLatency(device *grpc_device_go.DeviceId) (*grpc_device_manager_go.LatencyMeasure, derrors.Error) {
	return nil, derrors.NewUnimplementedError("not implemented yet")
}
func (m *Manager) GetLatencyList(group *grpc_device_go.DeviceGroupId) (*grpc_device_manager_go.LatencyMeasureList, derrors.Error) {
	return nil, derrors.NewUnimplementedError("not implemented yet")
}
