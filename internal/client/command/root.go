// Package command содержит корневую команду.
package command

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/casnerano/seckeep/internal/client/command/account"
	"github.com/casnerano/seckeep/internal/client/command/data"
	"github.com/casnerano/seckeep/internal/client/config"
	aService "github.com/casnerano/seckeep/internal/client/service/account"
	dService "github.com/casnerano/seckeep/internal/client/service/data"
	"github.com/casnerano/seckeep/internal/client/service/data/encryptor"
	"github.com/casnerano/seckeep/internal/client/service/storage"
	"github.com/casnerano/seckeep/internal/client/service/syncer"
	"github.com/casnerano/seckeep/pkg/cipher"
	"github.com/casnerano/seckeep/pkg/log"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

// Root структура коневой команды.
type Root struct {
	cmd *cobra.Command
}

// RootCommandContext контекст команды.
// Необходимо для удобной передачи дочерним командам.
type RootCommandContext struct {
	Config      *config.Config
	Logger      log.Loggable
	DataStorage *storage.Storage
}

// NewRoot конструктор корневой команды.
func NewRoot(ctx *RootCommandContext) *Root {
	httpClient := resty.New()
	httpClient.SetBaseURL(ctx.Config.Server.URL + "/api")

	dataService := dService.New(
		ctx.DataStorage,
		encryptor.New(
			cipher.New([]byte(ctx.Config.App.Encryptor.Secret)),
		),
	)

	sync := syncer.New(httpClient, ctx.DataStorage, ctx.Logger)

	cmd := &cobra.Command{
		Use:   "seckeep",
		Short: "Менеджер секретных данных",
		Long: "Приложение позволяет хранить секретные данные в зашифрованном виде,\n" +
			"и синхронизировать между несколькими клиентами. \n" +
			"Поджробная информация — https://github.com/casnerano/seckeep",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			welcome := "+  SecKeep — менеджер секретных данных  +"
			length := utf8.RuneCountInString(welcome)

			fmt.Println(strings.Repeat("+", length))
			fmt.Println(welcome)
			fmt.Println(strings.Repeat("+", length))

			tj := aService.NewTokenJar()
			if token, err := tj.ReadToken(); err == nil {
				httpClient.SetAuthToken(token)
			}

			if err := sync.PingServerHealth(); err != nil {
				if errors.Is(err, syncer.ErrUnauthorized) {
					fmt.Println("Отсутсвует авторизация.")
				} else {
					fmt.Println("Отсутсвует соединение с сервером.")
				}
				fmt.Println("Клиент работает в локальном режиме.")
			} else {
				fmt.Println("Соединение с сервером установлено.")
			}

			fmt.Println()
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(account.NewCmd(httpClient))
	cmd.AddCommand(data.NewCmd(dataService, sync))

	return &Root{
		cmd: cmd,
	}
}

// Execute метод запуска команды.
func (r Root) Execute() error {
	return r.cmd.Execute()
}
