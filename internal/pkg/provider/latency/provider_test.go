/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package latency

import (
	"github.com/google/uuid"
	"github.com/nalej/device-manager/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"time"
)

func RunTest(provider Provider) {
	// AddUser
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

	ginkgo.It("Should be able to get the last latency registry", func(){

		latency := &entities.Latency{
			OrganizationId: uuid.New().String(),
			DeviceGroupId: uuid.New().String(),
			DeviceId:  uuid.New().String(),
			Latency: 300,
			Inserted: time.Now().Unix(),
		}

		err := provider.AddPingLatency(*latency)
		gomega.Expect(err).To(gomega.Succeed())

		latencyLast := &entities.Latency{
			OrganizationId: latency.OrganizationId,
			DeviceGroupId: latency.DeviceGroupId,
			DeviceId:  latency.DeviceId,
			Latency: 200,
			Inserted: time.Now().Unix(),
		}

		err = provider.AddPingLatency(*latencyLast)
		gomega.Expect(err).To(gomega.Succeed())

		retrieved, err := provider.GetLastPingLatency(latency.OrganizationId, latency.DeviceGroupId, latency.DeviceId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved).NotTo(gomega.BeNil())
		gomega.Expect(latencyLast.Latency).Should(gomega.Equal(retrieved.Latency))


	})

	ginkgo.It("Should not be able to get the last latency registry", func(){

		nulTime := int64(0)
		latency := &entities.Latency{
			OrganizationId: uuid.New().String(),
			DeviceGroupId: uuid.New().String(),
			DeviceId:  uuid.New().String(),
			Latency: 300,
			Inserted: time.Now().Unix(),
		}

		retrieved, err := provider.GetLastPingLatency(latency.OrganizationId, latency.DeviceGroupId, latency.DeviceId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved).NotTo(gomega.BeNil())
		gomega.Expect(retrieved.Latency).Should(gomega.Equal(-1))
		gomega.Expect(retrieved.Inserted).Should(gomega.Equal(nulTime))


	})

}