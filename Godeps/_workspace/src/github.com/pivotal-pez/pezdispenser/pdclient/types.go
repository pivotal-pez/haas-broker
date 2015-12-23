package pdclient

import "net/http"

type (
	//PDClient - dispenser client object
	PDClient struct {
		APIKey string
		client clientDoer
		URL    string
	}
	clientDoer interface {
		Do(req *http.Request) (resp *http.Response, err error)
	}
	//LeaseRequestBody - request lease body object structure
	LeaseRequestBody struct {
		LeaseID              string                 `json:"lease_id"`
		InventoryID          string                 `json:"inventory_id"`
		Username             string                 `json:"username"`
		Sku                  string                 `json:"sku"`
		LeaseDuration        int64                  `json:"lease_duration"`
		LeaseEndDate         int64                  `json:"lease_end_date"`
		LeaseStartDate       int64                  `json:"lease_start_date"`
		LeaseProcurementMeta map[string]interface{} `json:"procurement_meta"`
	}
	//TaskResponse - a response object for a get task call
	TaskResponse struct {
		ID         string                 `json:"ID"`
		Timestamp  int64                  `json:"Timestamp"`
		Expires    int64                  `json:"Expires"`
		Status     string                 `json:"Status"`
		Profile    string                 `json:"Profile"`
		CallerName string                 `json:"CallerName"`
		MetaData   map[string]interface{} `json:"MetaData"`
	}
	ProvisionHostInfo struct {
		Data []ProvisionHostData `json:"data"`
	}
	ProvisionHostData struct {
		RequestID string `json:"requestid"`
	}
)
