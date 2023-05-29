package create

//go:generate mockgen -destination=mock/create.go -source=create.go

import (
	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/spf13/cobra"
)

// DataService интерфейс взаимодействия с данными.
type DataService interface {
	Create(dt model.DataTypeable) error
}

// SyncerService интерфейс синхронизации сервера и клиента.
type SyncerService interface {
	ServerHealthErr() error
	RunWithStatus()
}

// NewCmd конструктор базовой команды для создания данных.
// Содердит инициализацию дочерних команд.
func NewCmd(dataService DataService, syncer SyncerService) *cobra.Command {
	cmd := cobra.Command{
		Use:   "create",
		Short: "Создание",
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if syncer.ServerHealthErr() == nil {
				syncer.RunWithStatus()
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.PersistentFlags().StringSlice("meta", []string{}, "Мета данные")

	cmd.AddCommand(NewCredentialCmd(dataService))
	cmd.AddCommand(NewTextCmd(dataService))
	cmd.AddCommand(NewCardCmd(dataService))
	cmd.AddCommand(NewDocumentCmd(dataService))

	return &cmd
}
