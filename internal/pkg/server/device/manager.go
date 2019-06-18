/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package device

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"github.com/nalej/device-manager/internal/pkg/provider/latency"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"time"
)

const DeviceClientTimeout = time.Second * 5
const AuthxClientTimeout = time.Second * 5

// Manager structure with the required clients for node operations.
type Manager struct {
	authxClient grpc_authx_go.AuthxClient
	devicesClient grpc_device_go.DevicesClient
	appsClient grpc_application_go.ApplicationsClient
	threshold time.Duration
	latencyProvider latency.Provider
}

// NewManager creates a Manager using a set of clients.
func NewManager(authxClient grpc_authx_go.AuthxClient, deviceClient grpc_device_go.DevicesClient,
	appsClient grpc_application_go.ApplicationsClient, lProvider latency.Provider, threshold time.Duration) Manager {
	return Manager{
		authxClient: authxClient,
		devicesClient: deviceClient,
		appsClient: appsClient,
		latencyProvider:lProvider,
		threshold: threshold,
	}
}

// Retrieve a DeviceGroup with all information.
func (m*Manager) GetDeviceGroup(deviceGroupID *grpc_device_go.DeviceGroupId) (*grpc_device_manager_go.DeviceGroup, error){
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()
	dg, err := m.devicesClient.GetDeviceGroup(ctx, deviceGroupID)
	if err != nil{
		return nil, err
	}
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	dgc, err := m.authxClient.GetDeviceGroupCredentials(aCtx, deviceGroupID)
	if err != nil{
		return nil, err
	}
	return &grpc_device_manager_go.DeviceGroup{
		OrganizationId:            dg.OrganizationId,
		DeviceGroupId:             dg.DeviceGroupId,
		Name:                      dg.Name,
		Created:                   dg.Created,
		Labels:                    dg.Labels,
		Enabled:                   dgc.Enabled,
		DefaultDeviceConnectivity: dgc.DefaultDeviceConnectivity,
		DeviceGroupApiKey:         dgc.DeviceGroupApiKey,
	}, nil
}

func (m*Manager) addAuthInfoToDG(dg *grpc_device_go.DeviceGroup) (*grpc_device_manager_go.DeviceGroup, error){
	deviceGroupID := &grpc_device_go.DeviceGroupId{
		OrganizationId:       dg.OrganizationId,
		DeviceGroupId:        dg.DeviceGroupId,
	}
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	dgc, err := m.authxClient.GetDeviceGroupCredentials(aCtx, deviceGroupID)
	if err != nil{
		return nil, err
	}
	return &grpc_device_manager_go.DeviceGroup{
		OrganizationId:            dg.OrganizationId,
		DeviceGroupId:             dg.DeviceGroupId,
		Name:                      dg.Name,
		Created:                   dg.Created,
		Labels:                    dg.Labels,
		Enabled:                   dgc.Enabled,
		DefaultDeviceConnectivity: dgc.DefaultDeviceConnectivity,
		DeviceGroupApiKey:         dgc.DeviceGroupApiKey,
	}, nil
}

func (m*Manager) AddDeviceGroup(request *grpc_device_manager_go.AddDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error){
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()
	addDGRequest := &grpc_device_go.AddDeviceGroupRequest{
		OrganizationId:       request.OrganizationId,
		Name:                 request.Name,
		Labels:               nil,
	}
	added, err := m.devicesClient.AddDeviceGroup(ctx, addDGRequest)
	if err != nil{
		return nil, err
	}
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	addDGCredentialsRequest := &grpc_authx_go.AddDeviceGroupCredentialsRequest{
		OrganizationId:            request.OrganizationId,
		DeviceGroupId:             added.DeviceGroupId,
		Enabled:                   request.Enabled,
		DefaultDeviceConnectivity: request.DefaultDeviceConnectivity,
	}
	credentials, err := m.authxClient.AddDeviceGroupCredentials(aCtx, addDGCredentialsRequest)
	if err != nil{
		return nil, err
	}
	log.Debug().Interface("deviceGroup", added).Msg("device group has been added")
	return &grpc_device_manager_go.DeviceGroup{
		OrganizationId:            added.OrganizationId,
		DeviceGroupId:             added.DeviceGroupId,
		Name:                      added.Name,
		Created:                   added.Created,
		Labels:                    added.Labels,
		Enabled:                   credentials.Enabled,
		DefaultDeviceConnectivity: credentials.DefaultDeviceConnectivity,
		DeviceGroupApiKey:         credentials.DeviceGroupApiKey,
	}, nil
}

func (m*Manager) UpdateDeviceGroup(request *grpc_device_manager_go.UpdateDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error){
	toUpdate := &grpc_authx_go.UpdateDeviceGroupCredentialsRequest{
		OrganizationId:            request.OrganizationId,
		DeviceGroupId:             request.DeviceGroupId,
		UpdateEnabled:             request.UpdateEnabled,
		Enabled:                   request.Enabled,
		UpdateDeviceConnectivity:  request.UpdateDeviceConnectivity,
		DefaultDeviceConnectivity: request.DefaultDeviceConnectivity,
	}
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	_, err := m.authxClient.UpdateDeviceGroupCredentials(aCtx, toUpdate)
	if err != nil{
		return nil, err
	}
	dgID := &grpc_device_go.DeviceGroupId{
		OrganizationId:       request.OrganizationId,
		DeviceGroupId:        request.DeviceGroupId,
	}
	return m.GetDeviceGroup(dgID)
}

func (m*Manager) disableDeviceGroup(deviceGroupID *grpc_device_go.DeviceGroupId) error {
	toUpdate := &grpc_authx_go.UpdateDeviceGroupCredentialsRequest{
		OrganizationId:            deviceGroupID.OrganizationId,
		DeviceGroupId:             deviceGroupID.DeviceGroupId,
		UpdateEnabled:             true,
		Enabled:                   false,
		UpdateDeviceConnectivity:  true,
		DefaultDeviceConnectivity: false,
	}
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	_, err := m.authxClient.UpdateDeviceGroupCredentials(aCtx, toUpdate)
	if err != nil{
		return err
	}
	log.Debug().Interface("deviceGroupID", deviceGroupID).Msg("device group has been disabled")
	return nil
}

func (m*Manager) deviceGroupHasApps(deviceGroupID *grpc_device_go.DeviceGroupId) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()
	organizationID := &grpc_organization_go.OrganizationId{
		OrganizationId:       deviceGroupID.OrganizationId,
	}

	descriptors, err := m.appsClient.ListAppDescriptors(ctx, organizationID)
	if err != nil{
		return false, err
	}

	for _, desc := range descriptors.Descriptors{
		for _, rule := range desc.Rules{
			if rule.Access == grpc_application_go.PortAccess_DEVICE_GROUP{
				for _, dg := range rule.DeviceGroupIds{
					if dg == deviceGroupID.DeviceGroupId{
						return true, nil
					}
				}
			}
		}

	}
	return false, nil
}

func (m*Manager) removeDevice(deviceID *grpc_device_go.DeviceId) error {
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	_, err := m.authxClient.RemoveDeviceCredentials(aCtx, deviceID)
	if err != nil{
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()
	removeRequest := &grpc_device_go.RemoveDeviceRequest{
		OrganizationId:       deviceID.OrganizationId,
		DeviceGroupId:        deviceID.DeviceGroupId,
		DeviceId:             deviceID.DeviceId,
	}
	_, err = m.devicesClient.RemoveDevice(ctx, removeRequest)
	if err != nil{
		return err
	}
	log.Debug().Interface("deviceID", deviceID).Msg("device has been removed")
	return nil
}

func (m*Manager) removeDeviceGroupEntity(deviceGroupID *grpc_device_go.DeviceGroupId) error {
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	_, err := m.authxClient.RemoveDeviceGroupCredentials(aCtx, deviceGroupID)
	if err != nil{
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()
	removeRequest := &grpc_device_go.RemoveDeviceGroupRequest{
		OrganizationId:       deviceGroupID.OrganizationId,
		DeviceGroupId:        deviceGroupID.DeviceGroupId,
	}
	_, err = m.devicesClient.RemoveDeviceGroup(ctx, removeRequest)
	if err != nil{
		return err
	}
	log.Debug().Interface("deviceGroupID", deviceGroupID).Msg("device group entity has been removed")
	return nil
}

// If a device group is removed, we must remove all the attached devices.
func (m*Manager) RemoveDeviceGroup(deviceGroupID *grpc_device_go.DeviceGroupId) (*grpc_common_go.Success, error){
	apps, err := m.deviceGroupHasApps(deviceGroupID)
	if err != nil{
		log.Debug().Msg("cannot determine if a device group as it has apps linked to it")
		return nil, err
	}
	if apps{
		return nil, derrors.NewFailedPreconditionError("device group has application descriptors linked to it")
	}
	// First deactivate the device group so that sensors will
	err = m.disableDeviceGroup(deviceGroupID)
	if err != nil{
		log.Debug().Msg("cannot disable device group")
		return nil, err
	}

	devices, err := m.ListDevices(deviceGroupID)
	if err != nil{
		log.Debug().Msg("cannot retrieve the list of devices to be removed")
		return nil, err
	}
	log.Debug().Msg("removing devices")
	for _, d := range devices.Devices{
		dID := &grpc_device_go.DeviceId{
			OrganizationId:       d.OrganizationId,
			DeviceGroupId:        d.DeviceGroupId,
			DeviceId:             d.DeviceId,
		}
		err := m.removeDevice(dID)
		if err != nil{
			return nil, err
		}
	}
	log.Debug().Msg("removing device group")
	err = m.removeDeviceGroupEntity(deviceGroupID)
	if err != nil{
		return nil, err
	}
	log.Debug().Msg("device has been removed")
	return &grpc_common_go.Success{}, nil
}

func (m*Manager) ListDeviceGroups(organizationID *grpc_organization_go.OrganizationId) (*grpc_device_manager_go.DeviceGroupList, error){
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()
	dgs, err := m.devicesClient.ListDeviceGroups(ctx, organizationID)
	if err != nil{
		return nil, err
	}
	result := make([]*grpc_device_manager_go.DeviceGroup, 0)
	for _, dg := range dgs.Groups {
		toAdd, err := m.addAuthInfoToDG(dg)
		if err != nil{
			return nil, err
		}
		result = append(result, toAdd)
	}
	return &grpc_device_manager_go.DeviceGroupList{
		Groups:               result,
	}, nil
}

func (m*Manager) addDeviceEntity(request *grpc_device_manager_go.RegisterDeviceRequest) (*grpc_device_manager_go.RegisterResponse, error){
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()
	addRequest := &grpc_device_go.AddDeviceRequest{
		OrganizationId:       request.OrganizationId,
		DeviceGroupId:        request.DeviceGroupId,
		DeviceId:             request.DeviceId,
		Labels:               request.Labels,
		AssetInfo:            request.AssetInfo,
	}
	added, err := m.devicesClient.AddDevice(ctx, addRequest)
	if err != nil{
		return nil, err
	}
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	addCredentialsRequest := &grpc_authx_go.AddDeviceCredentialsRequest{
		OrganizationId:       request.OrganizationId,
		DeviceGroupId:        request.DeviceGroupId,
		DeviceId:             request.DeviceId,
	}
	credentials, err := m.authxClient.AddDeviceCredentials(aCtx, addCredentialsRequest)
	if err != nil{
		return nil, err
	}
	log.Debug().Interface("device", added).Msg("device has been added")
	return &grpc_device_manager_go.RegisterResponse{
		DeviceId:             credentials.DeviceId,
		DeviceApiKey:         credentials.DeviceApiKey,
	}, nil
}

func (m*Manager) RegisterDevice(request *grpc_device_manager_go.RegisterDeviceRequest) (*grpc_device_manager_go.RegisterResponse, error){
	// Check that the device group is usable
	log.Debug().Interface("request", request).Msg("adding device")
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	dgLoginRequest := &grpc_authx_go.DeviceGroupLoginRequest{
		OrganizationId:       request.OrganizationId,
		DeviceGroupApiKey:    request.DeviceGroupApiKey,
	}
	_, err := m.authxClient.DeviceGroupLogin(aCtx, dgLoginRequest)
	if err != nil{
		return nil, err
	}
	log.Debug().Str("deviceID", request.DeviceId).Msg("device group is valid")
	// Add the device
	return m.addDeviceEntity(request)
}

func (m*Manager) GetDevice(deviceID *grpc_device_go.DeviceId) (*grpc_device_manager_go.Device, error){
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()
	d, err := m.devicesClient.GetDevice(ctx, deviceID)
	if err != nil{
		return nil, err
	}
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	dc, err := m.authxClient.GetDeviceCredentials(aCtx, deviceID)

	status := grpc_device_manager_go.DeviceStatus_OFFLINE
	latency, err := m.latencyProvider.GetLastLatency(d.OrganizationId, d.DeviceGroupId, d.DeviceId)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error getting device latency")
	}else{
		status = m.fillDeviceStatus(latency)
	}


	//status := m.fillDeviceStatus(d.OrganizationId, d.DeviceGroupId, d.DeviceId )

	return &grpc_device_manager_go.Device{
		OrganizationId:       d.OrganizationId,
		DeviceGroupId:        d.DeviceGroupId,
		DeviceId:             d.DeviceId,
		RegisterSince:        d.RegisterSince,
		Labels:               d.Labels,
		Enabled:              dc.Enabled,
		DeviceApiKey:         dc.DeviceApiKey,
		DeviceStatus: 		  status,
		AssetInfo:            d.AssetInfo,
	}, nil
}

func (m * Manager) fillDeviceStatus (latency *entities.Latency) grpc_device_manager_go.DeviceStatus  { //(OrganizationId string, DeviceGroupId string, DeviceId string) grpc_device_manager_go.DeviceStatus  {
	status := grpc_device_manager_go.DeviceStatus_OFFLINE
	if latency != nil && latency.Latency != -1 { // if latency == -1 -> no ping found (no error, the device is OFFLINE)
		timeCalculated := time.Unix(latency.Inserted, 0).Add(m.threshold).Unix()
		if timeCalculated > time.Now().Unix(){
			status = grpc_device_manager_go.DeviceStatus_ONLINE
		}
	}
	return status
}

func (m*Manager) addAuthInfoToD(dg *grpc_device_go.Device) (*grpc_device_manager_go.Device, error){
	deviceID := &grpc_device_go.DeviceId{
		OrganizationId:       dg.OrganizationId,
		DeviceGroupId: dg.DeviceGroupId,
		DeviceId:        dg.DeviceId,
	}
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	dc, err := m.authxClient.GetDeviceCredentials(aCtx, deviceID)
	if err != nil{
		return nil, err
	}
	return &grpc_device_manager_go.Device{
		OrganizationId:       dg.OrganizationId,
		DeviceGroupId:        dg.DeviceGroupId,
		DeviceId:             dg.DeviceId,
		RegisterSince:        dg.RegisterSince,
		Labels:               dg.Labels,
		Enabled:              dc.Enabled,
		DeviceApiKey:         dc.DeviceApiKey,
		DeviceStatus:         grpc_device_manager_go.DeviceStatus_OFFLINE, // offline by default
		AssetInfo:            dg.AssetInfo,
	}, nil
}

func (m*Manager)updateStatus (devices []*grpc_device_manager_go.Device, latency entities.Latency) {
	for i:= 0; i< len(devices); i++ {
		if devices[i].DeviceId == latency.DeviceId {
			devices[i].DeviceStatus = m.fillDeviceStatus(&latency)
			//timeCalculated := time.Unix(latency.Inserted, 0).Add(m.threshold).Unix()
			//if timeCalculated > time.Now().Unix(){
		//		devices[i].DeviceStatus = grpc_device_manager_go.DeviceStatus_ONLINE
		//	}
		}
	}
}

func (m*Manager) ListDevices(deviceGroupID *grpc_device_go.DeviceGroupId) (*grpc_device_manager_go.DeviceList, error){
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()
	devices, err := m.devicesClient.ListDevices(ctx, deviceGroupID)
	if err != nil{
		return nil, err
	}
	result := make([]*grpc_device_manager_go.Device, 0)
	for _, d := range devices.Devices{
		toAdd, err := m.addAuthInfoToD(d)
		if err != nil{
			return nil, err
		}
		result = append(result, toAdd)
	}

	if len(result) > 0 {
		latencies, err := m.latencyProvider.GetGroupLastLatencies(deviceGroupID.OrganizationId, deviceGroupID.DeviceGroupId)
		if err != nil {
			log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error getting group latencies")
		}else{
			for _, latency := range latencies {
				m.updateStatus(result, *latency)
			}
		}
	}

	return &grpc_device_manager_go.DeviceList{
		Devices:              result,
	}, nil

}

func (m*Manager) AddLabelToDevice(request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error){
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()

	_, err := m.devicesClient.UpdateDevice(ctx, &grpc_device_go.UpdateDeviceRequest{
		OrganizationId: request.OrganizationId,
		DeviceGroupId: request.DeviceGroupId,
		DeviceId: request.DeviceId,
		AddLabels: true,
		RemoveLabels: false,
		Labels: request.Labels,
	})
	if err != nil {
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

func (m*Manager) RemoveLabelFromDevice(request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error){
	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()

	_, err := m.devicesClient.UpdateDevice(ctx, &grpc_device_go.UpdateDeviceRequest{
		OrganizationId: request.OrganizationId,
		DeviceGroupId: request.DeviceGroupId,
		DeviceId: request.DeviceId,
		AddLabels: false,
		RemoveLabels: true,
		Labels: request.Labels,
	})
	if err != nil {
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

func (m*Manager) UpdateDevice(request *grpc_device_manager_go.UpdateDeviceRequest) (*grpc_device_manager_go.Device, error){
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()
	updateRequest := &grpc_authx_go.UpdateDeviceCredentialsRequest{
		OrganizationId:       request.OrganizationId,
		DeviceGroupId:        request.DeviceGroupId,
		DeviceId:             request.DeviceId,
		Enabled:              request.Enabled,
	}
	_, err := m.authxClient.UpdateDeviceCredentials(aCtx, updateRequest)
	if err != nil{
		return nil, err
	}
	deviceID := &grpc_device_go.DeviceId{
		OrganizationId:       request.OrganizationId,
		DeviceGroupId:        request.DeviceGroupId,
		DeviceId:             request.DeviceId,
	}

	device, err := m.GetDevice(deviceID)
	device.Location = request.Location

	return device, err
}

func (m*Manager) RemoveDevice(deviceID *grpc_device_go.DeviceId) (*grpc_common_go.Success, error){
	aCtx, aCancel := context.WithTimeout(context.Background(), AuthxClientTimeout)
	defer aCancel()

	_, err := m.authxClient.RemoveDeviceCredentials(aCtx, deviceID)
	if err != nil{
		log.Warn().Interface("deviceID", deviceID).Msg("Device may be partially removed. Cannot remove auth information")
	}


	err = m.latencyProvider.RemoveLatency(deviceID.OrganizationId, deviceID.DeviceGroupId, deviceID.DeviceId)
	if err != nil {
		log.Warn().Interface("deviceID", deviceID).Msg("Device may be partially removed. Cannot remove device latencies")
	}

	ctx, cancel := context.WithTimeout(context.Background(), DeviceClientTimeout)
	defer cancel()

	removeRequest := &grpc_device_go.RemoveDeviceRequest{
		OrganizationId:       deviceID.OrganizationId,
		DeviceGroupId:        deviceID.DeviceGroupId,
		DeviceId:             deviceID.DeviceId,
	}

	return m.devicesClient.RemoveDevice(ctx, removeRequest)
}