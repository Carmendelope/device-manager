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
	latency map[string] entities.Latency
}

func NewMockupProvider() * MockupProvider {
	return &MockupProvider{
		latency:make(map[string]entities.Latency, 0),
	}
}

func (m * MockupProvider) getKey(latency entities.Latency) string {

	key := fmt.Sprintf("%s-%s-%s-%v", latency.OrganizationId, latency.DeviceGroupId, latency.DeviceId, latency.Inserted)
	return key
}

func (m * MockupProvider) AddPingLatency(latency entities.Latency ) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := m.getKey(latency)
	m.latency[key] = latency

	return nil
}