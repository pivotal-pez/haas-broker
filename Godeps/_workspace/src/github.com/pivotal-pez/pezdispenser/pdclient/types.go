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
	LeaseCreateResponseBody struct {
		ID         string                 `json:"id"`
		Timestamp  int64                  `json:"timestamp"`
		Expires    int64                  `json:"expires"`
		Status     string                 `json:"status"`
		Profile    string                 `json:"profile"`
		CallerName string                 `json:"caller_name"`
		MetaData   map[string]interface{} `json:"meta_data"`
	}
)
