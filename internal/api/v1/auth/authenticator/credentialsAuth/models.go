package credentialsAuth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type CredentialsAuthTypeJSON struct {
	Email        *string `json:"email"`
	PasswordHash *[]byte `json:"password"`
}

// Set Password updates a user's password
func (cat *CredentialsAuthTypeJSON) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	cat.PasswordHash = &hash
	return nil
}

func (cat *CredentialsAuthTypeJSON) CheckPassword(password string) error {
	if cat.PasswordHash == nil || len(*cat.PasswordHash) == 0 {
		return errors.New("password not set")
	}
	return bcrypt.CompareHashAndPassword(*cat.PasswordHash, []byte(password))
}

type CredentialsAuthTypeRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (data *CredentialsAuthTypeRequest) Verify() error {
	if data.Email == nil || len(*data.Email) == 0 {
		return errors.New("email is required")
	}

	if data.Password == nil || len(*data.Password) == 0 {
		return errors.New("password is required")
	}
	return nil
}
