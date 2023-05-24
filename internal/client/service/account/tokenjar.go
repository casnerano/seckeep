package account

import (
	"os"
)

const tokenJarFileName = "./cmd/client/var/token.jar"

// TokenJar структура для работы с хранилищем токена.
type TokenJar struct{}

// NewTokenJar конструктор.
func NewTokenJar() *TokenJar {
	return &TokenJar{}
}

// SetToken метод устанавливает (сохраняет) токен в локальный файл.
func (tj TokenJar) SetToken(token string) error {
	file, err := os.OpenFile(tokenJarFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		return nil
	}
	defer file.Close()

	_, err = file.WriteString(token)
	return err
}

// ReadToken метод читает токен из локального файла.
func (tj TokenJar) ReadToken() (string, error) {
	if _, err := os.Stat(tokenJarFileName); err != nil {
		return "", err
	}

	bToken, err := os.ReadFile(tokenJarFileName)
	if err != nil {
		return "", err
	}

	return string(bToken), nil
}
