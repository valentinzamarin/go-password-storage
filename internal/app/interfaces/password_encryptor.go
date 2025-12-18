package interfaces

type PasswordEncryptor interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
	GenerateSalt() ([]byte, error)
	DeriveKeyFromPassword(masterPassword string, salt []byte)
}
