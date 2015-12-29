package pdclient

import "errors"

var (
	//ErrInvalidDispenserResponse - error for invalid statscode response on
	//dispenser call
	ErrInvalidDispenserResponse = errors.New("invalid dispenser response statuscode")

	//ErrInvalidInnKeeperData - error for invalid data object response from innkeeper
	ErrInvalidInnKeeperData = errors.New("invalid innkeeper data object")

	//ProvisionHostInformationFieldname - map key name for provision host info in MetaData
	//of TaskResponse
	ProvisionHostInformationFieldname = "phinfo"
)
