package data

//go:generate mockgen -destination=mock/data.go -source=data.go

import (
	"github.com/casnerano/seckeep/internal/client/command/data/create"
	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/spf13/cobra"
)

// Service интерфейс взаимодействия с данными.
type Service interface {
	Create(dt model.DataTypeable) error
	Read(index int) (model.DataTypeable, error)
	GetList() map[int]model.DataTypeable
	Update(index int, dt model.DataTypeable) error
	Delete(index int) error
}

// SyncerService интерфейс синхронизации сервера и клиента.
type SyncerService interface {
	ServerHealthErr() error
	RunWithStatus()
}

// NewCmd конструктор базовой команды работы с данными.
// Содердит инициализацию дочерних команд.
func NewCmd(dataService Service, syncer SyncerService) *cobra.Command {
	cmd := cobra.Command{
		Use:   "data",
		Short: "Взаимодействие с данными",
	}

	cmd.AddCommand(create.NewCmd(dataService, syncer))
	cmd.AddCommand(NewReadCmd(dataService, syncer))
	cmd.AddCommand(NewListCmd(dataService, syncer))
	cmd.AddCommand(NewUpdateCmd(dataService, syncer))
	cmd.AddCommand(NewDeleteCmd(dataService, syncer))

	return &cmd
}
