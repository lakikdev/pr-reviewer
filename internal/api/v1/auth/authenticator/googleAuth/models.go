package googleAuth

import "errors"

type GoogleAuthTypeJSON struct {
	ProfileID *string `json:"profileID"`
}

type GoogleAuthTypeRequest struct {
	Code      *string `json:"code"`
	ProfileID *string `json:"profileID"`
}

func (data *GoogleAuthTypeRequest) Verify() error {
	if data.Code == nil || len(*data.Code) == 0 {
		return errors.New("code is required")
	}
	return nil
}

type googleResponse struct {
	ID    string `json:"sub"`
	Email string `json:"email"`
}
