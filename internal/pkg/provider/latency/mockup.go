/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
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
	list, exists := m.latencyGroup[fmt.Sprintf("%s-%s", organizationID, deviceGroupID)]
	if ! exists {
		return make([]*entities.Latency, 0), nil
	}
	return list, nil
}
