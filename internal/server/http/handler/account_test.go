package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mock_handler "github.com/casnerano/seckeep/internal/server/http/handler/mock"
	"github.com/casnerano/seckeep/internal/server/model"
	"github.com/casnerano/seckeep/internal/server/service/account"
	"github.com/casnerano/seckeep/pkg/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type AccountHandlerTestSuite struct {
	suite.Suite
	handler        *Account
	accountService *mock_handler.MockAccountService
}

func (s *AccountHandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.accountService = mock_handler.NewMockAccountService(ctrl)

	s.handler = NewAccount(s.accountService, log.NewStub())
}

func (s *AccountHandlerTestSuite) TestSignUpHandler() {
	rd := model.UserSignUpRequest{
		Login:    "ivan",
		Password: "ur3G28u%3fD",
		FullName: "Ivanov Ivan",
	}

	s.Run("New correct user", func() {
		s.accountService.EXPECT().SignUp(gomock.Any(), rd.Login, rd.Password, rd.FullName).Return("eyJhbGci.e30.Et9HFtf9R3GEM", nil)
		_, status := s.handler.SignUp(rd, httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/user/register", nil))
		s.Equal(http.StatusOK, status)
	})

	s.Run("Existing user", func() {
		s.accountService.EXPECT().SignUp(gomock.Any(), rd.Login, rd.Password, rd.FullName).Return("", account.ErrUserRegistered)
		_, status := s.handler.SignUp(rd, httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/user/register", nil))
		s.Equal(http.StatusConflict, status)
	})

	s.Run("Unknown error returning", func() {
		s.accountService.EXPECT().SignUp(gomock.Any(), rd.Login, rd.Password, rd.FullName).Return("", errors.New("unknown error"))
		_, status := s.handler.SignUp(rd, httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/user/register", nil))
		s.Equal(http.StatusInternalServerError, status)
	})
}

func (s *AccountHandlerTestSuite) TestSignInHandler() {
	rd := model.UserSignInRequest{
		Login:    "ivan",
		Password: "ur3G28u%3fD",
	}

	s.Run("User with correct credentials", func() {
		s.accountService.EXPECT().SignIn(gomock.Any(), rd.Login, rd.Password).Return("eyJhbGci.e30.Et9HFtf9R3GEM", nil)
		_, status := s.handler.SignIn(rd, httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/user/register", nil))
		s.Equal(http.StatusOK, status)
	})

	s.Run("User with incorrect credentials", func() {
		s.accountService.EXPECT().SignIn(gomock.Any(), rd.Login, rd.Password).Return("", account.ErrIncorrectCredentials)
		_, status := s.handler.SignIn(rd, httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/user/login", nil))
		s.Equal(http.StatusUnauthorized, status)
	})

	s.Run("Unknown error returning", func() {
		s.accountService.EXPECT().SignIn(gomock.Any(), rd.Login, rd.Password).Return("", errors.New("unknown error"))
		_, status := s.handler.SignIn(rd, httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/user/login", nil))
		s.Equal(http.StatusInternalServerError, status)
	})
}

func TestDataTestSuite(t *testing.T) {
	suite.Run(t, new(AccountHandlerTestSuite))
}
