package account

import "github.com/spf13/cobra"

// NewSignInCmd конструктор команда авторизации пользователя на сервере.
func NewSignInCmd(account Service) *cobra.Command {
	var login, password string

	cmd := cobra.Command{
		Use:   "sign-in",
		Short: "Авторизация",
		Run: func(cmd *cobra.Command, args []string) {
			if err := account.SignIn(login, password); err != nil {
				cmd.Println(err)
				return
			}
			cmd.Println("Успешная авторизация ;)")
		},
	}

	cmd.Flags().StringVarP(&login, "login", "l", "", "Логин")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Пароль")

	_ = cmd.MarkFlagRequired("login")
	_ = cmd.MarkFlagRequired("password")

	return &cmd
}
