package create

import (
	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/casnerano/seckeep/pkg/svalid"
	"github.com/spf13/cobra"
)

// NewTextCmd конструктор команды создания записи простого текста.
func NewTextCmd(dataService DataService) *cobra.Command {
	var value string

	cmd := cobra.Command{
		Use:   "text",
		Short: "Текстовая информация",
		Run: func(cmd *cobra.Command, args []string) {
			meta, _ := cmd.Flags().GetStringSlice("meta")

			d := model.DataText{
				Value: value,
				Meta:  meta,
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

			cmd.Println("Текстовые данные успешно добавлены.")
		},
	}

	cmd.Flags().StringVarP(&value, "value", "v", "", "Значение")

	_ = cmd.MarkFlagRequired("value")

	return &cmd
}
