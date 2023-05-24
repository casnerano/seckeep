package yaml

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testConfig struct {
	App struct {
		Encryptor struct {
			Secret string `yaml:"secret"`
		} `yaml:"encryptor"`
	} `yaml:"app"`
	Server struct {
		URL string `json:"url"`
	}
}

func TestLoadFromFile(t *testing.T) {
	yamlContent := `
app:
  encryptor:
    secret: "c4ca4238a0b923820dcc509a6f75849b"
server:
  url: http://127.0.0.1:8081
`

	want := testConfig{}
	want.App.Encryptor.Secret = "c4ca4238a0b923820dcc509a6f75849b"
	want.Server.URL = "http://127.0.0.1:8081"

	file, err := os.CreateTemp("", "config.yaml")
	require.NoError(t, err)
	_, err = file.WriteString(yamlContent)
	t.Cleanup(func() {
		_ = os.Remove(file.Name())
	})

	t.Run("Read correct yaml for matched structure", func(t *testing.T) {
		got := testConfig{}
		err = LoadFromFile(file.Name(), &got)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Read non-exist file", func(t *testing.T) {
		got := testConfig{}
		err = LoadFromFile("not-file", &got)
		assert.Error(t, err)
	})

	t.Run("Read incorrect yaml", func(t *testing.T) {
		got := testConfig{}
		_, err = file.WriteString("typo")
		require.NoError(t, err)

		err = LoadFromFile(file.Name(), &got)
		require.Error(t, err)
	})
}
