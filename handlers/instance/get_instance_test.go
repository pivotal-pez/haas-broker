package instance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/haas-broker/handlers/instance"
)

var _ = Describe("Get Instance based functionality", func() {
	Describe("given a GetRequestID() func", func() {
		Describe("given a valid instanceid format and collection", func() {
			Context("when a instanceid match is found", func() {
				var (
					collection    *fakeCol
					fakeRequestID = "xyz"
					requestID     string
					err           error
				)
				BeforeEach(func() {
					collection = new(fakeCol)
					collection.FakeResult = []InstanceModel{
						InstanceModel{
							RequestID: fakeRequestID,
						},
					}
					requestID, err = GetRequestID("garbageid", collection)
				})
				It("then it should return a valid taskID", func() {
					Ω(err).ShouldNot(HaveOccurred())
					Ω(requestID).Should(Equal(fakeRequestID))
				})
			})
		})
	})
	Describe("GetTaskID", func() {
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
})
