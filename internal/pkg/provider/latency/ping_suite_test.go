/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */
package latency

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestPingProviderPackage(t *testing.T){
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "PingProviders package suite")
}