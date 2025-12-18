package auth

import (
	"errors"
	"password-storage/internal/app/interfaces"

	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) interfaces.AuthRepository {
	return &AuthRepo{
		db: db,
	}
}

func (r *AuthRepo) GetMasterAuth() ([]byte, []byte, error) {
	var masterAuth AuthModel

	if err := r.db.First(&masterAuth).Error; err != nil {
		return nil, nil, err
	}

	return masterAuth.Salt, masterAuth.VerificationHash, nil
}

func (r *AuthRepo) CreateMasterAuth(salt []byte, verificationHash []byte) error {

	var count int64
	if err := r.db.Model(&AuthModel{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("master password already exists")
	}

	masterAuth := &AuthModel{
		Salt:             salt,
		VerificationHash: verificationHash,
	}

	return r.db.Create(masterAuth).Error
}
