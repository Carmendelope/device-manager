/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package device

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"time"
)

const DeviceClientTimeout = time.Second * 5

// Manager structure with the required clients for node operations.
type Manager struct {
	authxClient grpc_authx_go.AuthxClient
	devicesClient grpc_device_go.DevicesClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(authxClient grpc_authx_go.AuthxClient, deviceClient grpc_device_go.DevicesClient) Manager {
	return Manager{
		authxClient: authxClient,
		devicesClient: deviceClient,
	}
}

// RetrieveDeviceLabels retrieves the list of labels associated with the current device.
/*
func (m*Manager) RetrieveDeviceLabels(deviceId *grpc_device_go.DeviceId) (*grpc_common_go.Labels, error){
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()
	retrieved, err := m.deviceClient.GetDevice(ctx, deviceId)
	if err != nil{
		return nil, err
	}
	return &grpc_common_go.Labels{
		Labels: retrieved.Labels,
	}, nil
}
*/

func (m*Manager) AddDeviceGroup(request *grpc_device_manager_go.AddDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented"))
}

func (m*Manager) UpdateDeviceGroup(request *grpc_device_manager_go.UpdateDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented"))
}
func (m*Manager) RemoveDeviceGroup(deviceGroupID *grpc_device_go.DeviceGroupId) (*grpc_common_go.Success, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented"))
}
func (m*Manager) ListDeviceGroups(organizationID *grpc_organization_go.OrganizationId) (*grpc_device_manager_go.DeviceGroupList, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented"))
}
func (m*Manager) RegisterDevice(request *grpc_device_manager_go.RegisterDeviceRequest) (*grpc_device_manager_go.RegisterResponse, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented"))
}
func (m*Manager) GetDevice(deviceID *grpc_device_go.DeviceId) (*grpc_device_manager_go.Device, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented"))
}
func (m*Manager) ListDevices(deviceGroupID *grpc_device_go.DeviceGroupId) (*grpc_device_manager_go.DeviceList, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented"))
}
func (m*Manager) AddLabelToDevice(request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented"))
}
func (m*Manager) RemoveLabelFromDevice(request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented"))
}
func (m*Manager) UpdateDevice(request *grpc_device_manager_go.UpdateDeviceRequest) (*grpc_device_manager_go.Device, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented"))
}
