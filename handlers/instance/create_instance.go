package instance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//PutHandler - this is the actual handler method that will be used for the
//incoming request
func (s *InstanceCreator) PutHandler(w http.ResponseWriter, req *http.Request) {
	var (
		err          error
		bodyBytes    []byte
		statusCode   int
		responseBody string
	)

	if bodyBytes, err = ioutil.ReadAll(req.Body); err == nil {

		if err = json.Unmarshal(bodyBytes, &s.Model); err == nil {
			statusCode = http.StatusAccepted
			responseBody = fmt.Sprintf(`{"dashboard_url": "%s"}`, DashboardUrl)
		}
	}

	if err != nil {
		statusCode = http.StatusNotAcceptable
		responseBody = fmt.Sprintf(`{"error_message": "%s"}`, err.Error())
	}
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, responseBody)
}
