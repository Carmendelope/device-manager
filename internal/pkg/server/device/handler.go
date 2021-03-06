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

package device

import (
	"context"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

// Handler structure for the node requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) AddDeviceGroup(ctx context.Context, request *grpc_device_manager_go.AddDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error) {
	vErr := entities.ValidAddDeviceGroupRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.AddDeviceGroup(request)
}

func (h *Handler) UpdateDeviceGroup(ctx context.Context, request *grpc_device_manager_go.UpdateDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error) {
	vErr := entities.ValidUpdateDeviceGroupRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.UpdateDeviceGroup(request)
}

func (h *Handler) RemoveDeviceGroup(ctx context.Context, deviceGroupID *grpc_device_go.DeviceGroupId) (*grpc_common_go.Success, error) {
	vErr := entities.ValidDeviceGroupID(deviceGroupID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.RemoveDeviceGroup(deviceGroupID)
}

func (h *Handler) ListDeviceGroups(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_device_manager_go.DeviceGroupList, error) {
	vErr := entities.ValidOrganizationID(organizationID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.ListDeviceGroups(organizationID)
}

func (h *Handler) RegisterDevice(ctx context.Context, request *grpc_device_manager_go.RegisterDeviceRequest) (*grpc_device_manager_go.RegisterResponse, error) {
	vErr := entities.ValidRegisterDeviceRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.RegisterDevice(request)
}

func (h *Handler) GetDevice(ctx context.Context, deviceID *grpc_device_go.DeviceId) (*grpc_device_manager_go.Device, error) {
	vErr := entities.ValidDeviceID(deviceID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.GetDevice(deviceID)
}

func (h *Handler) ListDevices(ctx context.Context, deviceGroupID *grpc_device_go.DeviceGroupId) (*grpc_device_manager_go.DeviceList, error) {
	vErr := entities.ValidDeviceGroupID(deviceGroupID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.ListDevices(deviceGroupID)
}

func (h *Handler) AddLabelToDevice(ctx context.Context, request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error) {
	vErr := entities.ValidDeviceLabelRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.AddLabelToDevice(request)
}

func (h *Handler) RemoveLabelFromDevice(ctx context.Context, request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error) {
	vErr := entities.ValidDeviceLabelRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.RemoveLabelFromDevice(request)
}

func (h *Handler) UpdateDevice(ctx context.Context, request *grpc_device_manager_go.UpdateDeviceRequest) (*grpc_device_manager_go.Device, error) {
	vErr := entities.ValidUpdateDeviceRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.UpdateDevice(request)
}

func (h *Handler) UpdateDeviceLocation(ctx context.Context, request *grpc_device_manager_go.UpdateDeviceLocationRequest) (*grpc_device_manager_go.Device, error) {
	log.Debug().Interface("request", request).Msg("Update device location")
	vErr := entities.ValidUpdateDeviceLocationRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.UpdateDeviceLocation(request)
}

func (h *Handler) RemoveDevice(ctx context.Context, deviceID *grpc_device_go.DeviceId) (*grpc_common_go.Success, error) {
	vErr := entities.ValidDeviceID(deviceID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.RemoveDevice(deviceID)
}
