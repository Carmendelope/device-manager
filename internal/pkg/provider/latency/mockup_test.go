package latency

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup latency provider", func(){


	sp := NewMockupProvider()
	RunTest(sp)

})

