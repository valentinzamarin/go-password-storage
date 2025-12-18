package interfaces

import (
	command "password-storage/internal/app/command"
	"password-storage/internal/app/query"
)

type PasswordService interface {
	AddNewPassword(password *command.AddPasswordCommand) error
	GetAllPasswords() (*query.PasswordsQueryResult, error)
}
