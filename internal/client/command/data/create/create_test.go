package create

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	mock_create "github.com/casnerano/seckeep/internal/client/command/data/create/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

var (
	errUnknown = errors.New("unknown error")
)

type DataCreateCmdTestSuite struct {
	suite.Suite
	dataService   *mock_create.MockDataService
	syncerService *mock_create.MockSyncerService
}

func (s *DataCreateCmdTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.dataService = mock_create.NewMockDataService(ctrl)
	s.syncerService = mock_create.NewMockSyncerService(ctrl)
}

func (s *DataCreateCmdTestSuite) TestCreateCmd() {
	cmd := NewCmd(s.dataService, s.syncerService)
	s.True(cmd.HasSubCommands())
}

func (s *DataCreateCmdTestSuite) TestCredential() {
	login := "ivan"
	password := "ivanov"

	cmd := NewCredentialCmd(s.dataService)
	cmdBuf := bytes.NewBufferString("")
	cmd.SetOut(cmdBuf)

	s.Run("Success create", func() {
		s.dataService.EXPECT().Create(gomock.Any()).Return(nil)

		cmd.SetArgs([]string{"-l", login, "-p", password})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), "Учетная запись успешно добавлена")
	})

	s.Run("Invalid password (validate)", func() {
		cmd.SetArgs([]string{"-l", login, "-p", ""})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), "Error:Field validation for 'Password'")
	})

	s.Run("Invalid create", func() {
		s.dataService.EXPECT().Create(gomock.Any()).Return(errUnknown)

		cmd.SetArgs([]string{"-l", login, "-p", password})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), errUnknown.Error())
	})
}

func (s *DataCreateCmdTestSuite) TestText() {
	value := "Example text"

	cmd := NewTextCmd(s.dataService)
	cmdBuf := bytes.NewBufferString("")
	cmd.SetOut(cmdBuf)

	s.Run("Success create", func() {
		s.dataService.EXPECT().Create(gomock.Any()).Return(nil)

		cmd.SetArgs([]string{"-v", value})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), "Текстовые данные успешно добавлены")
	})

	s.Run("Invalid create", func() {
		s.dataService.EXPECT().Create(gomock.Any()).Return(errUnknown)

		cmd.SetArgs([]string{"-v", value})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), errUnknown.Error())
	})
}

func (s *DataCreateCmdTestSuite) TestCard() {
	number := "4969677832915892"
	monthYear := "01.02"
	owner := "Ivan Ivanov"
	cvv := "123"

	cmd := NewCardCmd(s.dataService)
	cmdBuf := bytes.NewBufferString("")
	cmd.SetOut(cmdBuf)

	s.Run("Success create", func() {
		s.dataService.EXPECT().Create(gomock.Any()).Return(nil)

		cmd.SetArgs([]string{"-n", number, "-m", monthYear, "-o", owner, "-c", cvv})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), "Данные кредитной карты успешно добавлены")
	})

	s.Run("Invalid create", func() {
		s.dataService.EXPECT().Create(gomock.Any()).Return(errUnknown)

		cmd.SetArgs([]string{"-n", number, "-m", monthYear, "-o", owner, "-c", cvv})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), errUnknown.Error())
	})
}

func (s *DataCreateCmdTestSuite) TestDocument() {
	name := "temp_data.txt"

	cmd := NewDocumentCmd(s.dataService)
	cmdBuf := bytes.NewBufferString("")
	cmd.SetOut(cmdBuf)

	file, err := os.CreateTemp("", name)
	s.Require().NoError(err)

	s.T().Cleanup(func() {
		_ = os.Remove(file.Name())
	})

	_, err = file.WriteString("Example data")
	s.Require().NoError(err)

	s.Run("Success create", func() {
		s.dataService.EXPECT().Create(gomock.Any()).Return(nil)

		cmd.SetArgs([]string{"-f", file.Name(), "-n", name})
		err = cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), "Документ успешно сохранен")
	})

	s.Run("Invalid create", func() {
		s.dataService.EXPECT().Create(gomock.Any()).Return(errUnknown)

		cmd.SetArgs([]string{"-f", file.Name(), "-n", name})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), errUnknown.Error())
	})
}

func TestDataCreateCmdTestSuite(t *testing.T) {
	suite.Run(t, new(DataCreateCmdTestSuite))
}
