package account

import "github.com/spf13/cobra"

// NewSignUpCmd конструктор команда регистрации пользователя на сервере.
func NewSignUpCmd(account Service) *cobra.Command {
	var login, password, name string

	cmd := cobra.Command{
		Use:   "sign-up",
		Short: "Регистрация",
		Run: func(cmd *cobra.Command, args []string) {
			if err := account.SignUp(login, password, name); err != nil {
				cmd.Println(err)
				return
			}
			cmd.Println("Успешная регистрация ;)")
		},
	}

	cmd.Flags().StringVarP(&login, "login", "l", "", "Логин")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Пароль")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Имя")

	_ = cmd.MarkFlagRequired("login")
	_ = cmd.MarkFlagRequired("password")
	_ = cmd.MarkFlagRequired("name")

	return &cmd
}
