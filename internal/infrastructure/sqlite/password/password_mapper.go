package passwords

import "password-storage/internal/domain/entities"

func toDBPassword(password entities.Password) *PasswordModel {
	return &PasswordModel{
		ID:                password.GetID(),
		URL:               password.GetURL(),
		Login:             password.GetLogin(),
		EncryptedPassword: []byte(password.GetPassword()),
		Description:       password.GetDescription(),
	}
}

func fromDBPassword(model *PasswordModel) *entities.Password {
	pw, _ := entities.NewPassword(
		model.URL,
		model.Login,
		string(model.EncryptedPassword),
		model.Description,
	)
	pw.SetID(model.ID)
	return pw
}
