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
	"github.com/google/uuid"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"math/rand"
	"time"
)

func RunTest(provider Provider) {
	ginkgo.It("Should be able to add a latency registry", func(){

		latency := &entities.Latency{
			OrganizationId: uuid.New().String(),
			DeviceGroupId: uuid.New().String(),
			DeviceId:  uuid.New().String(),
			Latency: 300,
			Inserted: time.Now().Unix(),
		}

		err := provider.AddPingLatency(*latency)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should be able to get the latencies of a device", func(){

		latency := &entities.Latency{
			OrganizationId: uuid.New().String(),
			DeviceGroupId: uuid.New().String(),
			DeviceId:  uuid.New().String(),
			Latency: 300,
			Inserted: time.Now().Unix(),
		}

		err := provider.AddPingLatency(*latency)
		gomega.Expect(err).To(gomega.Succeed())

		latency2 := &entities.Latency{
			OrganizationId: latency.OrganizationId,
			DeviceGroupId: latency.DeviceGroupId,
			DeviceId:  latency.DeviceId,
			Latency: 300,
			Inserted: time.Now().Unix()+4,
		}

		err = provider.AddPingLatency(*latency2)
		gomega.Expect(err).To(gomega.Succeed())

		retrieved, err := provider.GetLatency(latency.OrganizationId, latency.DeviceGroupId, latency.DeviceId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved).NotTo(gomega.BeNil())
		gomega.Expect(len(retrieved)).Should(gomega.Equal(2))

	})
	ginkgo.It("Should be able to remove a latency", func() {
		latency := &entities.Latency{
			OrganizationId: uuid.New().String(),
			DeviceGroupId: uuid.New().String(),
			DeviceId:  uuid.New().String(),
			Latency: 300,
			Inserted: time.Now().Unix(),
		}

		err := provider.AddPingLatency(*latency)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.RemoveLatency(latency.OrganizationId, latency.DeviceGroupId, latency.DeviceId)
		gomega.Expect(err).To(gomega.Succeed())
	})

	// ------------------------------
	ginkgo.It("Should be able to add a last latency registry", func(){

		latency := &entities.Latency{
			OrganizationId: uuid.New().String(),
			DeviceGroupId: uuid.New().String(),
			DeviceId:  uuid.New().String(),
			Latency: 300,
			Inserted: time.Now().Unix(),
		}

		err := provider.AddLastLatency(*latency)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should be able to get the last latency registry", func(){

		latency := &entities.Latency{
			OrganizationId: uuid.New().String(),
			DeviceGroupId: uuid.New().String(),
			DeviceId:  uuid.New().String(),
			Latency: 300,
			Inserted: time.Now().Unix(),
		}

		err := provider.AddLastLatency(*latency)
		gomega.Expect(err).To(gomega.Succeed())

		latencyLast := &entities.Latency{
			OrganizationId: latency.OrganizationId,
			DeviceGroupId: latency.DeviceGroupId,
			DeviceId:  latency.DeviceId,
			Latency: 200,
			Inserted: time.Now().Unix(),
		}

		err = provider.AddLastLatency(*latencyLast)
		gomega.Expect(err).To(gomega.Succeed())

		retrieved, err := provider.GetLastLatency(latency.OrganizationId, latency.DeviceGroupId, latency.DeviceId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved).NotTo(gomega.BeNil())
		gomega.Expect(latencyLast.Latency).Should(gomega.Equal(retrieved.Latency))


	})
	ginkgo.It("Should be able to get the latency list of a group", func(){

		organizationID := uuid.New().String()
		DeviceGroupId := uuid.New().String()
		numLatencies := 5
		for i:=0 ; i< numLatencies; i++ {
			latency := &entities.Latency{
				OrganizationId: organizationID,
				DeviceGroupId:  DeviceGroupId,
				DeviceId:       uuid.New().String(),
				Latency:        rand.Intn(500) +1,
				Inserted:       time.Now().Unix()+int64(i),
			}

			err := provider.AddLastLatency(*latency)
			gomega.Expect(err).To(gomega.Succeed())
		}

		list, err := provider.GetGroupLastLatencies(organizationID, DeviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(list).NotTo(gomega.BeNil())
		gomega.Expect(len(list)).Should(gomega.Equal(numLatencies))

	})
	ginkgo.It("Should be able to get an empty latency list of a non existing group", func(){

		organizationID := uuid.New().String()
		DeviceGroupId := uuid.New().String()

		list, err := provider.GetGroupLastLatencies(organizationID, DeviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(list).NotTo(gomega.BeNil())
		gomega.Expect(list).To(gomega.BeEmpty())

	})

}