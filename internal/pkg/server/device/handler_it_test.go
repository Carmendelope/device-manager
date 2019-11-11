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
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"github.com/nalej/device-manager/internal/pkg/provider/latency"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"math/rand"
	"os"
	"time"
)

func GetConnection(address string) (* grpc.ClientConn) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	gomega.Expect(err).To(gomega.Succeed())
	return conn
}

func CreateOrganization(name string, orgClient grpc_organization_go.OrganizationsClient) * grpc_organization_go.Organization {
	toAdd := &grpc_organization_go.AddOrganizationRequest{
		Name:                 fmt.Sprintf("%s-%d-%d", name, ginkgo.GinkgoRandomSeed(), rand.Int()),
	}
	added, err := orgClient.AddOrganization(context.Background(), toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	gomega.Expect(added).ToNot(gomega.BeNil())
	return added
}

func CreateDeviceGroup(client grpc_device_manager_go.DevicesClient, organizationID string, enabled bool, defaultConnectivity bool) * grpc_device_manager_go.DeviceGroup{
	addDGRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
		OrganizationId:            organizationID,
		Name:                      fmt.Sprintf("dg-%d", rand.Int()),
		Enabled:                   enabled,
		DefaultDeviceConnectivity: defaultConnectivity,
	}
	added, err := client.AddDeviceGroup(context.Background(), addDGRequest)
	gomega.Expect(err).To(gomega.Succeed())
	gomega.Expect(added).ShouldNot(gomega.BeNil())
	return added
}

var _ = ginkgo.Describe("Device service", func() {

	var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")

	if runIntegration != "true" {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress= os.Getenv("IT_SM_ADDRESS")
		authxAddress= os.Getenv("IT_AUTHX_ADDRESS")
	)

	if systemModelAddress == "" || authxAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_device_manager_go.DevicesClient

	// Providers
	var orgClient grpc_organization_go.OrganizationsClient
	var deviceClient grpc_device_go.DevicesClient
	var appClient grpc_application_go.ApplicationsClient
	var smConn *grpc.ClientConn
	var authxClient grpc_authx_go.AuthxClient
	var authxConn *grpc.ClientConn
	var latencyProvider *latency.MockupProvider

	// Target organization.
	var targetOrganization *grpc_organization_go.Organization

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		smConn = GetConnection(systemModelAddress)
		deviceClient = grpc_device_go.NewDevicesClient(smConn)
		appClient = grpc_application_go.NewApplicationsClient(smConn)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)

		authxConn = GetConnection(authxAddress)
		authxClient = grpc_authx_go.NewAuthxClient(authxConn)

		// provider
		latencyProvider = latency.NewMockupProvider()

		// Register the service
		d, _ := time.ParseDuration("3m")

		manager := NewManager(authxClient, deviceClient, appClient, latencyProvider, d)
		handler := NewHandler(manager)
		grpc_device_manager_go.RegisterDevicesServer(server, handler)
		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_device_manager_go.NewDevicesClient(conn)
		rand.Seed(ginkgo.GinkgoRandomSeed())
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func() {
		ginkgo.By("creating target entities", func() {
			// Initial data
			targetOrganization = CreateOrganization("device-manager-it", orgClient)
		})
	})

	ginkgo.Context("device groups", func(){
		ginkgo.It("should be able to create a device group", func(){
			addDGRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
				OrganizationId:            targetOrganization.OrganizationId,
				Name:                      fmt.Sprintf("dg-%d", rand.Int()),
				Enabled:                   false,
				DefaultDeviceConnectivity: false,
			}
			added, err := client.AddDeviceGroup(context.Background(), addDGRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.DeviceGroupApiKey).ShouldNot(gomega.BeEmpty())
		})
		ginkgo.It("should be able to list device groups", func(){
			addDGRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
				OrganizationId:            targetOrganization.OrganizationId,
				Name:                      fmt.Sprintf("dg-%d", rand.Int()),
				Enabled:                   false,
				DefaultDeviceConnectivity: false,
			}
			added, err := client.AddDeviceGroup(context.Background(), addDGRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())

			orgID := &grpc_organization_go.OrganizationId{
				OrganizationId:       targetOrganization.OrganizationId,
			}
			list, err := client.ListDeviceGroups(context.Background(), orgID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list.Groups)).Should(gomega.Equal(1))
			retrieved := list.Groups[0]
			gomega.Expect(retrieved.DeviceGroupApiKey).Should(gomega.Equal(added.DeviceGroupApiKey))
		})
		ginkgo.It("should be able to update a device group", func(){
			addDGRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
				OrganizationId:            targetOrganization.OrganizationId,
				Name:                      fmt.Sprintf("dg-%d", rand.Int()),
				Enabled:                   false,
				DefaultDeviceConnectivity: false,
			}
			added, err := client.AddDeviceGroup(context.Background(), addDGRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())

			updateRequest := &grpc_device_manager_go.UpdateDeviceGroupRequest{
				OrganizationId:            added.OrganizationId,
				DeviceGroupId:             added.DeviceGroupId,
				UpdateEnabled:             true,
				Enabled:                   true,
				UpdateDeviceConnectivity:  true,
				DefaultDeviceConnectivity: true,
			}
			updated, err := client.UpdateDeviceGroup(context.Background(), updateRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(updated.Enabled).Should(gomega.BeTrue())
			gomega.Expect(updated.DefaultDeviceConnectivity).Should(gomega.BeTrue())
		})
		ginkgo.It("should be able to remove a device group", func(){
			addDGRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
				OrganizationId:            targetOrganization.OrganizationId,
				Name:                      fmt.Sprintf("dg-%d", rand.Int()),
				Enabled:                   false,
				DefaultDeviceConnectivity: false,
			}
			added, err := client.AddDeviceGroup(context.Background(), addDGRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			toRemove := &grpc_device_go.DeviceGroupId{
				OrganizationId:       targetOrganization.OrganizationId,
				DeviceGroupId:        added.DeviceGroupId,
			}
			success, err := client.RemoveDeviceGroup(context.Background(), toRemove)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())
		})
		ginkgo.PIt("should not be able to remove a device group with associated app descriptors", func(){

		})

	})
	ginkgo.Context("devices", func(){
		ginkgo.It("should be able to register a device", func(){
			dg := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, true)
			registerRequest := &grpc_device_manager_go.RegisterDeviceRequest{
				OrganizationId:       dg.OrganizationId,
				DeviceGroupId:        dg.DeviceGroupId,
				DeviceGroupApiKey:    dg.DeviceGroupApiKey,
				DeviceId:             fmt.Sprintf("d-%s-%d", dg.DeviceGroupId, rand.Int()),
				Labels:               nil,
			}
			added, err := client.RegisterDevice(context.Background(), registerRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.DeviceApiKey).ShouldNot(gomega.BeNil())
		})
		ginkgo.It("new devices should follow the device group policy", func(){
			// first try with default connectivity disabled
			dgF := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, false)
			registerRequest1 := &grpc_device_manager_go.RegisterDeviceRequest{
				OrganizationId:       dgF.OrganizationId,
				DeviceGroupId:        dgF.DeviceGroupId,
				DeviceGroupApiKey:    dgF.DeviceGroupApiKey,
				DeviceId:             fmt.Sprintf("d-%s-%d", dgF.DeviceGroupId, rand.Int()),
				Labels:               nil,
			}
			added1, err := client.RegisterDevice(context.Background(), registerRequest1)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added1).ShouldNot(gomega.BeNil())
			deviceID1 := &grpc_device_go.DeviceId{
				OrganizationId:       dgF.OrganizationId,
				DeviceGroupId:        dgF.DeviceGroupId,
				DeviceId:             added1.DeviceId,
			}
			retrieved1, err := client.GetDevice(context.Background(), deviceID1)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved1.Enabled).Should(gomega.BeFalse())
			// then, try with default connectivity enabled
			dgT := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, true)
			registerRequest2 := &grpc_device_manager_go.RegisterDeviceRequest{
				OrganizationId:       dgT.OrganizationId,
				DeviceGroupId:        dgT.DeviceGroupId,
				DeviceGroupApiKey:    dgT.DeviceGroupApiKey,
				DeviceId:             fmt.Sprintf("d-%s-%d", dgT.DeviceGroupId, rand.Int()),
				Labels:               nil,
			}
			added2, err := client.RegisterDevice(context.Background(), registerRequest2)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added2).ShouldNot(gomega.BeNil())
			deviceID2 := &grpc_device_go.DeviceId{
				OrganizationId:       dgT.OrganizationId,
				DeviceGroupId:        dgT.DeviceGroupId,
				DeviceId:             added2.DeviceId,
			}
			retrieved2, err := client.GetDevice(context.Background(), deviceID2)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved2.Enabled).Should(gomega.BeTrue())
		})
		ginkgo.PIt("should not be able to register a device with incorrect DGAK", func(){

		})
		ginkgo.PIt("should not be able to register a device with incorrect device group id", func(){

		})
		ginkgo.It("should be able to retrieve a device", func(){
			dg := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, true)
			registerRequest := &grpc_device_manager_go.RegisterDeviceRequest{
				OrganizationId:       dg.OrganizationId,
				DeviceGroupId:        dg.DeviceGroupId,
				DeviceGroupApiKey:    dg.DeviceGroupApiKey,
				DeviceId:             fmt.Sprintf("d-%s-%d", dg.DeviceGroupId, rand.Int()),
				Labels:               nil,
			}
			added, err := client.RegisterDevice(context.Background(), registerRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())

			deviceID := &grpc_device_go.DeviceId{
				OrganizationId:       dg.OrganizationId,
				DeviceGroupId:        dg.DeviceGroupId,
				DeviceId:             added.DeviceId,
			}
			retrieved, err := client.GetDevice(context.Background(), deviceID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved.DeviceApiKey).Should(gomega.Equal(added.DeviceApiKey))
		})
		ginkgo.It("should be able to list devices", func(){
			dg := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, true)
			registerRequest := &grpc_device_manager_go.RegisterDeviceRequest{
				OrganizationId:       dg.OrganizationId,
				DeviceGroupId:        dg.DeviceGroupId,
				DeviceGroupApiKey:    dg.DeviceGroupApiKey,
				DeviceId:             fmt.Sprintf("d-%s-%d", dg.DeviceGroupId, rand.Int()),
				Labels:               nil,
			}
			added, err := client.RegisterDevice(context.Background(), registerRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())

			deviceGroupID := &grpc_device_go.DeviceGroupId{
				OrganizationId:       dg.OrganizationId,
				DeviceGroupId:        dg.DeviceGroupId,
			}
			list, err := client.ListDevices(context.Background(), deviceGroupID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list.Devices)).Should(gomega.Equal(1))
			retrieved := list.Devices[0]
			gomega.Expect(retrieved.DeviceApiKey).Should(gomega.Equal(added.DeviceApiKey))
		})
		ginkgo.It("should be able to add a label to a device", func(){
			dg := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, true)
			registerRequest := &grpc_device_manager_go.RegisterDeviceRequest{
				OrganizationId:       dg.OrganizationId,
				DeviceGroupId:        dg.DeviceGroupId,
				DeviceGroupApiKey:    dg.DeviceGroupApiKey,
				DeviceId:             fmt.Sprintf("d-%s-%d", dg.DeviceGroupId, rand.Int()),
				Labels:               nil,
			}
			added, err := client.RegisterDevice(context.Background(), registerRequest)
			gomega.Expect(err).To(gomega.Succeed())

			success, err := client.AddLabelToDevice(context.Background(), &grpc_device_manager_go.DeviceLabelRequest{
				OrganizationId: dg.OrganizationId,
				DeviceGroupId:dg.DeviceGroupId,
				DeviceId: added.DeviceId,
				Labels: map[string]string{"label1":"value1", "label2":"value2"},
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

			// check if the update works

			device, err := client.GetDevice(context.Background(), &grpc_device_go.DeviceId{
				OrganizationId: dg.OrganizationId,
				DeviceGroupId: dg.DeviceGroupId,
				DeviceId: added.DeviceId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(device).NotTo(gomega.BeNil())
			gomega.Expect(len(device.Labels)).Should(gomega.Equal(2))
		})
		ginkgo.It("should be able to remove a label from a device", func(){
			dg := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, true)
			registerRequest := &grpc_device_manager_go.RegisterDeviceRequest{
				OrganizationId:       dg.OrganizationId,
				DeviceGroupId:        dg.DeviceGroupId,
				DeviceGroupApiKey:    dg.DeviceGroupApiKey,
				DeviceId:             fmt.Sprintf("d-%s-%d", dg.DeviceGroupId, rand.Int()),
				Labels:               map[string]string{"label1":"value1", "label2":"value2"},
			}
			added, err := client.RegisterDevice(context.Background(), registerRequest)
			gomega.Expect(err).To(gomega.Succeed())

			success, err := client.RemoveLabelFromDevice(context.Background(), &grpc_device_manager_go.DeviceLabelRequest{
				OrganizationId: dg.OrganizationId,
				DeviceGroupId:dg.DeviceGroupId,
				DeviceId: added.DeviceId,
				Labels: map[string]string{"label1":"value1"},
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

			// check if the update works

			device, err := client.GetDevice(context.Background(), &grpc_device_go.DeviceId{
				OrganizationId: dg.OrganizationId,
				DeviceGroupId: dg.DeviceGroupId,
				DeviceId: added.DeviceId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(device).NotTo(gomega.BeNil())
			gomega.Expect(len(device.Labels)).Should(gomega.Equal(1))
		})
		ginkgo.It("should not be able to remove a label from a non exiting device", func(){
			dg := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, true)

			success, err := client.RemoveLabelFromDevice(context.Background(), &grpc_device_manager_go.DeviceLabelRequest{
				OrganizationId: dg.OrganizationId,
				DeviceGroupId:dg.DeviceGroupId,
				DeviceId: uuid.New().String(),
				Labels: map[string]string{"label1":"value1"},
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())

		})
		ginkgo.It("should be able to update a device", func(){
			dg := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, false)
			registerRequest := &grpc_device_manager_go.RegisterDeviceRequest{
				OrganizationId:       dg.OrganizationId,
				DeviceGroupId:        dg.DeviceGroupId,
				DeviceGroupApiKey:    dg.DeviceGroupApiKey,
				DeviceId:             fmt.Sprintf("d-%s-%d", dg.DeviceGroupId, rand.Int()),
				Labels:               nil,
			}
			added, err := client.RegisterDevice(context.Background(), registerRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())

			updateRequest := &grpc_device_manager_go.UpdateDeviceRequest{
				OrganizationId:       dg.OrganizationId,
				DeviceGroupId:        dg.DeviceGroupId,
				DeviceId:             added.DeviceId,
				Enabled:              true,
			}
			updated, err := client.UpdateDevice(context.Background(), updateRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(updated.Enabled).Should(gomega.BeTrue())
		})
	})

	ginkgo.Context("interaction device group and device", func(){
		ginkgo.PIt("should remove devices on device group removal", func(){

		})
	})

	ginkgo.Context("Checking the device status", func() {
		ginkgo.It("should be able to get the device (Status OFFLINE)", func() {
			dg := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, true)
			registerRequest := &grpc_device_manager_go.RegisterDeviceRequest{
				OrganizationId:    dg.OrganizationId,
				DeviceGroupId:     dg.DeviceGroupId,
				DeviceGroupApiKey: dg.DeviceGroupApiKey,
				DeviceId:          fmt.Sprintf("d-%s-%d", dg.DeviceGroupId, rand.Int()),
				Labels:            nil,
			}
			added, err := client.RegisterDevice(context.Background(), registerRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.DeviceApiKey).ShouldNot(gomega.BeNil())

			// adding a ping
			ping := entities.Latency{
				OrganizationId: dg.OrganizationId,
				DeviceGroupId:  dg.DeviceGroupId,
				DeviceId:       registerRequest.DeviceId,
				Latency:        30,
				Inserted:       time.Now().Add(- time.Duration(4)*time.Minute).Unix()  ,
			}
			err = latencyProvider.AddPingLatency(ping)
			gomega.Expect(err).To(gomega.Succeed())

			toRetrieve := &grpc_device_go.DeviceId{
				OrganizationId: dg.OrganizationId,
				DeviceGroupId:  dg.DeviceGroupId,
				DeviceId:       registerRequest.DeviceId,
			}
			retrieved, err := client.GetDevice(context.Background(), toRetrieve)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved.DeviceStatus).Should(gomega.Equal(grpc_device_manager_go.DeviceStatus_OFFLINE))
		})
		ginkgo.It("should be able to get the device (Status ONLINE)", func() {
			dg := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, true)
			registerRequest := &grpc_device_manager_go.RegisterDeviceRequest{
				OrganizationId:    dg.OrganizationId,
				DeviceGroupId:     dg.DeviceGroupId,
				DeviceGroupApiKey: dg.DeviceGroupApiKey,
				DeviceId:          fmt.Sprintf("d-%s-%d", dg.DeviceGroupId, rand.Int()),
				Labels:            nil,
			}
			added, err := client.RegisterDevice(context.Background(), registerRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.DeviceApiKey).ShouldNot(gomega.BeNil())

			// adding a ping
			ping := entities.Latency{
				OrganizationId: dg.OrganizationId,
				DeviceGroupId:  dg.DeviceGroupId,
				DeviceId:       registerRequest.DeviceId,
				Latency:        30,
				Inserted:       time.Now().Add(- time.Duration(2)*time.Minute).Unix()  ,
			}
			err = latencyProvider.AddPingLatency(ping)
			gomega.Expect(err).To(gomega.Succeed())

			toRetrieve := &grpc_device_go.DeviceId{
				OrganizationId: dg.OrganizationId,
				DeviceGroupId:  dg.DeviceGroupId,
				DeviceId:       registerRequest.DeviceId,
			}
			retrieved, err := client.GetDevice(context.Background(), toRetrieve)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved.DeviceStatus).Should(gomega.Equal(grpc_device_manager_go.DeviceStatus_ONLINE))
		})
		ginkgo.It("should be able to list devices with the correct status (ONLINE/OFFLINE)", func(){
			dg := CreateDeviceGroup(client, targetOrganization.OrganizationId, true, true)
			for i:=1; i<=2; i++ {
				registerRequest := &grpc_device_manager_go.RegisterDeviceRequest{
					OrganizationId:    dg.OrganizationId,
					DeviceGroupId:     dg.DeviceGroupId,
					DeviceGroupApiKey: dg.DeviceGroupApiKey,
					DeviceId:          fmt.Sprintf("d-%s-%d", dg.DeviceGroupId, i),
					Labels:            nil,
				}
				added, err := client.RegisterDevice(context.Background(), registerRequest)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
			}

			// add some pings (device1)
			for i:= 0; i<=5; i++ {
				ping := entities.Latency{
					OrganizationId: dg.OrganizationId,
					DeviceGroupId:  dg.DeviceGroupId,
					DeviceId:       fmt.Sprintf("d-%s-1", dg.DeviceGroupId),
					Latency:        30,
					Inserted:       time.Now().Unix()+int64(i),
				}
				errLat := latencyProvider.AddPingLatency(ping)
				gomega.Expect(errLat).To(gomega.Succeed())
			}
			// add a ping (expired ping for device2)
			ping := entities.Latency{
				OrganizationId: dg.OrganizationId,
				DeviceGroupId:  dg.DeviceGroupId,
				DeviceId:       fmt.Sprintf("d-%s-2", dg.DeviceGroupId),
				Latency:        30,
				Inserted:       time.Now().Add(-time.Duration(5)*time.Hour).Unix()  ,
			}
			errLat := latencyProvider.AddPingLatency(ping)
			gomega.Expect(errLat).To(gomega.Succeed())

			// get devices
			deviceGroupID := &grpc_device_go.DeviceGroupId{
				OrganizationId:       dg.OrganizationId,
				DeviceGroupId:        dg.DeviceGroupId,
			}

			list, err := client.ListDevices(context.Background(), deviceGroupID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			for _, device := range list.Devices {
				if device.DeviceId == fmt.Sprintf("d-%s-1", dg.DeviceGroupId) {
					gomega.Expect(device.DeviceStatus).Should(gomega.Equal(grpc_device_manager_go.DeviceStatus_ONLINE))
				}else{
					gomega.Expect(device.DeviceStatus).Should(gomega.Equal(grpc_device_manager_go.DeviceStatus_OFFLINE))
				}
			}

		})
	})


})