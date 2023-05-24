package print

import (
	"bytes"
	"testing"

	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/stretchr/testify/suite"
)

type DataPrintTestSuite struct {
	suite.Suite
	dt     map[int]model.DataTypeable
	print  *Print
	output *bytes.Buffer
}

func (s *DataPrintTestSuite) SetupSuite() {
	s.dt = map[int]model.DataTypeable{
		0: &model.DataText{Value: "Example #1 Text", Meta: []string{"Tag1", "Tag2"}},
		1: &model.DataText{Value: "Example #2 Text", Meta: []string{"Tag1"}},
		2: &model.DataCredential{Login: "example-l", Password: "example-p", Meta: nil},
		3: &model.DataCard{Number: "123456789123456", MonthYear: "01.02", CVV: "123", Owner: "Ivan Ivanov", Meta: nil},
		4: &model.DataDocument{Name: "Example.Name", Content: []byte("Example content"), Meta: nil},
	}
	s.output = new(bytes.Buffer)
	s.print = New(s.output)
}

func (s *DataPrintTestSuite) SetupTestSuite() {
	s.output.Reset()
}

func (s *DataPrintTestSuite) TestGroupedList() {
	s.print.GroupedList(s.dt)
	stOutput := s.output.String()

	s.Run("Text items data output", func() {
		s.Contains(stOutput, "Текстовые данные:")
		s.Contains(stOutput, "#0 [ Значение: Exam | Мета: Tag1; Tag2 ]")
		s.Contains(stOutput, "#1 [ Значение: Exam | Мета: Tag1 ]")
	})

	s.Run("Credential items data output", func() {
		s.Contains(stOutput, "Учетные записи:")
		s.Contains(stOutput, "#2 [ Логин: example-l | Пароль: ***** | Мета: — ]")
	})

	s.Run("Card items data output", func() {
		s.Contains(stOutput, "Данные кредитных карт:")
		s.Contains(stOutput, "#3 [ Номер: 123456789123456 | Месяц/Год: 01.02 | CVV: *** | Держатель: Ivan Ivanov | Мета: — ]")
	})

	s.Run("Document items data output", func() {
		s.Contains(stOutput, "Документы:")
		s.Contains(stOutput, "#4 [ Название: Example.Name | Мета: — ]")
	})
}

func (s *DataPrintTestSuite) TestDetail() {
	s.Run("Text detail data output", func() {
		s.print.Detail(0, s.dt[0])

		stOutput := s.output.String()
		s.output.Reset()

		s.Contains(stOutput, "Индекс: #0\nЗначение: Example #1 Text\nМета: Tag1; Tag2")
	})

	s.Run("Credential detail data output", func() {
		s.print.Detail(2, s.dt[2])

		stOutput := s.output.String()
		s.output.Reset()

		s.Contains(stOutput, "Индекс: #2\nЛогин: example-l\nПароль: example-p\nМета: —")
	})

	s.Run("Card detail data output", func() {
		s.print.Detail(3, s.dt[3])

		stOutput := s.output.String()
		s.output.Reset()

		s.Contains(stOutput, "Индекс: #3\nНомер: 123456789123456\nМесяц/Год: 01.02\nCVV: 123\nДержатель: Ivan Ivanov\nМета: —")
	})

	s.Run("Document detail data output", func() {
		s.print.Detail(4, s.dt[4])

		stOutput := s.output.String()
		s.output.Reset()

		s.Contains(stOutput, "Индекс: #4\nНазвание: Example.Name\nКонтент:\n=====\nExample content\n=====\nМета: —")
	})
}

func (s *DataPrintTestSuite) TestJoinedMetaString() {
	tests := []struct {
		name string
		meta []string
		want string
	}{
		{"Single meta", []string{"Single"}, "Single"},
		{"Many meta", []string{"One", "Two", "Three"}, "One; Two; Three"},
		{"Empty", []string{}, "—"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Equal(tt.want, s.print.JoinedMetaString(tt.meta))
		})
	}
}

func TestDataTestSuite(t *testing.T) {
	suite.Run(t, new(DataPrintTestSuite))
}
