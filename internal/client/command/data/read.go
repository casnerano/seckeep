package data

import (
	"github.com/casnerano/seckeep/internal/client/service/data/print"
	"github.com/spf13/cobra"
)

// NewReadCmd конструктор команда вывода записи по индексу.
func NewReadCmd(dataService Service, syncer SyncerService) *cobra.Command {
	var index int

	cmd := cobra.Command{
		Use:   "read",
		Short: "Чтение",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Parent() != nil && cmd.Parent().Parent() != nil {
				cmd.Parent().Parent().PersistentPreRun(cmd.Parent(), args)
			}

			if syncer.ServerHealthErr() == nil {
				syncer.RunWithStatus()
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			d, err := dataService.Read(index)
			if err != nil {
				cmd.Println(err.Error())
				return
			}
			p := print.New(cmd.OutOrStdout())
			p.Detail(index, d)
		},
	}

	cmd.Flags().IntVarP(&index, "index", "i", 0, "Индекс (номер) записи")
	_ = cmd.MarkFlagRequired("index")

	return &cmd
}
