package create

import (
	"os"
	"path/filepath"

	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/casnerano/seckeep/pkg/svalid"
	"github.com/spf13/cobra"
)

// NewDocumentCmd конструктор команды создания записи документа.
func NewDocumentCmd(dataService DataService) *cobra.Command {
	var name, file string

	cmd := cobra.Command{
		Use:   "document",
		Short: "Документ (файл)",
		Run: func(cmd *cobra.Command, args []string) {
			meta, _ := cmd.Flags().GetStringSlice("meta")

			bContent, err := os.ReadFile(file)
			if err != nil {
				cmd.Println("Не удалось прочитать файл.")
				return
			}

			if name == "" {
				name = filepath.Base(file)
			}

			d := model.DataDocument{
				Name:    name,
				Content: bContent,
				Meta:    meta,
			}

			validator := svalid.New()
			if err := validator.Validate(d); err != nil {
				cmd.Println(err.Error())
				return
			}

			err = dataService.Create(&d)

			if err != nil {
				cmd.Println(err)
				return
			}

			cmd.Println("Документ успешно сохранен.")
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "Путь к документу")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Название")

	_ = cmd.MarkFlagRequired("file")

	return &cmd
}
