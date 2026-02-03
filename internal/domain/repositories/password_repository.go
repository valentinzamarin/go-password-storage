package repositories

import (
	"context"
	"password-storage/internal/domain/entities"
)

type PasswordCommands interface {
	Add(ctx context.Context, password *entities.Password) error
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, password *entities.Password) error
}

type PasswordQueries interface {
	GetAll(ctx context.Context) ([]*entities.Password, error)
	GetByID(ctx context.Context, id uint) (*entities.Password, error)
}

type PasswordRepo interface {
	PasswordCommands
	PasswordQueries
}
