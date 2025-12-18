package passwords

import (
	"fmt"
	"password-storage/internal/domain/entities"
	"password-storage/internal/domain/repositories"

	"gorm.io/gorm"
)

type GormPasswordRepository struct {
	db *gorm.DB
}

func NewGormPasswordRepository(db *gorm.DB) repositories.PasswordRepo {
	return &GormPasswordRepository{
		db: db,
	}
}

func (p *GormPasswordRepository) Add(password *entities.Password) error {
	dbPassword := toDBPassword(*password)
	err := p.db.Create(dbPassword).Error
	if err != nil {
		fmt.Printf("DB save failed: %v\n", err)
		return err
	}
	return nil
}

func (p *GormPasswordRepository) GetAll() ([]*entities.Password, error) {
	var models []*PasswordModel
	if err := p.db.Find(&models).Error; err != nil {
		return nil, err
	}
	passwords := make([]*entities.Password, 0, len(models))

	for _, model := range models {
		password := fromDBPassword(model)
		passwords = append(passwords, password)
	}

	return passwords, nil
}
