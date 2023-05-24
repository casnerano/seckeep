package create

import (
	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/casnerano/seckeep/pkg/svalid"
	"github.com/spf13/cobra"
)

// NewCredentialCmd конструктор команды создания записи учетной записи.
func NewCredentialCmd(dataService DataService) *cobra.Command {
	var login, password string

	cmd := cobra.Command{
		Use:   "credential",
		Short: "Учетные данные",
		Run: func(cmd *cobra.Command, args []string) {
			meta, _ := cmd.Flags().GetStringSlice("meta")

			d := model.DataCredential{
				Login:    login,
				Password: password,
				Meta:     meta,
			}

			validator := svalid.New()
			if err := validator.Validate(d); err != nil {
				cmd.Println(err.Error())
				return
			}

			err := dataService.Create(d)

			if err != nil {
				cmd.Println(err)
				return
			}

			cmd.Println("Учетная запись успешно добавлена.")
		},
	}

	cmd.Flags().StringVarP(&login, "login", "l", "", "Логин")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Пароль")

	_ = cmd.MarkFlagRequired("login")
	_ = cmd.MarkFlagRequired("password")

	return &cmd
}
