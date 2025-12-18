package interfaces

import (
	"errors"
	"password-storage/internal/app/command"
	"password-storage/internal/app/query"
)

type AuthService interface {
	IsMasterPasswordSet(q *query.IsMasterPasswordSetQuery) (*query.IsMasterPasswordSetQueryResult, error)
	CreateMasterPassword(cmd *command.CreateMasterPasswordCommand) error
	Authenticate(cmd *command.AuthenticateCommand) error
}

var ErrNotFound = errors.New("record not found")
