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
	"context"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-controller-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
)

// Handler structure for the node requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h * Handler) RegisterLatency (ctx context.Context, request *grpc_device_controller_go.RegisterLatencyRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidRegisterLatencyRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	err =  h.Manager.RegisterLatency(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}
// TODO: change getLatency to GetDeviceLatencies
func (h * Handler) GetLatency(ctx context.Context, device *grpc_device_go.DeviceId) (*grpc_device_manager_go.LatencyMeasure, error) {
	err := entities.ValidDeviceID(device)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	list, err :=  h.Manager.GetLatency(device)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return list, nil
}
// TODO: change getLatency to GetDeviceGroupLatencies
func (h * Handler) GetLatencyList(ctx context.Context, group *grpc_device_go.DeviceGroupId) (*grpc_device_manager_go.LatencyMeasureList, error){
	err := entities.ValidDeviceGroupID(group)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	list, err :=  h.Manager.GetLatencyList(group)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return list, nil
}