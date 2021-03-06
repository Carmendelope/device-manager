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
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"sync"
)

type MockupProvider struct {
	sync.Mutex
	// latency indexed by organization_id, device_group_id, device_id, inserted
	latency map[string][]*entities.Latency
	// lastLatency indexed by organization, device_group_id + device_id
	lastLatency map[string]map[string]*entities.Latency
}

func NewMockupProvider() *MockupProvider {
	return &MockupProvider{
		latency:     make(map[string][]*entities.Latency, 0),
		lastLatency: make(map[string]map[string]*entities.Latency, 0),
	}
}

func (m *MockupProvider) getKey(organizationID string, deviceGroupID string, deviceID string) string {

	key := fmt.Sprintf("%s-%s-%s", organizationID, deviceGroupID, deviceID)
	return key
}

func (m *MockupProvider) getShortKey(organizationID string, deviceGroupID string) string {

	key := fmt.Sprintf("%s-%s", organizationID, deviceGroupID)
	return key
}

// AddPingLatency adds a new latency
func (m *MockupProvider) AddPingLatency(latency entities.Latency) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := m.getKey(latency.OrganizationId, latency.DeviceGroupId, latency.DeviceId)

	_, exists := m.latency[key]
	if !exists {
		m.latency[key] = make([]*entities.Latency, 0)
	}
	m.latency[key] = append(m.latency[key], &latency)

	return nil
}

func (m *MockupProvider) GetLatency(organizationID string, deviceGroupID string, deviceID string) ([]*entities.Latency, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	list, exists := m.latency[m.getKey(organizationID, deviceGroupID, deviceID)]
	if !exists {
		latencies := make([]*entities.Latency, 0)
		return latencies, nil
	}
	return list, nil
}

func (m *MockupProvider) RemoveLatency(organizationID string, deviceGroupID string, deviceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()

	delete(m.latency, m.getKey(organizationID, deviceGroupID, deviceID))

	return nil
}

// GetLastPingLatency get the las latency measure of a device
func (m *MockupProvider) GetLastLatency(organizationID string, deviceGroupID string, deviceID string) (*entities.Latency, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	key := m.getShortKey(organizationID, deviceGroupID)

	latencies, exists := m.lastLatency[key]
	if !exists {
		return entities.NewEmptyLatency(), nil
	}

	latency, exists := latencies[deviceID]
	if !exists {
		return entities.NewEmptyLatency(), nil
	}
	return latency, nil
}

func (m *MockupProvider) AddLastLatency(latency entities.Latency) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := m.getShortKey(latency.OrganizationId, latency.DeviceGroupId)

	latencies, exists := m.lastLatency[key]
	if !exists {
		latencies = make(map[string]*entities.Latency, 0)
		m.lastLatency[key] = latencies
	}
	latencies[latency.DeviceId] = &latency

	return nil
}

func (m *MockupProvider) GetGroupLastLatencies(organizationID string, deviceGroupID string) ([]*entities.Latency, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	latencies := make([]*entities.Latency, 0)
	key := m.getShortKey(organizationID, deviceGroupID)

	list, exists := m.lastLatency[key]
	if !exists {
		return latencies, nil

	}
	for _, latency := range list {
		latencies = append(latencies, latency)
	}

	return latencies, nil
}
