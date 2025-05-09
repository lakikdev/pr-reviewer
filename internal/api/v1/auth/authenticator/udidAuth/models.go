package udidAuth

import "errors"

type UDIDAuthTypeJSON struct {
	UDID *string `json:"udid"`
}

type UDIDAuthTypeRequest struct {
	UDID *string `json:"udid"`
}

func (data *UDIDAuthTypeRequest) Verify() error {
	if data.UDID == nil || len(*data.UDID) == 0 {
		return errors.New("udid is required")
	}
	return nil
}
