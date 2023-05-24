// Package print дает методы для печати данных в разных представлениях.
package print

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/casnerano/seckeep/internal/shared"
)

// Print структура печати данных.
type Print struct {
	writer io.Writer
}

// New конструктор.
func New(writer io.Writer) *Print {
	return &Print{writer: writer}
}

// GroupedList метод печает набор данных сгрупированные по типу.
func (p *Print) GroupedList(dt map[int]model.DataTypeable) {
	groups := make(map[shared.DataType]map[int]model.DataTypeable)

	for index := range dt {
		_, ok := groups[dt[index].Type()]
		if !ok {
			groups[dt[index].Type()] = make(map[int]model.DataTypeable)
		}
		groups[dt[index].Type()][index] = dt[index]
	}

	grIndex := 0
	grCount := len(groups)

	for key := range groups {
		switch key {
		case shared.DataTypeCredential:
			fmt.Fprintln(p.writer, "Учетные записи:")
		case shared.DataTypeText:
			fmt.Fprintln(p.writer, "Текстовые данные:")
		case shared.DataTypeCard:
			fmt.Fprintln(p.writer, "Данные кредитных карт:")
		case shared.DataTypeDocument:
			fmt.Fprintln(p.writer, "Документы:")
		}

		sortedIndex := make([]int, 0, len(groups[key]))
		for index := range groups[key] {
			sortedIndex = append(sortedIndex, index)
		}

		sort.Ints(sortedIndex)

		for _, index := range sortedIndex {
			switch key {
			case shared.DataTypeCredential:
				if data, ok := groups[key][index].(*model.DataCredential); ok {
					fmt.Fprintf(
						p.writer,
						"#%d [ Логин: %s | Пароль: ***** | Мета: %s ]\n",
						index,
						data.Login,
						p.JoinedMetaString(data.Meta),
					)
				}
			case shared.DataTypeText:
				if data, ok := groups[key][index].(*model.DataText); ok {
					length := float64(len(data.Value))
					fmt.Fprintf(
						p.writer,
						"#%d [ Значение: %s | Мета: %s ]\n",
						index,
						data.Value[:int(length-math.Ceil(length/100*70))],
						p.JoinedMetaString(data.Meta),
					)
				}
			case shared.DataTypeCard:
				if data, ok := groups[key][index].(*model.DataCard); ok {
					ownerValue := "—"
					if data.Owner != "" {
						ownerValue = data.Owner
					}
					fmt.Fprintf(
						p.writer,
						"#%d [ Номер: %s | Месяц/Год: %s | CVV: *** | Держатель: %s | Мета: %s ]\n",
						index,
						data.Number,
						data.MonthYear,
						ownerValue,
						p.JoinedMetaString(data.Meta),
					)
				}
			case shared.DataTypeDocument:
				if data, ok := groups[key][index].(*model.DataDocument); ok {
					fmt.Fprintf(
						p.writer,
						"#%d [ Название: %s | Мета: %s ]\n",
						index,
						data.Name,
						p.JoinedMetaString(data.Meta),
					)
				}
			}
		}

		grIndex++
		if grIndex < grCount {
			fmt.Fprintln(p.writer)
		}
	}
}

// Detail метод печает детальную информацию данных.
func (p *Print) Detail(index int, dt model.DataTypeable) {
	switch dt.Type() {
	case shared.DataTypeCredential:
		if data, ok := dt.(*model.DataCredential); ok {
			fmt.Fprintf(
				p.writer,
				"Индекс: #%d\nЛогин: %s\nПароль: %s\nМета: %s",
				index,
				data.Login,
				data.Password,
				p.JoinedMetaString(data.Meta),
			)
		}
	case shared.DataTypeText:
		if data, ok := dt.(*model.DataText); ok {
			fmt.Fprintf(
				p.writer,
				"Индекс: #%d\nЗначение: %s\nМета: %s",
				index,
				data.Value,
				p.JoinedMetaString(data.Meta),
			)
		}
	case shared.DataTypeCard:
		if data, ok := dt.(*model.DataCard); ok {
			ownerValue := "—"
			if data.Owner != "" {
				ownerValue = data.Owner
			}
			fmt.Fprintf(
				p.writer,
				"Индекс: #%d\nНомер: %s\nМесяц/Год: %s\nCVV: %s\nДержатель: %s\nМета: %s",
				index,
				data.Number,
				data.MonthYear,
				data.CVV,
				ownerValue,
				p.JoinedMetaString(data.Meta),
			)
		}
	case shared.DataTypeDocument:
		if data, ok := dt.(*model.DataDocument); ok {
			fmt.Fprintf(
				p.writer,
				"Индекс: #%d\nНазвание: %s\nКонтент:\n=====\n%s\n=====\nМета: %s",
				index,
				data.Name,
				data.Content,
				p.JoinedMetaString(data.Meta),
			)
		}
	}
}

// JoinedMetaString метод объеденяет слайс тегов (строк) в строку.
func (p *Print) JoinedMetaString(meta []string) string {
	if len(meta) == 0 {
		return "—"
	}
	return strings.Join(meta, "; ")
}
