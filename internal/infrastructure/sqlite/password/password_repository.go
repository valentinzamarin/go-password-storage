package passwords

import (
	"context"
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

func (p *GormPasswordRepository) Add(ctx context.Context, password *entities.Password) error {
	dbPassword := toDBPassword(*password)
	err := p.db.WithContext(ctx).Create(dbPassword).Error
	if err != nil {
		fmt.Printf("DB save failed: %v\n", err)
		return err
	}
	return nil
}
func (p *GormPasswordRepository) GetAll(ctx context.Context) ([]*entities.Password, error) {
	var models []*PasswordModel
	if err := p.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	passwords := make([]*entities.Password, 0, len(models))

	for _, model := range models {
		password := fromDBPassword(model)
		passwords = append(passwords, password)
	}

	return passwords, nil
}

func (p *GormPasswordRepository) Delete(ctx context.Context, id uint) error {
	if err := p.db.WithContext(ctx).Delete(&PasswordModel{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (p *GormPasswordRepository) Update(ctx context.Context, password *entities.Password) error {
	dbPassword := toDBPassword(*password)
	if err := p.db.WithContext(ctx).Save(dbPassword).Error; err != nil {
		fmt.Printf("DB update failed: %v\n", err)
		return err
	}
	return nil
}

func (p *GormPasswordRepository) GetByID(ctx context.Context, id uint) (*entities.Password, error) {
	var model PasswordModel
	if err := p.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return fromDBPassword(&model), nil
}
