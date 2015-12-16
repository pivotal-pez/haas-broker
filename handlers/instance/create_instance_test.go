package instance_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-pez/haas-broker/handlers/instance"
	"github.com/pivotal-pez/pezdispenser/pdclient/fake"
)

var _ = Describe("InstanceCreator", func() {

	Describe("given PutHandler() method", func() {
		Context("when called with a new service instance request", func() {
			var (
				instanceCreator    *InstanceCreator
				responseWriter     *httptest.ResponseRecorder
				dispenserResponse  string = `{"id": "560ede8bfccecc0072000001"}`
				controlRequestBody string = `{
					"organization_guid": "org-guid-here",
					"plan_id":           "plan-guid-here",
					"service_id":        "service-guid-here",
					"space_guid":        "space-guid-here",
					"parameters":        {
						"parameter1": 1,
						"parameter2": "value"
					}
				}`
			)

			BeforeEach(func() {
				instanceCreator = new(InstanceCreator)
				instanceCreator.Collection = new(fakeCol)
				responseWriter = httptest.NewRecorder()
				request := &http.Request{
					Body: ioutil.NopCloser(bytes.NewBufferString(controlRequestBody)),
				}
				HttpClient = &fake.ClientDoer{
					Response: &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewBufferString(dispenserResponse)),
					},
				}
				instanceCreator.PutHandler(responseWriter, request)
			})
			It("then it should return a proper statuscode for async calls", func() {
				Ω(responseWriter.Code).Should(Equal(http.StatusAccepted))
			})

			It("then it should return a response body containing the dashboardurl", func() {
				body, _ := ioutil.ReadAll(responseWriter.Body)
				Ω(body).Should(ContainSubstring(DashboardUrl))
			})

			It("then it should properly parse the request body", func() {
				Ω(instanceCreator.Model.OrganizationGUID).Should(Equal("org-guid-here"))
				Ω(instanceCreator.Model.PlanID).Should(Equal("plan-guid-here"))
				Ω(instanceCreator.Model.ServiceID).Should(Equal("service-guid-here"))
				Ω(instanceCreator.Model.SpaceGUID).Should(Equal("space-guid-here"))
				Ω(instanceCreator.Model.TaskGUID).Should(Equal("560ede8bfccecc0072000001"))
			})

			Context("and request has an invalid body", func() {
				BeforeEach(func() {
					responseWriter = httptest.NewRecorder()
					brokenJson := `{?::"parameters":1234}`
					request := &http.Request{
						Body: ioutil.NopCloser(bytes.NewBufferString(brokenJson)),
					}
					instanceCreator.PutHandler(responseWriter, request)
				})
				It("then it should return an error response", func() {
					Ω(responseWriter.Code).Should(Equal(http.StatusNotAcceptable))
				})

				It("then it should return a error status", func() {
					body, _ := ioutil.ReadAll(responseWriter.Body)
					Ω(body).Should(ContainSubstring("error_message"))
				})
			})
		})
	})
})
