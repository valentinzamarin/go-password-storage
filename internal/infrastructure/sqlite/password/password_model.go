package passwords

type PasswordModel struct {
	ID                uint   `gorm:"primaryKey"`
	URL               string `gorm:"not null"`
	Login             string `gorm:"not null"`
	EncryptedPassword []byte `gorm:"not null;column:password"`
	Description       string `gorm:""`
}

/*
GORM generates table names
table becomes "passwords_model
to avoid this - rename it.
*/
func (PasswordModel) TableName() string {
	return "passwords"
}
