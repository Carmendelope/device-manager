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
}