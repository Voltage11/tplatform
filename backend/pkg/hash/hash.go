package hash

import (
	"crypto/rand"
	"math/big"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetRandomString() string {
	// Попробуем сгенерировать UUID
	uid, err := uuid.NewRandom()
	if err == nil {
		return uid.String()
	}

	// Резерв: генерация 32-символьной строки
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const length = 32
	result := make([]byte, length)

	for i := range result {
		// Генерируем случайное число от 0 до len(chars)-1
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[n.Int64()]
	}

	return string(result)
}
