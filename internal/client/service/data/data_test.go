package data

import (
	"errors"
	"testing"

	"github.com/casnerano/seckeep/internal/client/model"
	mock_data "github.com/casnerano/seckeep/internal/client/service/data/mock"
	smodel "github.com/casnerano/seckeep/internal/pkg/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

var (
	errUnknown = errors.New("unknown error")
)

type DataTestSuite struct {
	suite.Suite
	storage    *mock_data.MockStorage
	encryptor  *mock_data.MockEncryptor
	dataSerice *Data
}

func (s *DataTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.storage = mock_data.NewMockStorage(ctrl)
	s.encryptor = mock_data.NewMockEncryptor(ctrl)
	s.dataSerice = New(s.storage, s.encryptor)
}

func (s *DataTestSuite) TestCreate() {
	textDt := model.DataText{
		Value: "Example text",
		Meta:  nil,
	}

	s.Run("Correct data", func() {
		s.encryptor.EXPECT().Encrypt(textDt).Return([]byte{1, 2, 3, 4, 5}, nil)
		s.storage.EXPECT().Create(gomock.Any()).Return(nil)
		err := s.dataSerice.Create(textDt)

		s.NoError(err)
	})

	s.Run("Encryptor unknown error", func() {
		s.encryptor.EXPECT().Encrypt(textDt).Return(nil, errUnknown)
		err := s.dataSerice.Create(textDt)

		s.ErrorIs(err, errUnknown)
	})

	s.Run("Unknown storage error", func() {
		s.encryptor.EXPECT().Encrypt(textDt).Return([]byte{1, 2, 3, 4, 5}, nil)
		s.storage.EXPECT().Create(gomock.Any()).Return(errUnknown)
		err := s.dataSerice.Create(textDt)

		s.ErrorIs(err, errUnknown)
	})
}

func (s *DataTestSuite) TestRead() {
	index := 1
	storeData := model.StoreData{
		UUID:  "c9d5577e-f8cf-11ed-be56-0242ac120002",
		Type:  smodel.DataTypeText,
		Value: []byte("Example text"),
	}

	dataTypeList := []model.DataTypeable{
		&model.DataCredential{},
		&model.DataText{},
		&model.DataCard{},
		&model.DataDocument{},
	}

	for key, valType := range dataTypeList {
		s.Run("Correct data", func() {
			storeData.Type = valType.Type()
			s.storage.EXPECT().Read(index).Return(&storeData, nil)
			s.encryptor.EXPECT().Decrypt(storeData.Value, dataTypeList[key]).Return(nil)
			gotDt, err := s.dataSerice.Read(index)

			s.Equal(dataTypeList[key], gotDt)
			s.NoError(err)
		})
	}

	s.Run("Unknown storage DataType", func() {
		s.storage.EXPECT().Read(index).Return(&model.StoreData{Type: "unknown"}, nil)
		gotDt, err := s.dataSerice.Read(index)

		s.Nil(gotDt)
		s.Error(err)
	})

	s.Run("Incorrect storage read", func() {
		s.storage.EXPECT().Read(index).Return(nil, errUnknown)
		gotDt, err := s.dataSerice.Read(index)

		s.Nil(gotDt)
		s.ErrorIs(err, errUnknown)
	})

	s.Run("Incorrect decrypt", func() {
		storeData.Type = smodel.DataTypeText
		s.storage.EXPECT().Read(index).Return(&storeData, nil)
		s.encryptor.EXPECT().Decrypt(storeData.Value, &model.DataText{}).Return(errUnknown)
		gotDt, err := s.dataSerice.Read(index)

		s.Nil(gotDt)
		s.ErrorIs(err, errUnknown)
	})
}

func (s *DataTestSuite) TestGetList() {
	dtItems := []*model.StoreData{
		{Deleted: true},
		{Deleted: false},
	}

	s.storage.EXPECT().GetList().Return(dtItems)
	s.storage.EXPECT().Read(gomock.Any()).Return(&model.StoreData{Type: smodel.DataTypeText}, nil)
	s.encryptor.EXPECT().Decrypt(gomock.Any(), &model.DataText{}).Return(nil)
	result := s.dataSerice.GetList()

	s.NotEmpty(result)
}

func (s *DataTestSuite) TestUpdate() {
	textDt := model.DataText{
		Value: "Example text",
		Meta:  nil,
	}

	s.Run("Correct data", func() {
		encrypted := []byte{1, 2, 3, 4, 5}
		s.encryptor.EXPECT().Encrypt(textDt).Return(encrypted, nil)
		s.storage.EXPECT().Update(1, encrypted, gomock.Any()).Return(nil)
		err := s.dataSerice.Update(1, textDt)

		s.NoError(err)
	})

	s.Run("Encryptor unknown error", func() {
		s.encryptor.EXPECT().Encrypt(textDt).Return(nil, errUnknown)
		err := s.dataSerice.Update(1, textDt)

		s.ErrorIs(err, errUnknown)
	})
}

func (s *DataTestSuite) TestDelete() {
	s.storage.EXPECT().Delete(1).Return(nil)
	err := s.dataSerice.Delete(1)

	s.NoError(err)
}

func TestDataTestSuite(t *testing.T) {
	suite.Run(t, new(DataTestSuite))
}
