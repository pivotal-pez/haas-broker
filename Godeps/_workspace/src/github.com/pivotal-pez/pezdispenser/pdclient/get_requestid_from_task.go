package pdclient

import (
	"encoding/json"

	"github.com/xchapter7x/lo"
)

func GetRequestIDFromTaskResponse(taskResponse TaskResponse) (requestID string, err error) {
	var provisionHostInfoBytes []byte
	firstRecordIndex := 0
	meta := taskResponse.MetaData
	provisionHostInfo := ProvisionHostInfo{}
	lo.G.Debug("taskResponse: ", taskResponse)
	lo.G.Debug("metadata: ", meta)

	if provisionHostInfoBytes, err = json.Marshal(meta[ProvisionHostInformationFieldname]); err == nil {

		if err = json.Unmarshal(provisionHostInfoBytes, &provisionHostInfo); err == nil {

			if len(provisionHostInfo.Data) > firstRecordIndex {
				requestID = provisionHostInfo.Data[firstRecordIndex].RequestID

			} else {
				lo.G.Error("no request id found in: ", provisionHostInfo)
			}

		} else {
			lo.G.Error("error unmarshalling: ", err, meta)
			lo.G.Error("metadata: ", meta)
		}

	} else {
		lo.G.Error("error marshalling: ", err)
	}
	return
}
