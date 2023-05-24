// Package encryptor дамет методы шифрования и дешифрования данных.
package encryptor

//go:generate mockgen -destination=mock/encryptor.go -source=encryptor.go

import (
	"encoding/json"

	"github.com/casnerano/seckeep/internal/client/model"
)

// Cipher интерфейс шифровщика и дешифроващика.
type Cipher interface {
	Encrypt(src []byte) ([]byte, error)
	Decrypt(dst []byte) ([]byte, error)
}

// Encryptor структура шифрования и дешифрования.
type Encryptor struct {
	cipher Cipher
}

// New конструктор.
func New(cipher Cipher) *Encryptor {
	return &Encryptor{
		cipher: cipher,
	}
}

// Encrypt метод шифрует данные.
func (e *Encryptor) Encrypt(dt model.DataTypeable) ([]byte, error) {
	bJSON, err := json.Marshal(dt)
	if err != nil {
		return nil, err
	}

	encrypted, err := e.cipher.Encrypt(bJSON)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

// Decrypt метод дешифрует данные.
func (e *Encryptor) Decrypt(encrypted []byte, dt model.DataTypeable) error {
	decrypted, err := e.cipher.Decrypt(encrypted)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decrypted, dt)
	if err != nil {
		return err
	}

	return nil
}
