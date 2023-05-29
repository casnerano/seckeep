package cipher

import (
	"bytes"
	"testing"
)

func TestCipher(t *testing.T) {
	tests := []struct {
		name    string
		rawText []byte
	}{
		{"Short text", []byte("Lorem")},
		{"Middle text", []byte("Lorem ipsum")},
		{"Large text", []byte("Lorem ipsum is placeholder text commonly used in the graphic")},
	}

	c := New([]byte("example key"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cipherText, err := c.Encrypt(tt.rawText)
			if err != nil {
				t.Errorf("Encrypt() error = %v", err)
				return
			}
			rawText, err := c.Decrypt(cipherText)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
				return
			}

			if !bytes.Equal(tt.rawText, rawText) {
				t.Errorf("The decrypted text does not match the encrypted.")
			}
		})
	}
}
