// Package yaml модержит функции загрузчика yaml-конфигурации.
package yaml

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadFromFile загрузка конфигурации из файла.
func LoadFromFile(fName string, config any) error {
	file, err := os.ReadFile(fName)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return err
	}

	return nil
}
