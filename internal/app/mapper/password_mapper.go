package mapper

import (
	"password-storage/internal/app/common"
	"password-storage/internal/domain/entities"
)

func ToPasswordResult(p *entities.Password, decryptedPass string) *common.PasswordsResult {
	return &common.PasswordsResult{
		ID:          p.GetID(),
		URL:         p.GetURL(),
		Login:       p.GetLogin(),
		Password:    decryptedPass,
		Description: p.GetDescription(),
	}
}
