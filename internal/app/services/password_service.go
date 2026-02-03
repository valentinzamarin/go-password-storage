package services

import (
	"context"
	"fmt"
	command "password-storage/internal/app/command"
	"password-storage/internal/app/common"
	"password-storage/internal/app/interfaces"
	"password-storage/internal/app/mapper"
	"password-storage/internal/app/query"
	"password-storage/internal/domain/entities"
	"password-storage/internal/domain/repositories"
)

type PasswordService struct {
	passwordRepo repositories.PasswordRepo
	encrypt      interfaces.PasswordEncryptor
}

func NewPasswordService(passwordRepo repositories.PasswordRepo, encrypt interfaces.PasswordEncryptor) interfaces.PasswordService {
	return &PasswordService{
		passwordRepo: passwordRepo,
		encrypt:      encrypt,
	}
}

func (ps *PasswordService) AddNewPassword(ctx context.Context, addCommand *command.AddPasswordCommand) error {

	encryptedPass, err := ps.encrypt.Encrypt([]byte(addCommand.Password))
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	newPassword, err := entities.NewPassword(
		addCommand.URL,
		addCommand.Login,
		string(encryptedPass),
		addCommand.Description,
	)
	if err != nil {
		return err
	}

	if err := ps.passwordRepo.Add(ctx, newPassword); err != nil {
		return fmt.Errorf("failed to save password: %w", err)
	}

	return nil
}

func (ps *PasswordService) GetAllPasswords(ctx context.Context) (*query.PasswordsQueryResult, error) {
	allPasswords, err := ps.passwordRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get passwords: %w", err)
	}

	passwordsToView := make([]*common.PasswordsResult, 0, len(allPasswords))

	for _, p := range allPasswords {
		decryptedPass, err := ps.encrypt.Decrypt([]byte(p.GetPassword()))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt password: %w", err)
		}

		passwordsToView = append(passwordsToView, mapper.ToPasswordResult(p, string(decryptedPass)))
	}

	return &query.PasswordsQueryResult{
		Result: passwordsToView,
	}, nil
}

func (ps *PasswordService) DeletePassword(ctx context.Context, cmd *command.DeletePasswordCommand) error {
	if cmd == nil {
		return fmt.Errorf("delete command is nil")
	}
	if err := ps.passwordRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete password: %w", err)
	}
	return nil
}

func (ps *PasswordService) UpdatePassword(ctx context.Context, cmd *command.UpdatePasswordCommand) error {
	if cmd == nil {
		return fmt.Errorf("update command is nil")
	}
	// load existing entity
	existing, err := ps.passwordRepo.GetByID(ctx, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to load password: %w", err)
	}

	encryptedPass, err := ps.encrypt.Encrypt([]byte(cmd.Password))
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	if err := existing.Update(cmd.URL, cmd.Login, string(encryptedPass), cmd.Description); err != nil {
		return fmt.Errorf("invalid update: %w", err)
	}

	if err := ps.passwordRepo.Update(ctx, existing); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
