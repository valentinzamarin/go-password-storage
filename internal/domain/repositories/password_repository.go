package repositories

import "password-storage/internal/domain/entities"

type PasswordCommands interface {
	Add(password *entities.Password) error
	Delete(id uint) error
}

type PasswordQueries interface {
	GetAll() ([]*entities.Password, error)
}

type PasswordRepo interface {
	PasswordCommands
	PasswordQueries
}
