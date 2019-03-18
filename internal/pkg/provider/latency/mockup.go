/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package latency

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"sync"
	"time"
)

type MockupProvider struct {
	sync.Mutex
	// latency indexed by organization_id, device_group_id, device_id, inserted
	latency map[string] entities.Latency
	// last inserted indexed by organization_id, device_group_id, device_id
	lastInserted map[string] int64
	// latency indexed by organization_id, device_group_id
	latencyGroup map[string][]*entities.Latency
}

func NewMockupProvider() * MockupProvider {
	return &MockupProvider{
		latency:make(map[string]entities.Latency, 0),
		lastInserted: make(map[string]int64, 0),
		latencyGroup: make(map[string][]*entities.Latency),
	}
}

func (m * MockupProvider) getKey(organizationID string, deviceGroupID string, deviceID string, inserted int64) string {

	key := fmt.Sprintf("%s-%s-%s-%v", organizationID, deviceGroupID, deviceID, inserted)
	return key
}

func (m * MockupProvider) getShortKey(organizationID string, deviceGroupID string, deviceID string) string {

	key := fmt.Sprintf("%s-%s-%s", organizationID, deviceGroupID, deviceID)
	return key
}

// GetLastPingLatency get the las latency measure of a device
func (m * MockupProvider) 	GetLastPingLatency (organizationID string, deviceGroupID string, deviceID string) (*entities.Latency, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	// get the last time inserted
	shortKey := m.getShortKey(organizationID, deviceGroupID, deviceID)
	inserted, exists := m.lastInserted[shortKey]
	if ! exists {
		return entities.NewEmptyLatency(), nil
	}
	key := m.getKey(organizationID, deviceGroupID, deviceID, inserted)
	latency, exists := m.latency[key]
	if ! exists {
		return entities.NewEmptyLatency(), nil
	}

	return &latency, nil
}

// AddPingLatency adds a new latency
func (m * MockupProvider) AddPingLatency(latency entities.Latency ) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := m.getKey(latency.OrganizationId, latency.DeviceGroupId, latency.DeviceId, latency.Inserted)
	m.latency[key] = latency

	shortKey := m.getShortKey(latency.OrganizationId, latency.DeviceGroupId, latency.DeviceId)
	m.lastInserted[shortKey] = latency.Inserted

	 list, exists := m.latencyGroup[fmt.Sprintf("%s-%s", latency.OrganizationId, latency.DeviceGroupId)]
	 if !exists {
	 	list = make ([]*entities.Latency, 0)

	 }
	list = append(list, &latency)
	m.latencyGroup[fmt.Sprintf("%s-%s", latency.OrganizationId, latency.DeviceGroupId)] = list

	return nil
}

func (m * MockupProvider) 	GetGroupLatency (organizationID string, deviceGroupID string) ([]*entities.Latency, derrors.Error){
	return m.GetGroupIntervalLatencies(organizationID, deviceGroupID, limitTime)
}

func (m * MockupProvider) GetGroupIntervalLatencies (organizationID string, deviceGroupID string, duration time.Duration) ([]*entities.Latency, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	threshold := time.Now().Add(-1 * duration).Unix()
	list, exists := m.latencyGroup[fmt.Sprintf("%s-%s", organizationID, deviceGroupID)]
	result := make ([]*entities.Latency, 0)
	if ! exists {
		return result, nil
	}

	for _, latency := range list{
		if latency.Inserted > threshold {
			result = append(result, latency)
		}
	}

	return result, nil
}

func (m * MockupProvider) GetLatency(organizationID string, deviceGroupID string, deviceID string) ([]*entities.Latency, derrors.Error){
	m.Lock()
	defer m.Unlock()

	latencies := make([]*entities.Latency, 0)
	list, exists := m.latencyGroup[fmt.Sprintf("%s-%s", organizationID, deviceGroupID)]
	if ! exists {
		return latencies, nil
	}
	for _, latency := range list {
		if latency.DeviceId == deviceID {
			latencies = append(latencies, latency)
		}
	}
	return latencies, nil
}

func (m *MockupProvider) RemoveLatency(organizationID string, deviceGroupID string, deviceID string) derrors.Error{
	m.Lock()
	defer m.Unlock()

	shortKey := m.getShortKey(organizationID, deviceGroupID, deviceID)
	inserted, exists := m.lastInserted[shortKey]
	if exists {
		key := m.getKey(organizationID, deviceGroupID, deviceID, inserted)
		delete(m.latency, key)
	}
	delete(m.lastInserted, shortKey)

	groupKey := fmt.Sprintf("%s-%s", organizationID, deviceGroupID)
	newLatencyList := make([]*entities.Latency, 0)
	for _, l := range m.latencyGroup[groupKey]{
		if l.DeviceId != deviceID {
			newLatencyList = append(newLatencyList, l)
		}
	}
	m.latencyGroup[groupKey] = newLatencyList

	return nil
}