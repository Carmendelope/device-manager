/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package device

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestDevicePackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Device package suite")
}