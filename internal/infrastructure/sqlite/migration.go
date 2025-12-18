package sqlite

import (
	"password-storage/internal/infrastructure/sqlite/auth"
	passwords "password-storage/internal/infrastructure/sqlite/password"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&passwords.PasswordModel{},
		&auth.AuthModel{},
	)
}
