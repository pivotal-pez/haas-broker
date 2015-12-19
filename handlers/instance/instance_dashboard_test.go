package instance_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-pez/haas-broker/handlers"
	. "github.com/pivotal-pez/haas-broker/handlers/instance"
	"github.com/pivotal-pez/pezdispenser/pdclient/fake"
	"github.com/unrolled/render"
)

var _ = Describe("Management Dashboard", func() {
	Describe("given GetDashboard() function", func() {
		Context("when the instance is not yet provisioned", func() {
			runGetDashboardHandlerContext("running", PendingResponse)
		})
		Context("when the instance is provisioned", func() {
			runGetDashboardHandlerContext("complete", "name")
			runGetDashboardHandlerContext("complete", "oob_ip")
			runGetDashboardHandlerContext("complete", "oob_pw")
			runGetDashboardHandlerContext("complete", "oob_user")
		})
		Context("when the instance has failed provisioning", func() {
			runGetDashboardHandlerContext("failed", FailureResponse)
		})
	})
})

func runGetDashboardHandlerContext(dispenserStatus string, brokerResponseEvaluator string) {
	var (
		responseWriter     *httptest.ResponseRecorder
		controlRequestBody        = "{}"
		dispenserResponse  string = fmt.Sprintf(`{
					"ID": "567471e1c19475001d000001","Timestamp": 1450471905595633562,"Expires": 0,
					"Status": "%s",
					"Profile": "agent_task_long_running","CallerName": "m1.small","MetaData": {
						"status": {
				      "data": {
				        "credentials": {
				          "name": "host-07-16",
				          "oob_ip": "10.65.70.116",
				          "oob_pw": "d3v0ps!",
				          "oob_user": "pezuser"
				        },
				        "status": "complete"
				      },
				      "message": "ok",
				      "status": "success"
						},
						"phinfo": {"data": [
								{
									"requestid": "2676f04e-a5c9-11e5-88f7-0050569b9b57"
								}],
							"message": "ok",
							"status": "success"
						}}}`, dispenserStatus)
	)
	BeforeEach(func() {
		responseWriter = httptest.NewRecorder()
		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewBufferString(controlRequestBody)),
		}
		fakeClientDoer := &fake.ClientDoer{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(dispenserResponse)),
			},
		}
		BuildDashboardHandler(handlers.DispenserCreds{}, new(fakeCol), render.New(), fakeClientDoer)(responseWriter, request)
	})
	It(fmt.Sprintf("then it should return a %s message", brokerResponseEvaluator), func() {
		body, _ := ioutil.ReadAll(responseWriter.Body)
		Ω(body).Should(ContainSubstring(brokerResponseEvaluator))
	})
}
