package data

import (
	"fmt"

	"github.com/casnerano/seckeep/internal/client/model"
	smodel "github.com/casnerano/seckeep/internal/pkg/model"
	"github.com/casnerano/seckeep/pkg/svalid"
	"github.com/spf13/cobra"
)

type uQuestion struct {
	title        string
	currentValue string
	value        string
}

type uQuestions map[string]*uQuestion

func (uq uQuestions) ask() {
	for _, v := range uq {
		fmt.Printf("%s (текущее значение: %s) > ", v.title, v.currentValue)
		if _, err := fmt.Scanf("%s\n", &v.value); err != nil {
			v.value = v.currentValue
		}
	}
}

// NewUpdateCmd конструктор команда обновления записи по индексу.
func NewUpdateCmd(dataService Service, syncer SyncerService) *cobra.Command {
	var index int

	cmd := cobra.Command{
		Use:   "update",
		Short: "Обновление",
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
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

			var updatedData model.DataTypeable

			switch d.Type() {
			case smodel.DataTypeCredential:
				if credential, ok := d.(*model.DataCredential); ok {
					questions := uQuestions{
						"login":    {title: "Логин", currentValue: credential.Login},
						"password": {title: "Пароль", currentValue: credential.Password},
					}

					questions.ask()

					credential.Login = questions["login"].value
					credential.Password = questions["password"].value

					updatedData = credential
				}
			case smodel.DataTypeText:
				if text, ok := d.(*model.DataText); ok {
					questions := uQuestions{
						"value": {title: "Значение", currentValue: text.Value},
					}

					questions.ask()

					text.Value = questions["value"].value

					updatedData = text
				}
			case smodel.DataTypeCard:
				if card, ok := d.(*model.DataCard); ok {
					questions := uQuestions{
						"number":     {title: "Номер карты", currentValue: card.Number},
						"month-year": {title: "Месяц/Год", currentValue: card.MonthYear},
						"owner":      {title: "Держатель", currentValue: card.Owner},
						"cvv":        {title: "CVV", currentValue: card.CVV},
					}

					questions.ask()

					card.Number = questions["number"].value
					card.MonthYear = questions["month-year"].value
					card.CVV = questions["owner"].value
					card.Owner = questions["cvv"].value

					updatedData = card
				}
			case smodel.DataTypeDocument:
				if document, ok := d.(*model.DataDocument); ok {
					questions := uQuestions{
						"name": {title: "Название", currentValue: document.Name},
					}

					questions.ask()

					document.Name = questions["name"].value

					updatedData = document
				}
			default:
				cmd.Println("Неизвестный тип данных.")
				return
			}

			validator := svalid.New()
			if err = validator.Validate(updatedData); err != nil {
				cmd.Println(err.Error())
				return
			}

			if err = dataService.Update(index, updatedData); err != nil {
				cmd.Println(err.Error())
			}

			cmd.Println("Данные успешно обновлены.")
		},
	}

	cmd.Flags().IntVarP(&index, "index", "i", 0, "Индекс (номер) записи")
	_ = cmd.MarkFlagRequired("index")

	return &cmd
}
