package pdclient

import "errors"

var (
	//ErrInvalidDispenserResponse - error for invalid statscode response on
	//dispenser call
	ErrInvalidDispenserResponse = errors.New("invalid dispenser response statuscode")

	//ProvisionHostInformationFieldname - map key name for provision host info in MetaData
	//of TaskResponse
	ProvisionHostInformationFieldname = "phinfo"
)
