package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a bcrypt hash of the plaintext password.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// ComparePassword checks a plaintext password against a stored bcrypt hash.
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
