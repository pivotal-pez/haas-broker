package pdclient

import (
	"encoding/json"

	"github.com/xchapter7x/lo"
)

func GetRequestIDFromTaskResponse(taskResponse TaskResponse) (requestID string) {
	firstRecordIndex := 0
	meta := taskResponse.MetaData
	provisionHostInfo := ProvisionHostInfo{}

	if provisionHostInfoBytes, err := json.Marshal(meta[ProvisionHostInformationFieldname]); err == nil {

		if err = json.Unmarshal(provisionHostInfoBytes, &provisionHostInfo); err == nil {
			requestID = provisionHostInfo.Data[firstRecordIndex].RequestID

		} else {
			lo.G.Error("error unmarshalling: ", err, meta)
			lo.G.Error("metadata: ", meta)
		}

	} else {
		lo.G.Error("error marshalling: ", err)
	}
	return
}