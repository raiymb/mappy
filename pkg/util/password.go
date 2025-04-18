package util

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns a bcrypt‑hashed string.
func HashPassword(pw string) (string, error) {
	out, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(out), err
}

// CheckPassword compares hash with plain pwd.
func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

// UUID generates a 16‑byte random hex (good enough for JTI).
func UUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
