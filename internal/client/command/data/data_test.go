package data

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"testing"

	mock_data "github.com/casnerano/seckeep/internal/client/command/data/mock"
	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

var (
	errUnknown = errors.New("unknown error")
)

type DataCmdTestSuite struct {
	suite.Suite
	dataService   *mock_data.MockDataService
	syncerService *mock_data.MockSyncerService
}

func (s *DataCmdTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.dataService = mock_data.NewMockDataService(ctrl)
	s.syncerService = mock_data.NewMockSyncerService(ctrl)
}

func (s *DataCmdTestSuite) TestDataCmd() {
	cmd := NewCmd(s.dataService, s.syncerService)
	s.True(cmd.HasSubCommands())
}

func (s *DataCmdTestSuite) TestDelete() {
	s.syncerService.EXPECT().ServerHealthErr().Return(nil).AnyTimes()
	s.syncerService.EXPECT().RunWithStatus().AnyTimes()
	index := 3

	cmd := NewDeleteCmd(s.dataService, s.syncerService)
	cmdBuf := bytes.NewBufferString("")
	cmd.SetOut(cmdBuf)

	cmd.SetArgs([]string{"-i", strconv.Itoa(index)})

	s.Run("Success delete", func() {
		s.dataService.EXPECT().Delete(index).Return(nil)

		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), "Запись успешно удалена")
	})

	s.Run("Invalid delete", func() {
		s.dataService.EXPECT().Delete(index).Return(errUnknown)

		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), errUnknown.Error())
	})
}

func (s *DataCmdTestSuite) TestRead() {
	s.syncerService.EXPECT().ServerHealthErr().Return(nil).AnyTimes()
	s.syncerService.EXPECT().RunWithStatus().AnyTimes()
	index := 3

	cmd := NewReadCmd(s.dataService, s.syncerService)
	cmdBuf := bytes.NewBufferString("")
	cmd.SetOut(cmdBuf)

	cmd.SetArgs([]string{"-i", strconv.Itoa(index)})

	s.Run("Success delete", func() {
		dt := model.DataText{
			Value: "Example text",
			Meta:  nil,
		}

		s.dataService.EXPECT().Read(index).Return(&dt, nil)

		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), dt.Value)
	})

	s.Run("Invalid delete", func() {
		s.dataService.EXPECT().Read(index).Return(nil, errUnknown)

		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), errUnknown.Error())
	})
}

func (s *DataCmdTestSuite) TestList() {
	s.syncerService.EXPECT().ServerHealthErr().Return(nil).AnyTimes()
	s.syncerService.EXPECT().RunWithStatus().AnyTimes()

	cmd := NewListCmd(s.dataService, s.syncerService)
	cmdBuf := bytes.NewBufferString("")
	cmd.SetOut(cmdBuf)

	s.Run("Good list", func() {
		dt := map[int]model.DataTypeable{
			0: &model.DataText{Value: "Example #1 Text", Meta: []string{"Tag1", "Tag2"}},
			1: &model.DataText{Value: "Example #2 Text", Meta: []string{"Tag1"}},
		}

		s.dataService.EXPECT().GetList().Return(dt)

		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.NotEmpty(string(out))
	})

	s.Run("Empty list", func() {
		dt := map[int]model.DataTypeable{}
		s.dataService.EXPECT().GetList().Return(dt)

		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), "Список записей пуст")
	})
}

func (s *DataCmdTestSuite) TestUpdate() {
	s.syncerService.EXPECT().ServerHealthErr().Return(nil).AnyTimes()
	s.syncerService.EXPECT().RunWithStatus().AnyTimes()

	cmd := NewUpdateCmd(s.dataService, s.syncerService)
	cmdBuf := bytes.NewBufferString("")
	cmd.SetOut(cmdBuf)

	dataTypeList := []model.DataTypeable{
		&model.DataText{Value: "Example #1 Text", Meta: []string{"Tag1", "Tag2"}},
		&model.DataCredential{Login: "example-l", Password: "example-p", Meta: nil},
		&model.DataCard{Number: "4969677832915892", MonthYear: "01.02", CVV: "123", Owner: "Ivan Ivanov", Meta: nil},
		&model.DataDocument{Name: "Example.Name", Content: []byte("Example content"), Meta: nil},
	}

	dataCount := len(dataTypeList)

	for index := range dataTypeList {
		s.Run("Success update", func() {
			s.dataService.EXPECT().Read(index).Return(dataTypeList[index], nil).MaxTimes(dataCount)
			s.dataService.EXPECT().Update(index, gomock.Any()).Return(nil).MaxTimes(dataCount)

			cmd.SetArgs([]string{"-i", strconv.Itoa(index)})
			err := cmd.Execute()
			s.Require().NoError(err)

			out, err := io.ReadAll(cmdBuf)
			s.Require().NoError(err)

			s.Contains(string(out), "Данные успешно обновлены")
		})
	}

	s.Run("Incorrect read", func() {
		index := dataCount + 1
		s.dataService.EXPECT().Read(index).Return(nil, errUnknown)

		cmd.SetArgs([]string{"-i", strconv.Itoa(index)})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), errUnknown.Error())
	})
}

func TestDataCmdTestSuite(t *testing.T) {
	suite.Run(t, new(DataCmdTestSuite))
}
