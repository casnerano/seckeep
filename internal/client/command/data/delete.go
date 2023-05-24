package data

import (
	"github.com/spf13/cobra"
)

// NewDeleteCmd конструктор команда удаления записи по индексу.
func NewDeleteCmd(dataService Service, syncer SyncerService) *cobra.Command {
	var index int

	cmd := cobra.Command{
		Use:   "delete",
		Short: "Удаление",
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if syncer.ServerHealthErr() == nil {
				syncer.RunWithStatus()
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := dataService.Delete(index); err != nil {
				cmd.Println(err.Error())
				return
			}
			cmd.Println("Запись успешно удалена.")
		},
	}

	cmd.Flags().IntVarP(&index, "index", "i", 0, "Индекс (номер) записи")
	_ = cmd.MarkFlagRequired("index")

	return &cmd
}
