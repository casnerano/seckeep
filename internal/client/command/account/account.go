package account

//go:generate mockgen -destination=mock/account.go -source=account.go

import (
	"github.com/casnerano/seckeep/internal/client/service/account"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// Service интерфейс взаимодействия с аккаунтом.
type Service interface {
	SignUp(login, password, fullName string) error
	SignIn(login, password string) error
}

// NewCmd конструктор базовой команды взаимодействия с аккаунтом пользователя.
// Содердит инициализацию дочерних команд.
func NewCmd(client *resty.Client) *cobra.Command {
	cmd := cobra.Command{
		Use:              "account",
		Short:            "Взаимодействие с аккаунтом пользователя",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	}

	accountService := account.New(client, account.NewTokenJar())
	cmd.AddCommand(NewSignInCmd(accountService))
	cmd.AddCommand(NewSignUpCmd(accountService))

	return &cmd
}
