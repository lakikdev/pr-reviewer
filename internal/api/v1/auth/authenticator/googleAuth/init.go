package googleAuth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"pr-reviewer/internal/api/utils"
	"pr-reviewer/internal/api/v1/auth/authenticator"
	"pr-reviewer/internal/api/v1/auth/authenticator/errorAuth"
	"pr-reviewer/internal/database"
	"pr-reviewer/internal/database/dbHelper"
	"pr-reviewer/internal/model"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuth struct {
	name       string
	DB         *database.DB
	authConfig *oauth2.Config
	profileURL string
}

func New(db *database.DB, clientID, clientSecret string) authenticator.IAuthType {
	if clientID == "" && clientSecret == "" {
		return errorAuth.New("the google auth requires a client id and secret")
	}

	authConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{"profile"},
		Endpoint:     google.Endpoint,
	}

	return &GoogleAuth{
		name:       "google",
		DB:         db,
		authConfig: authConfig,
		profileURL: "https://www.googleapis.com/oauth2/v3/userinfo",
	}
}

func (auth *GoogleAuth) Name() string {
	return auth.name
}

func (auth *GoogleAuth) CreateIfNotExists() bool {
	return true
}

func (auth *GoogleAuth) GetRequestData() interface{} {
	return &GoogleAuthTypeRequest{}
}

func (auth *GoogleAuth) GetUserAuth(tx database.TxInterface, requestData interface{}) (*model.UserAuth, error) {
	data := requestData.(*GoogleAuthTypeRequest)

	t, err := auth.authConfig.Exchange(context.TODO(), *data.Code)
	if err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500x101", "Error exchanging code - (server internal error)")
	}

	var response googleResponse
	if err := auth.getProfile(t.AccessToken, &response); err != nil {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500x101", "Error loading google profile - (server internal error)")
	}

	data.ProfileID = &response.ID

	// Check if Google Profile already in use
	userAuth, err := tx.UserAuth().GetByData(auth.Name(), "profileID", data.ProfileID)
	if err != nil && err != dbHelper.ErrItemNotFound {
		return nil, utils.NewHttpError(err, http.StatusInternalServerError, "Err500x102", "Error getting user data - (server internal error)")
	}

	return userAuth, nil
}

func (auth *GoogleAuth) ValidateAuth(userAuth *model.UserAuth, requestData interface{}) error {
	data := requestData.(*GoogleAuthTypeRequest)

	var userAuthData GoogleAuthTypeJSON
	err := json.Unmarshal(userAuth.Data, &userAuthData)
	if err != nil {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500x102", "Error checking auth data - (server internal error)")
	}

	//It should be the same all the time since we selected by ProfileID,
	//Exists as example for other Auth Types and in case we will add additional validation conditions
	if data.ProfileID != nil && userAuthData.ProfileID != nil && *data.ProfileID != *userAuthData.ProfileID {
		return utils.NewHttpError(err, http.StatusInternalServerError, "Err500x102", "Error checking auth data - (server internal error)")
	}

	return nil
}

func (auth *GoogleAuth) BuildAuthData(requestData interface{}) interface{} {
	data := requestData.(*GoogleAuthTypeRequest)
	return &GoogleAuthTypeJSON{
		ProfileID: data.ProfileID,
	}
}

func (auth *GoogleAuth) getProfile(token string, response interface{}) error {
	// Build client:
	client := auth.authConfig.Client(context.TODO(), &oauth2.Token{
		AccessToken: token,
		Expiry:      time.Now().Add(24 * time.Hour),
	})

	// Execute request:
	resp, err := client.Get(auth.profileURL)
	if err != nil {
		return errors.Wrapf(err, "request to provider endpoint %s failed", auth.profileURL)
	}
	defer resp.Body.Close()

	// Deserialize response:
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(response); err != nil {
		return err
	}
	return nil
}

func (auth *GoogleAuth) UpdateUser(user *model.User, requestData interface{}) error {
	return nil
}
