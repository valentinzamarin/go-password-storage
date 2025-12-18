package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewConnection(baseName string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(baseName), &gorm.Config{})
}
