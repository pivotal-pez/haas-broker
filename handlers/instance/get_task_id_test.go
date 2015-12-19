package instance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/haas-broker/handlers/instance"
)

var _ = Describe("GetTaskID", func() {
	Describe("given a valid instanceid format and collection", func() {
		Context("when a instanceid match is found", func() {
			var (
				collection   *fakeCol
				fakeTaskGUID = "xyz"
				taskID       string
				err          error
			)
			BeforeEach(func() {
				collection = new(fakeCol)
				collection.FakeResult = []InstanceModel{
					InstanceModel{
						TaskGUID: fakeTaskGUID,
					},
				}
				taskID, err = GetTaskID("garbageid", collection)
			})
			It("then it should return a valid taskID", func() {
				Ω(err).ShouldNot(HaveOccurred())
				Ω(taskID).Should(Equal(fakeTaskGUID))
			})
		})
	})
})
