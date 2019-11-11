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
)

type Provider interface {

	// ------------- //
	// -- Latency -- //
	// ------------- //
	// AddPingLatency adds a new latency
	AddPingLatency(entities.Latency) derrors.Error

	GetLatency(organizationID string, deviceGroupID string, deviceID string) ([]*entities.Latency, derrors.Error)

	// RemoveLatency removes the entries associated with a given device.
	RemoveLatency(organizationID string, deviceGroupID string, deviceID string) derrors.Error

	// ------------------ //
	// -- Last Latency -- //
	// ------------------ //
	// AddLastLatency
	AddLastLatency(latency entities.Latency) derrors.Error

	// GetLastPingLatency get the las latency measure of a device
	GetLastLatency(organizationID string, deviceGroupID string, deviceID string) (*entities.Latency, derrors.Error)

	// GetGroupLastLatencies get all the last latencies of the devices in the group
	GetGroupLastLatencies(organizationID string, deviceGroupID string) ([]*entities.Latency, derrors.Error)
}
