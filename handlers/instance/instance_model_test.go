package instance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/haas-broker/handlers/instance"
)

var _ = Describe("InstanceModel", func() {
	var (
		instanceModel  InstanceModel
		fakeCollection *fakeCol
	)
	Describe("given Save() method", func() {
		Context("when given a valid collection", func() {
			BeforeEach(func() {
				fakeCollection = new(fakeCol)
				instanceModel = InstanceModel{
					TaskGUID:         "asdf",
					OrganizationGUID: "123456",
					PlanID:           "ahbaobe",
					ServiceID:        "iiiibdodb",
					SpaceGUID:        "spacey",
				}
				instanceModel.Save(fakeCollection)
			})
			It("then it should save the model to the given collection", func() {
				Ω(len(fakeCollection.SpyID.(string))).Should(Equal(len("5671bda239df83c73fff9068")))
				Ω(fakeCollection.SpyUpdate).Should(Equal(instanceModel))
			})
		})
	})
})
