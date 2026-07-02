package password

import (
    "golang.org/x/crypto/bcrypt"
)

const cost = 12 

// Hash возвращает bcrypt-хеш пароля
func Hash(plain string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(plain), cost)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

// Verify сравнивает пароль с bcrypt-хешом
func Verify(hashed, plain string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
    return err == nil
}