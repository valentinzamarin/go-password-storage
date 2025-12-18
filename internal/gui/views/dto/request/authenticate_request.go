package request

import (
	"errors"
	"password-storage/internal/app/command"
)

var ErrAuthPasswordEmpty = errors.New("password cannot be empty")

type AuthenticateRequest struct {
	Password string
}

func (req *AuthenticateRequest) Validate() error {
	if req.Password == "" {
		return ErrAuthPasswordEmpty
	}
	return nil
}

func (req *AuthenticateRequest) ToCommand() (*command.AuthenticateCommand, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &command.AuthenticateCommand{
		Password: req.Password,
	}, nil
}
