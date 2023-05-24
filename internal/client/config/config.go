package config

// FileName дефолтный путь к файлу конфигурации клиента.
const FileName = "./configs/client.yml"

// Config конфигурация клиента.
type Config struct {
	App struct {
		Encryptor struct {
			Secret string `yaml:"secret"`
		} `yaml:"encryptor"`
	} `yaml:"app"`
	Server struct {
		URL string `yaml:"url"`
	} `yaml:"server"`
}
