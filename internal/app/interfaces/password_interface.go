package interfaces

import (
	"context"
	command "password-storage/internal/app/command"
	"password-storage/internal/app/query"
)

type PasswordService interface {
	AddNewPassword(ctx context.Context, password *command.AddPasswordCommand) error
	GetAllPasswords(ctx context.Context) (*query.PasswordsQueryResult, error)
	DeletePassword(ctx context.Context, cmd *command.DeletePasswordCommand) error
	UpdatePassword(ctx context.Context, cmd *command.UpdatePasswordCommand) error
}
