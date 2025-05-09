package auth

import (
	"net/http"

	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/database"
	"pr-reviewer/internal/helper/stringext"
	"pr-reviewer/internal/model"
)

func (api *API) issueTokens(tx database.TxInterface, userID model.UserID, sessionData model.SessionData) (*authHelper.Tokens, error) {
	tokens, err := authHelper.IssueToken(model.Principal{UserID: userID})
	if err != nil || tokens == nil {
		return nil, err
	}

	session := model.Session{
		UserID:     userID,
		DeviceID:   sessionData.DeviceID,
		DeviceType: sessionData.DeviceType,
		OSVersion:  sessionData.OSVersion,
		AppVersion: sessionData.AppVersion,
		IPAddress:  sessionData.IPAddress,
		ExpiresAt:  tokens.RefreshTokenExpiresAt,
	}

	if err := session.SetToken(tokens.RefreshToken); err != nil {
		return nil, err
	}

	if err := tx.Session().Save(session); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (api *API) setIPAddress(r *http.Request, sessionData *model.SessionData) error {
	xForwardedFor := r.Header["X-Forwarded-For"]
	if len(xForwardedFor) == 0 {
		sessionData.IPAddress = stringext.New("unknown")
		return nil
	}

	ipAddress := xForwardedFor[0]
	sessionData.IPAddress = &ipAddress
	return nil
}
