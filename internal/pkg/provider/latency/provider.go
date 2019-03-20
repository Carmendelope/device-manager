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

	// ------------- //
	// -- Latency -- //
	// ------------- //
	// AddPingLatency adds a new latency
	AddPingLatency(entities.Latency ) derrors.Error

	// TODO: delete this !
	// GetGroupIntervalLatencies get the last latency measures of a device_group (time indicated in duration parameter of latencies)
	GetGroupIntervalLatencies(organizationID string, deviceGroupID string, duration time.Duration) ([]*entities.Latency, derrors.Error)

	GetLatency(organizationID string, deviceGroupID string, deviceID string) ([]*entities.Latency, derrors.Error)

	// RemoveLatency removes the entries associated with a given device.
	RemoveLatency(organizationID string, deviceGroupID string, deviceID string) derrors.Error

	// ------------------ //
	// -- Last Latency -- //
	// ------------------ //
	// AddLastLatency
	AddLastLatency (latency entities.Latency) derrors.Error

	// GetLastPingLatency get the las latency measure of a device
	GetLastLatency (organizationID string, deviceGroupID string, deviceID string) (*entities.Latency, derrors.Error)

	// GetGroupLastLatencies get all the last latencies of the devices in the group
	GetGroupLastLatencies(organizationID string, deviceGroupID string)([]*entities.Latency, derrors.Error)
}