/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
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