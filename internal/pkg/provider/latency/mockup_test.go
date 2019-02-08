package latency

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup role provider", func(){


	sp := NewMockupProvider()
	RunTest(sp)

})

