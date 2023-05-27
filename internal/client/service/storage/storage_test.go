package storage

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/casnerano/seckeep/internal/client/model"
	"github.com/casnerano/seckeep/internal/pkg"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	suite.Suite
	storageService  *Storage
	tempStorageFile *os.File
}

func (s *StorageTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	var err error
	s.tempStorageFile, err = os.CreateTemp("", "temp_data.registry")
	s.Require().NoError(err)

	line1 := `{"uuid":"69bb39ff-26e8-44d4-a7f5-d55210a41e0c","type":"TEXT","value":"","version":"2023-05-18T15:41:30.115096Z","created_at":"2023-05-18T15:41:30.115096Z","deleted":false}`
	line2 := `{"uuid":"e68a13dc-b17f-4d2b-839b-c305240827b8","type":"TEXT","value":"","version":"2023-05-18T15:42:37.517163Z","created_at":"2023-05-18T15:41:30.017702Z","deleted":false}`

	_, err = fmt.Fprintf(s.tempStorageFile, "%s\n%s\n", line1, line2)
	s.Require().NoError(err)

	storage, err := New(s.tempStorageFile.Name())
	s.Require().NoError(err)

	s.storageService = storage

	s.Require().Len(s.storageService.memStore, 2)
}

func (s *StorageTestSuite) TearDownSuite() {
	os.Remove(s.tempStorageFile.Name())
}

func (s *StorageTestSuite) TestLen() {
	s.Len(s.storageService.memStore, s.storageService.Len())
}

func (s *StorageTestSuite) TestCreate() {
	dt := model.StoreData{
		UUID:      "8bba5bca-f95f-11ed-be56-0242ac120002",
		Type:      pkg.DataTypeText,
		Value:     []byte{},
		Version:   time.Now(),
		CreatedAt: time.Now(),
		Deleted:   false,
	}
	err := s.storageService.Create(&dt)

	s.NoError(err)

	foundInStore := false
	for key := range s.storageService.memStore {
		if s.storageService.memStore[key].UUID == dt.UUID {
			foundInStore = true
			break
		}
	}

	s.True(foundInStore)
}

func (s *StorageTestSuite) TestRead() {
	s.Run("Good index", func() {
		index := 1
		dt, err := s.storageService.Read(index)

		s.NoError(err)
		s.Equal(s.storageService.memStore[index], dt)
	})

	s.Run("Out of range index", func() {
		index := 100
		dt, err := s.storageService.Read(index)

		s.ErrorIs(err, ErrOutOfRangeStore)
		s.Nil(dt)
	})
}

func (s *StorageTestSuite) TestGetList() {
	result := s.storageService.GetList()
	s.Equal(s.storageService.memStore, result)
}

func (s *StorageTestSuite) TestUpdate() {
	newVersion := time.Now().Add(time.Hour)
	err := s.storageService.Update(0, []byte{}, newVersion)

	s.NoError(err)
	s.Equal(s.storageService.memStore[0].Version, newVersion)
}

func (s *StorageTestSuite) TestDelete() {
	err := s.storageService.Delete(0)
	s.NoError(err)
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
