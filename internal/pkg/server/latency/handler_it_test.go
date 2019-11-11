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
	"github.com/google/uuid"
	"github.com/nalej/device-manager/internal/pkg/provider/latency"
	"github.com/nalej/grpc-device-controller-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"math/rand"
)

var _ = ginkgo.Describe("Latency Register", func() {

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_device_manager_go.LatencyClient

	// Provider
	var lProvider latency.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		// Create providers
		lProvider = latency.NewMockupProvider()

		manager := NewManager(lProvider)
		handler := NewHandler(manager)
		grpc_device_manager_go.RegisterLatencyServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_device_manager_go.NewLatencyClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.It("should be able to register a ping Latency", func() {
		toAdd := &grpc_device_controller_go.RegisterLatencyRequest{
			OrganizationId: uuid.New().String(),
			DeviceGroupId:  uuid.New().String(),
			DeviceId:       uuid.New().String(),
			Latency:        rand.Int31n(1000) + 1,
		}

		success, err := client.RegisterLatency(context.Background(), toAdd)
		gomega.Expect(err).Should(gomega.Succeed())
		gomega.Expect(success).ShouldNot(gomega.BeNil())
	})

})
