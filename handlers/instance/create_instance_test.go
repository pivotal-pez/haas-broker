package instance_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-pez/cfmgo"
	. "github.com/pivotal-pez/haas-broker/handlers/instance"
	"github.com/pivotal-pez/pezdispenser/pdclient/fake"
)

var _ = Describe("InstanceCreator", func() {
	Describe("given GetHandler() method", func() {
		Context("when the instance is not yet provisioned", func() {
			runGetHandlerContext("running", "in progress")
		})
		Context("when the instance is provisioned", func() {
			runGetHandlerContext("complete", "succeeded")
		})
		Context("when the instance has failed provisioning", func() {
			runGetHandlerContext("failed", "failed")
		})
	})
	Describe("given PutHandler() method", func() {
		Context("when called with a new service instance request", func() {
			var (
				instanceCreator    *InstanceCreator
				responseWriter     *httptest.ResponseRecorder
				controlRequestID          = "2676f04e-a5c9-11e5-88f7-0050569b9b57"
				controlTaskID             = "560ede8bfccecc0072000001"
				controlStatus             = "complete"
				dispenserResponse         = makeDispenserResponse(controlTaskID, controlRequestID, controlStatus)
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
				instanceCreator.ClientDoer = &fake.ClientDoer{
					Response: &http.Response{
						StatusCode: http.StatusCreated,
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
				Ω(body).Should(ContainSubstring("dashboard_url"))
			})

			It("then it should properly parse the request body", func() {
				Ω(instanceCreator.Model.OrganizationGUID).Should(Equal("org-guid-here"))
				Ω(instanceCreator.Model.PlanID).Should(Equal("plan-guid-here"))
				Ω(instanceCreator.Model.ServiceID).Should(Equal("service-guid-here"))
				Ω(instanceCreator.Model.SpaceGUID).Should(Equal("space-guid-here"))
				Ω(instanceCreator.Model.TaskGUID).Should(Equal(controlTaskID))
				Ω(instanceCreator.Model.RequestID).Should(Equal(controlRequestID))
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

func runGetHandlerContext(dispenserStatus string, brokerResponseEvaluator string) {
	var (
		instanceCreator    *InstanceCreator
		responseWriter     *httptest.ResponseRecorder
		controlRequestBody = "{}"
		controlRequestID   = "2676f04e-a5c9-11e5-88f7-0050569b9b57"
		controlTaskID      = "560ede8bfccecc0072000001"
		dispenserResponse  = makeDispenserResponse(controlTaskID, controlRequestID, dispenserStatus)
	)
	var origGetTaskID func(instanceID string, collection cfmgo.Collection) (taskID string, err error)
	AfterEach(func() {
		GetTaskID = origGetTaskID
	})
	BeforeEach(func() {
		origGetTaskID = GetTaskID
		GetTaskID = func(instanceID string, collection cfmgo.Collection) (taskID string, err error) {
			return "567471e1c19475001d000001", nil
		}
		instanceCreator = new(InstanceCreator)
		instanceCreator.Collection = new(fakeCol)
		responseWriter = httptest.NewRecorder()
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewBufferString(controlRequestBody)),
		}
		instanceCreator.ClientDoer = &fake.ClientDoer{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(dispenserResponse)),
			},
		}
		instanceCreator.GetHandler(responseWriter, request)
	})
	It(fmt.Sprintf("then it should return a %s message", brokerResponseEvaluator), func() {
		body, _ := ioutil.ReadAll(responseWriter.Body)
		Ω(body).Should(ContainSubstring(brokerResponseEvaluator))
	})
}

func makeDispenserResponse(id string, requestID string, status string) string {
	return fmt.Sprintf(`{
		"ID": "%s",
		"Timestamp": 1450479310174826830,
		"Expires": 0,
		"Status": "%s",
		"Profile": "agent_task_long_running",
		"CallerName": "m1.small",
		"MetaData": {
			"phinfo": {
				"data": [
					{
						"requestid": "%s"
					}
				],
				"message": "ok",
				"status": "success"
			}
		}
	}`, id, status, requestID)
}
