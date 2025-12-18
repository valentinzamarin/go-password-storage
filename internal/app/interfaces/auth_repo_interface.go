package interfaces

type AuthRepository interface {
	GetMasterAuth() (salt []byte, verificationHash []byte, err error)
	CreateMasterAuth(salt, verificationHash []byte) error
}
