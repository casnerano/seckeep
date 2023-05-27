package data

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/casnerano/seckeep/internal/pkg/model"
	"github.com/casnerano/seckeep/internal/server/repository"
	mock_repository "github.com/casnerano/seckeep/internal/server/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

var (
	errUnknown = errors.New("unknown error")
)

type DataTestSuite struct {
	suite.Suite
	dataService *Data
	dataRepo    *mock_repository.MockData
}

func (s *DataTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.dataRepo = mock_repository.NewMockData(ctrl)
	s.dataService = New(s.dataRepo)
}

func (s *DataTestSuite) TestCreate() {
	wantData := model.Data{
		UUID:      "f9bd9622-f730-11ed-b67e-0242ac120000",
		UserUUID:  "f9bd9622-f730-11ed-b67e-0242ac120002",
		Type:      model.DataTypeText,
		Value:     []byte(""),
		Version:   time.Now(),
		CreatedAt: time.Now(),
	}

	s.dataRepo.EXPECT().Add(gomock.Any(), wantData).Return(&wantData, nil)
	gotData, err := s.dataService.Create(context.Background(), wantData)

	s.NoError(err)
	s.Equal(wantData, *gotData)
}

func (s *DataTestSuite) TestFindByUUID() {
	wantData := model.Data{
		UUID:      "f9bd9622-f730-11ed-b67e-0242ac120000",
		UserUUID:  "f9bd9622-f730-11ed-b67e-0242ac120002",
		Type:      model.DataTypeText,
		Value:     []byte(""),
		Version:   time.Now(),
		CreatedAt: time.Now(),
	}

	s.Run("Data is exist", func() {
		s.dataRepo.EXPECT().FindByUUID(gomock.Any(), wantData.UserUUID, wantData.UUID).Return(&wantData, nil)
		gotData, err := s.dataService.FindByUUID(context.Background(), wantData.UserUUID, wantData.UUID)

		s.NoError(err)
		s.Equal(wantData, *gotData)
	})

	s.Run("Data is not exist", func() {
		s.dataRepo.EXPECT().FindByUUID(gomock.Any(), wantData.UserUUID, wantData.UUID).Return(nil, repository.ErrNotFound)
		gotData, err := s.dataService.FindByUUID(context.Background(), wantData.UserUUID, wantData.UUID)

		s.Nil(gotData)
		s.ErrorIs(err, ErrNotFound)
	})

	s.Run("Unknown error", func() {
		s.dataRepo.EXPECT().FindByUUID(gomock.Any(), wantData.UserUUID, wantData.UUID).Return(nil, errUnknown)
		gotData, err := s.dataService.FindByUUID(context.Background(), wantData.UserUUID, wantData.UUID)

		s.Nil(gotData)
		s.ErrorIs(err, errUnknown)
	})
}

func (s *DataTestSuite) TestFindByUserUUID() {
	userUUID := "f9bd9622-f730-11ed-b67e-0242ac000000"
	wantDataList := []*model.Data{
		{
			UUID:      "f9bd9622-f730-11ed-b67e-0242ac120000",
			UserUUID:  "f9bd9622-f730-11ed-b67e-0242ac120002",
			Type:      model.DataTypeText,
			Value:     []byte("1"),
			Version:   time.Now(),
			CreatedAt: time.Now(),
		},
		{
			UUID:      "f9bd9622-f730-11ed-b67e-0242ac130000",
			UserUUID:  "f9bd9622-f730-11ed-b67e-0242ac130002",
			Type:      model.DataTypeCard,
			Value:     []byte("2"),
			Version:   time.Now(),
			CreatedAt: time.Now(),
		},
	}

	s.dataRepo.EXPECT().FindByUserUUID(gomock.Any(), userUUID).Return(wantDataList, nil)
	gotDataList, err := s.dataService.FindByUserUUID(context.Background(), userUUID)

	s.NoError(err)
	s.Equal(wantDataList, gotDataList)
}

func (s *DataTestSuite) TestUpdate() {
	wantData := model.Data{
		UUID:      "f9bd9622-f730-11ed-b67e-0242ac120000",
		UserUUID:  "f9bd9622-f730-11ed-b67e-0242ac120002",
		Type:      model.DataTypeText,
		Value:     []byte(""),
		Version:   time.Now(),
		CreatedAt: time.Now(),
	}

	s.Run("Data is exist", func() {
		s.dataRepo.EXPECT().Update(gomock.Any(), wantData.UserUUID, wantData.UUID, wantData.Value, wantData.Version).Return(&wantData, nil)
		gotData, err := s.dataService.Update(context.Background(), wantData.UserUUID, wantData.UUID, wantData.Value, wantData.Version)

		s.NoError(err)
		s.Equal(wantData, *gotData)
	})

	s.Run("Data is not exist", func() {
		s.dataRepo.EXPECT().Update(gomock.Any(), wantData.UserUUID, wantData.UUID, wantData.Value, wantData.Version).Return(nil, repository.ErrNotFound)
		gotData, err := s.dataService.Update(context.Background(), wantData.UserUUID, wantData.UUID, wantData.Value, wantData.Version)

		s.Nil(gotData)
		s.ErrorIs(err, ErrNotFound)
	})

	s.Run("Unknown error", func() {
		s.dataRepo.EXPECT().Update(gomock.Any(), wantData.UserUUID, wantData.UUID, wantData.Value, wantData.Version).Return(nil, errUnknown)
		gotData, err := s.dataService.Update(context.Background(), wantData.UserUUID, wantData.UUID, wantData.Value, wantData.Version)

		s.Nil(gotData)
		s.ErrorIs(err, errUnknown)
	})
}

func (s *DataTestSuite) TestDelete() {
	uuid := "f9bd9622-f730-11ed-b67e-0242ac120000"
	userUUID := "f9bd9622-f730-11ed-b67e-0242ac120002"

	s.Run("Data is exist", func() {
		s.dataRepo.EXPECT().Delete(gomock.Any(), userUUID, uuid).Return(nil)
		err := s.dataService.Delete(context.Background(), userUUID, uuid)

		s.NoError(err)
	})

	s.Run("Data is not exist", func() {
		s.dataRepo.EXPECT().Delete(gomock.Any(), userUUID, uuid).Return(repository.ErrNotFound)
		err := s.dataService.Delete(context.Background(), userUUID, uuid)

		s.ErrorIs(err, ErrNotFound)
	})

	s.Run("Unknown error", func() {
		s.dataRepo.EXPECT().Delete(gomock.Any(), userUUID, uuid).Return(errUnknown)
		err := s.dataService.Delete(context.Background(), userUUID, uuid)

		s.ErrorIs(err, errUnknown)
	})
}

func TestDataTestSuite(t *testing.T) {
	suite.Run(t, new(DataTestSuite))
}
