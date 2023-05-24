package account

import (
	"bytes"
	"io"
	"testing"

	mock_account "github.com/casnerano/seckeep/internal/client/command/account/mock"
	"github.com/casnerano/seckeep/internal/client/service/account"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type AccountTestSuite struct {
	suite.Suite
	accountService *mock_account.MockAccountService
}

func (s *AccountTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.accountService = mock_account.NewMockAccountService(ctrl)
}

func (s *AccountTestSuite) TestAccountCmd() {
	cmd := NewCmd(resty.New())
	s.True(cmd.HasSubCommands())
}

func (s *AccountTestSuite) TestSignIn() {
	cmd := NewSignInCmd(s.accountService)
	cmdBuf := bytes.NewBufferString("")
	cmd.SetOut(cmdBuf)

	s.Run("Success sing-in", func() {
		login := "ivan"
		password := "example"

		s.accountService.EXPECT().SignIn(login, password).Return(nil)

		cmd.SetArgs([]string{"-l", login, "-p", password})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), "Успешная авторизация")
	})

	s.Run("Incorrect credentials", func() {
		login := "ivan"
		password := "example"

		s.accountService.EXPECT().SignIn(login, password).Return(account.ErrIncorrectCredentials)

		cmd.SetArgs([]string{"-l", login, "-p", password})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), account.ErrIncorrectCredentials.Error())
	})

	s.Run("Incorrect args", func() {
		login := "ivan"
		password := "example"

		s.accountService.EXPECT().SignIn(login, password).Return(account.ErrIncorrectCredentials)

		cmd.SetArgs([]string{"-l", login, "--typo", password})
		err := cmd.Execute()

		s.Error(err)
	})
}

func (s *AccountTestSuite) TestSignUp() {
	cmd := NewSignUpCmd(s.accountService)
	cmdBuf := bytes.NewBufferString("")
	cmd.SetOut(cmdBuf)

	s.Run("Success sing-up", func() {
		login := "ivan"
		password := "example"
		fullName := "Ivan Ivanov"

		s.accountService.EXPECT().SignUp(login, password, fullName).Return(nil)

		cmd.SetArgs([]string{"-l", login, "-p", password, "-n", fullName})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), "Успешная регистрация")
	})

	s.Run("Existing user", func() {
		login := "ivan"
		password := "example"
		fullName := "Ivan Ivanov"

		s.accountService.EXPECT().SignUp(login, password, fullName).Return(account.ErrUserRegistered)

		cmd.SetArgs([]string{"-l", login, "-p", password, "-n", fullName})
		err := cmd.Execute()
		s.Require().NoError(err)

		out, err := io.ReadAll(cmdBuf)
		s.Require().NoError(err)

		s.Contains(string(out), account.ErrUserRegistered.Error())
	})

	s.Run("Incorrect args", func() {
		login := "ivan"
		password := "example"
		fullName := "Ivan Ivanov"

		s.accountService.EXPECT().SignUp(login, password, fullName).Return(account.ErrIncorrectCredentials)

		cmd.SetArgs([]string{"-l", login, "-p", password, "--typo", fullName})
		err := cmd.Execute()

		s.Error(err)
	})
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}
