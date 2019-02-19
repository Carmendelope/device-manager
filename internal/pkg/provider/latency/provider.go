/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package latency

import (
	"github.com/nalej/derrors"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"time"
)

type Provider interface {

	// AddPingLatency adds a new latency
	AddPingLatency(entities.Latency ) derrors.Error

	// GetLastPingLatency get the las latency measure of a device
	GetLastPingLatency (organizationID string, deviceGroupID string, deviceID string) (*entities.Latency, derrors.Error)

	// GetGroupLatency get the last latency measures of a device_group (5 minutes of latencies)
	GetGroupLatency (organizationID string, deviceGroupID string) ([]*entities.Latency, derrors.Error)

	// GetGroupIntervalLatencies get the last latency measures of a device_group (time indicated in duration parameter of latencies)
	GetGroupIntervalLatencies(organizationID string, deviceGroupID string, duration time.Duration) ([]*entities.Latency, derrors.Error)

	GetLatency(organizationID string, deviceGroupID string, deviceID string) ([]*entities.Latency, derrors.Error)

}