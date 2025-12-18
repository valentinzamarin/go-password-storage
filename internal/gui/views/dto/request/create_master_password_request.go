package request

import (
	"errors"
	"password-storage/internal/app/command"
)

const MinMasterPasswordLength = 2

var (
	ErrMasterPasswordEmpty    = errors.New("password cannot be empty")
	ErrMasterPasswordTooShort = errors.New("password must be at least 2 characters")
	ErrPasswordsDoNotMatch    = errors.New("passwords do not match")
)

type CreateMasterPasswordRequest struct {
	Password        string
	ConfirmPassword string
}

func (req *CreateMasterPasswordRequest) Validate() error {
	if req.Password == "" {
		return ErrMasterPasswordEmpty
	}

	if len(req.Password) < MinMasterPasswordLength {
		return ErrMasterPasswordTooShort
	}

	if req.Password != req.ConfirmPassword {
		return ErrPasswordsDoNotMatch
	}

	return nil
}

func (req *CreateMasterPasswordRequest) ToCreateMasterPasswordCommand() (*command.CreateMasterPasswordCommand, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &command.CreateMasterPasswordCommand{
		Password: req.Password,
	}, nil
}
