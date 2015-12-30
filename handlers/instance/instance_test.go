package instance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-pez/haas-broker/handlers"
	. "github.com/pivotal-pez/haas-broker/handlers/instance"
)

var _ = Describe("Instance", func() {
	Describe("given NewInstanceCreator() function", func() {
		Context("when passed a valid collection and dispenserCreds objects ", func() {
			var collection *fakeCol
			var dispenserCreds handlers.DispenserCreds
			var instanceCreator *InstanceCreator

			BeforeEach(func() {
				collection = new(fakeCol)
				dispenserCreds = handlers.DispenserCreds{
					ApiKey: "fake-api-key",
					URL:    "fake-url.com",
				}
				instanceCreator = NewInstanceCreator(collection, dispenserCreds)
			})
			It("then it should properly initialize a instancecreator and return it.", func() {
				Ω(instanceCreator.Dispenser).Should(Equal(dispenserCreds))
				Ω(instanceCreator.Collection).Should(Equal(collection))
			})
		})
	})
})
