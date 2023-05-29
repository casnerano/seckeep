package create

import (
	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/casnerano/seckeep/pkg/svalid"
	"github.com/spf13/cobra"
)

// NewCardCmd конструктор команды создания записи банковской карты.
func NewCardCmd(dataService DataService) *cobra.Command {
	var number, monthYear, cvv, owner string

	cmd := cobra.Command{
		Use:   "card",
		Short: "Кредитная карта",
		Run: func(cmd *cobra.Command, args []string) {
			meta, _ := cmd.Flags().GetStringSlice("meta")

			d := model.DataCard{
				Number:    number,
				MonthYear: monthYear,
				CVV:       cvv,
				Owner:     owner,
				Meta:      meta,
			}

			validator := svalid.New()
			if err := validator.Validate(d); err != nil {
				cmd.Println(err.Error())
				return
			}

			err := dataService.Create(d)

			if err != nil {
				cmd.Println(err.Error())
				return
			}

			cmd.Println("Данные кредитной карты успешно добавлены.")
		},
	}

	cmd.Flags().StringVarP(&number, "number", "n", "", "Номер карты")
	cmd.Flags().StringVarP(&monthYear, "month-year", "m", "", "Месяц/Год")
	cmd.Flags().StringVarP(&owner, "owner", "o", "", "Держатель")
	cmd.Flags().StringVarP(&cvv, "cvv", "c", "", "CVV")

	_ = cmd.MarkFlagRequired("number")
	_ = cmd.MarkFlagRequired("month-year")
	_ = cmd.MarkFlagRequired("cvv")

	return &cmd
}
