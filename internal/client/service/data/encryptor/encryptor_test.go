package encryptor

import (
	"errors"
	"testing"

	"github.com/casnerano/seckeep/internal/client/model"
	mock_encryptor "github.com/casnerano/seckeep/internal/client/service/data/encryptor/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

var (
	errUnknown = errors.New("unknown error")
)

type EncryptorTestSuite struct {
	suite.Suite
	cipher    *mock_encryptor.MockCipher
	encryptor *Encryptor
}

func (s *EncryptorTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.cipher = mock_encryptor.NewMockCipher(ctrl)
	s.encryptor = New(s.cipher)
}

func (s *EncryptorTestSuite) TestEncrypt() {
	bJSON := []byte(`{"value":"Example text","meta":null}`)
	bEncrypted := []byte{1, 2, 3, 4, 5}

	dt := model.DataText{
		Value: "Example text",
		Meta:  nil,
	}

	s.Run("Correct data", func() {
		s.cipher.EXPECT().Encrypt(bJSON).Return(bEncrypted, nil)
		gotEncrypted, err := s.encryptor.Encrypt(dt)

		s.Equal(gotEncrypted, bEncrypted)
		s.NoError(err)
	})

	s.Run("Unknown error", func() {
		s.cipher.EXPECT().Encrypt(bJSON).Return(nil, errUnknown)
		gotEncrypted, err := s.encryptor.Encrypt(dt)

		s.Nil(gotEncrypted)
		s.ErrorIs(err, errUnknown)
	})
}

func (s *EncryptorTestSuite) TestDecrypt() {
	bJSON := []byte(`{"value":"Example text","meta":null}`)
	bEncrypted := []byte{1, 2, 3, 4, 5}

	dt := model.DataText{
		Value: "Example text",
		Meta:  nil,
	}

	s.Run("Correct data", func() {
		s.cipher.EXPECT().Decrypt(bEncrypted).Return(bJSON, nil)
		dtResult := model.DataText{}
		err := s.encryptor.Decrypt(bEncrypted, &dtResult)

		s.Equal(dt, dtResult)
		s.NoError(err)
	})

	s.Run("Unknown error", func() {
		s.cipher.EXPECT().Decrypt(bEncrypted).Return(nil, errUnknown)
		dtResult := model.DataText{}
		err := s.encryptor.Decrypt(bEncrypted, &dtResult)

		s.Equal(model.DataText{}, dtResult)
		s.ErrorIs(err, errUnknown)
	})

	s.Run("Unmarshal error", func() {
		s.cipher.EXPECT().Decrypt(bEncrypted).Return([]byte{1}, nil)
		dtResult := model.DataText{}
		err := s.encryptor.Decrypt(bEncrypted, &dtResult)

		s.Equal(model.DataText{}, dtResult)
		s.Error(err)
	})
}

func TestEncryptorTestSuite(t *testing.T) {
	suite.Run(t, new(EncryptorTestSuite))
}
