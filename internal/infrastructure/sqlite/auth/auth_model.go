package auth

type AuthModel struct {
	ID               uint   `gorm:"primaryKey"`
	Salt             []byte `gorm:"not null"`
	VerificationHash []byte `gorm:"not null"`
}
