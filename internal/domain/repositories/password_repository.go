package repositories

import "password-storage/internal/domain/entities"

type PasswordCommands interface {
	Add(password *entities.Password) error
}

type PasswordQueries interface {
	GetAll() ([]*entities.Password, error)
}

type PasswordRepo interface {
	PasswordCommands
	PasswordQueries
}
