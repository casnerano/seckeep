package config

// FileName дефолтный путь к файлу конфигурации сервера.
const FileName = "./configs/server.yml"

// Config конфигурация сервера.
type Config struct {
	App struct {
		Authenticator struct {
			Secret string `yaml:"secret"`
		} `yaml:"authenticator"`
	} `yaml:"app"`
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	Database struct {
		DSN string `yaml:"dsn"`
	} `yaml:"database"`
}
