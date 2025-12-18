package encrypt

import (
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	keySize   = 32
	saltSize  = 16
	nonceSize = 24

	argonTime    = 1
	argonMemory  = 64 * 1024 
	argonThreads = 4
)

var (
	ErrKeyNotSet        = errors.New("encryption key is not set")
	ErrDataTooShort     = errors.New("encrypted data is too short")
	ErrDecryptionFailed = errors.New("decryption failed: invalid master password or corrupted data")
)

type PasswordEncrypt struct {
	encryptionKey *[keySize]byte
}

func NewPasswordEncrypt() *PasswordEncrypt {
	return &PasswordEncrypt{}
}

func (s *PasswordEncrypt) Encrypt(data []byte) ([]byte, error) {
	if s.encryptionKey == nil {
		return nil, ErrKeyNotSet
	}

	var nonce [nonceSize]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	encrypted := secretbox.Seal(nonce[:], data, &nonce, s.encryptionKey)
	return encrypted, nil
}

func (s *PasswordEncrypt) Decrypt(encryptedData []byte) ([]byte, error) {
	if s.encryptionKey == nil {
		return nil, ErrKeyNotSet
	}

	if len(encryptedData) < nonceSize {
		return nil, ErrDataTooShort
	}

	var nonce [nonceSize]byte
	copy(nonce[:], encryptedData[:nonceSize])

	decrypted, ok := secretbox.Open(nil, encryptedData[nonceSize:], &nonce, s.encryptionKey)
	if !ok {
		return nil, ErrDecryptionFailed
	}

	return decrypted, nil
}


func (s *PasswordEncrypt) DeriveKeyFromPassword(masterPassword string, salt []byte) {
	key := argon2.IDKey(
		[]byte(masterPassword),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		keySize,
	)

	s.encryptionKey = &[keySize]byte{}
	copy(s.encryptionKey[:], key)
}


func (s *PasswordEncrypt) IsKeySet() bool {
	return s.encryptionKey != nil
}


func (s *PasswordEncrypt) GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func (s *PasswordEncrypt) ClearKey() {
	if s.encryptionKey != nil {
		for i := range s.encryptionKey {
			s.encryptionKey[i] = 0
		}
		s.encryptionKey = nil
	}
}
