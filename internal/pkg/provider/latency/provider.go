/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package latency

import (
	"github.com/nalej/derrors"
	"github.com/nalej/device-manager/internal/pkg/entities"
)

type Provider interface {

	// AddPingLatency adds a new latency
	AddPingLatency(entities.Latency ) derrors.Error

}