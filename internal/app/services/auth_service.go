package services

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"password-storage/internal/app/command"
	"password-storage/internal/app/interfaces"
	"password-storage/internal/app/query"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
)

type AuthService struct {
	authRepo  interfaces.AuthRepository
	encryptor interfaces.PasswordEncryptor
}

func NewAuthService(authRepo interfaces.AuthRepository, encryptor interfaces.PasswordEncryptor) interfaces.AuthService {
	return &AuthService{
		authRepo:  authRepo,
		encryptor: encryptor,
	}
}

func (s *AuthService) hashPassword(password string, salt []byte) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)
}

func (s *AuthService) IsMasterPasswordSet(q *query.IsMasterPasswordSetQuery) (*query.IsMasterPasswordSetQueryResult, error) {
	_, _, err := s.authRepo.GetMasterAuth()
	if err != nil {
		if errors.Is(err, interfaces.ErrNotFound) {
			return &query.IsMasterPasswordSetQueryResult{IsSet: false}, nil
		}
		return nil, err
	}
	return &query.IsMasterPasswordSetQueryResult{IsSet: true}, nil
}

func (s *AuthService) CreateMasterPassword(cmd *command.CreateMasterPasswordCommand) error {
	salt, err := s.encryptor.GenerateSalt()
	if err != nil {
		return fmt.Errorf("could not generate salt: %w", err)
	}

	verificationHash := s.hashPassword(cmd.Password, salt)

	if err := s.authRepo.CreateMasterAuth(salt, verificationHash); err != nil {
		return fmt.Errorf("could not save master auth data: %w", err)
	}

	s.encryptor.DeriveKeyFromPassword(cmd.Password, salt)
	return nil
}

func (s *AuthService) Authenticate(cmd *command.AuthenticateCommand) error {
	salt, verificationHash, err := s.authRepo.GetMasterAuth()
	if err != nil {
		return fmt.Errorf("could not get master auth data: %w", err)
	}

	hashToCompare := s.hashPassword(cmd.Password, salt)

	if subtle.ConstantTimeCompare(verificationHash, hashToCompare) != 1 {
		return errors.New("invalid master password")
	}

	s.encryptor.DeriveKeyFromPassword(cmd.Password, salt)
	return nil
}
