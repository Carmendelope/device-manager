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

package server

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/device-manager/internal/pkg/provider/latency"
	"github.com/nalej/device-manager/internal/pkg/server/device"
	lat "github.com/nalej/device-manager/internal/pkg/server/latency"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

// Service structure with the configuration and the gRPC server.
type Service struct {
	Configuration Config
}

// NewService creates a new system model service.
func NewService(conf Config) *Service {
	return &Service{
		conf,
	}
}

// Providers structure with all the providers in the system.
type Providers struct {
	pProvider latency.Provider
}

// CreateInMemoryProviders returns a set of in-memory providers.
func (s *Service) CreateInMemoryProviders() *Providers {
	return &Providers{
		pProvider: latency.NewMockupProvider(),
	}
}

// CreateDBScyllaProviders returns a set of in-memory providers.
func (s *Service) CreateDBScyllaProviders() *Providers {
	return &Providers{
		pProvider: latency.NewScyllaProvider(
			s.Configuration.ScyllaDBAddress, s.Configuration.ScyllaDBPort, s.Configuration.KeySpace),
	}
}

// GetProviders builds the providers according to the selected backend.
func (s *Service) GetProviders() *Providers {
	if s.Configuration.UseInMemoryProviders {
		return s.CreateInMemoryProviders()
	} else if s.Configuration.UseDBScyllaProviders {
		return s.CreateDBScyllaProviders()
	}
	log.Fatal().Msg("unsupported type of provider")
	return nil
}

// Clients structure with the gRPC clients for remote services.
type Clients struct {
	AuthxClient   grpc_authx_go.AuthxClient
	DevicesClient grpc_device_go.DevicesClient
	AppsClient    grpc_application_go.ApplicationsClient
}

// GetClients creates the required connections with the remote clients.
func (s *Service) GetClients() (*Clients, derrors.Error) {
	authxConn, err := grpc.Dial(s.Configuration.AuthxAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the authx component")
	}

	smConn, err := grpc.Dial(s.Configuration.SystemModelAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the system model component")
	}

	aClient := grpc_authx_go.NewAuthxClient(authxConn)
	dClient := grpc_device_go.NewDevicesClient(smConn)
	appsClient := grpc_application_go.NewApplicationsClient(smConn)

	return &Clients{aClient, dClient, appsClient}, nil
}

// Run the service, launch the REST service handler.
func (s *Service) Run() error {
	cErr := s.Configuration.Validate()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("invalid configuration")
	}
	s.Configuration.Print()
	clients, cErr := s.GetClients()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("Cannot create clients")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	prov := s.GetProviders()

	// Create handlers
	manager := device.NewManager(clients.AuthxClient, clients.DevicesClient, clients.AppsClient, prov.pProvider, s.Configuration.Threshold)
	handler := device.NewHandler(manager)

	pManager := lat.NewManager(prov.pProvider)
	pHandler := lat.NewHandler(pManager)

	grpcServer := grpc.NewServer()

	grpc_device_manager_go.RegisterDevicesServer(grpcServer, handler)
	grpc_device_manager_go.RegisterLatencyServer(grpcServer, pHandler)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}
