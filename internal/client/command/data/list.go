package data

import (
	"github.com/casnerano/seckeep/internal/client/service/data/print"
	"github.com/spf13/cobra"
)

// NewListCmd конструктор команда вывода списка записей.
func NewListCmd(dataService Service, syncer SyncerService) *cobra.Command {
	cmd := cobra.Command{
		Use:   "list",
		Short: "Список",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().Parent() != nil {
				cmd.Parent().Parent().PersistentPreRun(cmd.Parent(), args)
			}

			if syncer.ServerHealthErr() == nil {
				syncer.RunWithStatus()
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			dList := dataService.GetList()
			if len(dList) == 0 {
				cmd.Println("Список записей пуст.")
				return
			}
			p := print.New(cmd.OutOrStdout())
			p.GroupedList(dList)
		},
	}

	return &cmd
}
