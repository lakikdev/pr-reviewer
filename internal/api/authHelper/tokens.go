package authHelper

import (
	"errors"
	"time"

	"pr-reviewer/internal/config"
	"pr-reviewer/internal/model"

	jwt "github.com/dgrijalva/jwt-go"
)

var accessTokenDuration = time.Duration(30) * time.Minute   // 30 min
var refreshTokenDuration = time.Duration(30*24) * time.Hour //30 days

type Claims struct {
	UserID model.UserID `json:"userID"`
	Type   *string      `json:"type,omitempty"`
	jwt.StandardClaims
}

// Tokens is wrapper for access and refresh tokens
type Tokens struct {
	AccessToken           string `json:"accessToken,omitempty"`
	AccessTokenExpiresAt  int64  `json:"accessTokenExpiresAt,omitempty"` //we return only access token's expires at time
	RefreshToken          string `json:"refreshToken,omitempty"`
	RefreshTokenExpiresAt int64  `json:"refreshTokenExpiresAt,omitempty"` //we will store this time in database with refresh token
}

// IssueToken generate Access and Refresh tokens
func IssueToken(principal model.Principal) (*Tokens, error) {
	if principal.UserID == model.NilUserID {
		return nil, errors.New("invalid principal")
	}

	accessToken, accessTokenExpiresAt, err := generateToken(principal, accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenExpiresAt, err := generateToken(principal, refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	tokens := Tokens{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}

	return &tokens, nil
}

func generateToken(principal model.Principal, duration time.Duration) (string, int64, error) {
	now := time.Now()
	//Generate Access Token
	claims := &Claims{
		UserID: principal.UserID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.UTC().Unix(),
			ExpiresAt: now.UTC().Add(duration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(*config.JWTSecretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, claims.ExpiresAt, nil
}

func VerifyToken(token string) (*model.Principal, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(*config.JWTSecretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}

	principal := &model.Principal{
		UserID: claims.UserID,
	}

	//we want to return principal even if token invalid because we need to get userID
	if !tkn.Valid {
		return principal, err
	}

	return principal, nil

}
