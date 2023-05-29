package account

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/casnerano/seckeep/internal/server/model"
	"github.com/casnerano/seckeep/internal/server/repository"
	mock_repository "github.com/casnerano/seckeep/internal/server/repository/mock"
	mock_account "github.com/casnerano/seckeep/internal/server/service/account/mock"
	"github.com/casnerano/seckeep/pkg/jwtoken"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

var (
	errUnknown = errors.New("unknown error")
)

type AccountTestSuite struct {
	suite.Suite
	accountService *Account
	userRepo       *mock_repository.MockUser
	jwt            *mock_account.MockJWT
	secret         string
}

func (s *AccountTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.userRepo = mock_repository.NewMockUser(ctrl)
	s.jwt = mock_account.NewMockJWT(ctrl)
	s.secret = "for-example-secret"

	s.accountService = New(s.userRepo, s.jwt, s.secret)
}

func (s *AccountTestSuite) TestSignUp() {
	user := model.User{
		UUID:      "f9bd9622-f730-11ed-b67e-0242ac120002",
		Login:     "ivan",
		Password:  "iVm20%02fD5O",
		FullName:  "Ivanov Ivan",
		CreatedAt: time.Now(),
	}

	s.Run("Non-existing user", func() {
		jwtPayload := jwtoken.Payload{
			UUID:     user.UUID,
			FullName: user.FullName,
		}

		wantToken := "eyJhbGci.e30.Et9HFtf9R3GEM"

		s.userRepo.EXPECT().Add(gomock.Any(), user.Login, gomock.Any(), user.FullName).Return(&user, nil)
		s.jwt.EXPECT().Create(jwtPayload, jwtTTL, []byte(s.secret)).Return(wantToken, nil)

		gotToken, err := s.accountService.SignUp(context.Background(), user.Login, user.Password, user.FullName)
		s.Require().NoError(err)

		s.Equal(wantToken, gotToken)
	})

	s.Run("Existing user", func() {
		s.userRepo.EXPECT().Add(gomock.Any(), user.Login, gomock.Any(), user.FullName).Return(nil, repository.ErrAlreadyExist)
		gotToken, err := s.accountService.SignUp(context.Background(), user.Login, user.Password, user.FullName)

		s.Empty(gotToken)
		s.ErrorIs(err, ErrUserRegistered)
	})

	s.Run("Unknown error", func() {
		s.userRepo.EXPECT().Add(gomock.Any(), user.Login, gomock.Any(), user.FullName).Return(nil, errUnknown)
		gotToken, err := s.accountService.SignUp(context.Background(), user.Login, user.Password, user.FullName)

		s.Empty(gotToken)
		s.ErrorIs(err, errUnknown)
	})

	s.Run("BCrypt password generate error", func() {
		gotToken, err := s.accountService.SignUp(context.Background(), user.Login, strings.Repeat(user.Password, 70), user.FullName)

		s.Empty(gotToken)
		s.Error(err)
	})
}

func (s *AccountTestSuite) TestSignIn() {
	rawPassword := "iVm20%02fD5O"
	user := model.User{
		UUID:      "f9bd9622-f730-11ed-b67e-0242ac120002",
		Login:     "ivan",
		Password:  rawPassword,
		FullName:  "Ivanov Ivan",
		CreatedAt: time.Now(),
	}

	s.Run("Existing user with correct credentials", func() {
		jwtPayload := jwtoken.Payload{
			UUID:     user.UUID,
			FullName: user.FullName,
		}

		wantToken := "eyJhbGci.e30.Et9HFtf9R3GEM"

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		s.Require().NoError(err)

		user.Password = string(hashedPassword)

		s.userRepo.EXPECT().FindByLogin(gomock.Any(), user.Login).Return(&user, nil)
		s.jwt.EXPECT().Create(jwtPayload, jwtTTL, []byte(s.secret)).Return(wantToken, nil)

		gotToken, err := s.accountService.SignIn(context.Background(), user.Login, rawPassword)
		s.Require().NoError(err)

		s.Equal(wantToken, gotToken)
	})

	s.Run("Non-existing user", func() {
		s.userRepo.EXPECT().FindByLogin(gomock.Any(), user.Login).Return(nil, repository.ErrNotFound)
		gotToken, err := s.accountService.SignIn(context.Background(), user.Login, rawPassword)

		s.Empty(gotToken)
		s.ErrorIs(err, ErrUserNotFound)
	})

	s.Run("Non-existing user with unknown error", func() {
		s.userRepo.EXPECT().FindByLogin(gomock.Any(), user.Login).Return(nil, errUnknown)
		gotToken, err := s.accountService.SignIn(context.Background(), user.Login, rawPassword)

		s.Empty(gotToken)
		s.ErrorIs(err, errUnknown)
	})

	s.Run("Existing user with incorrect credentials", func() {
		s.userRepo.EXPECT().FindByLogin(gomock.Any(), user.Login).Return(&user, nil)
		gotToken, err := s.accountService.SignIn(context.Background(), user.Login, rawPassword+"typo")

		s.Empty(gotToken)
		s.ErrorIs(err, ErrIncorrectCredentials)
	})
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}
