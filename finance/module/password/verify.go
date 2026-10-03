package password

import "golang.org/x/crypto/bcrypt"

func Verify(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
