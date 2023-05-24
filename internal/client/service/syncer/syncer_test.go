package syncer

import (
	"errors"
	"net/http"
	"testing"

	"github.com/casnerano/seckeep/internal/client/model"
	mock_syncer "github.com/casnerano/seckeep/internal/client/service/syncer/mock"
	"github.com/casnerano/seckeep/pkg/log"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

var (
	errUnknown = errors.New("unknown error")
)

type DataTestSuite struct {
	suite.Suite
	client        *resty.Client
	storage       *mock_syncer.MockStorage
	syncerService *Syncer
}

func (s *DataTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.storage = mock_syncer.NewMockStorage(ctrl)

	s.client = resty.New()
	s.client.SetBaseURL("http://127.0.0.1/api")

	s.syncerService = New(s.client, s.storage, log.NewStub())
}

func (s *DataTestSuite) SetupTest() {
	httpmock.ActivateNonDefault(s.client.GetClient())
}

func (s *DataTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (s *DataTestSuite) SetupSubTest() {
	httpmock.Reset()
}

func (s *DataTestSuite) TestPingServerHealth() {
	s.Run("Error response", func() {
		httpmock.RegisterResponder(
			http.MethodGet, s.client.BaseURL+"/ping",
			httpmock.NewErrorResponder(errUnknown),
		)

		err := s.syncerService.PingServerHealth()

		s.ErrorIs(err, errUnknown)
	})

	s.Run("Error response", func() {
		httpmock.RegisterResponder(
			http.MethodGet, s.client.BaseURL+"/ping",
			httpmock.NewStringResponder(http.StatusUnauthorized, ""),
		)

		err := s.syncerService.PingServerHealth()

		s.ErrorIs(err, ErrUnauthorized)
	})

	s.Run("Unknown response status", func() {
		httpmock.RegisterResponder(
			http.MethodGet, s.client.BaseURL+"/ping",
			httpmock.NewStringResponder(http.StatusTeapot, ""),
		)

		err := s.syncerService.PingServerHealth()

		s.Error(err)
	})

	s.Run("Good response", func() {
		httpmock.RegisterResponder(
			http.MethodGet, s.client.BaseURL+"/ping",
			httpmock.NewStringResponder(http.StatusOK, ""),
		)

		err := s.syncerService.PingServerHealth()

		s.NoError(err)
	})
}

func (s *DataTestSuite) TestRun() {
	s.Run("Good run sync", func() {
		s.storage.EXPECT().Len().Return(0)

		responder, err := httpmock.NewJsonResponder(http.StatusOK, []struct{}{})
		s.Require().NoError(err)

		httpmock.RegisterResponder(
			http.MethodGet, s.client.BaseURL+"/data",
			responder,
		)

		s.storage.EXPECT().GetList().Return([]*model.StoreData{
			{UUID: ""},
			{UUID: "u2", Deleted: true},
			{UUID: "u3"},
		})

		httpmock.RegisterResponder(
			http.MethodPost, s.client.BaseURL+"/data",
			httpmock.NewStringResponder(http.StatusOK, ""),
		)

		httpmock.RegisterResponder(
			http.MethodDelete, "=~^"+s.client.BaseURL+"/data/\\w+",
			httpmock.NewStringResponder(http.StatusOK, ""),
		)

		httpmock.RegisterResponder(
			http.MethodPut, "=~^"+s.client.BaseURL+"/data/\\w+",
			httpmock.NewStringResponder(http.StatusOK, ""),
		)

		httpmock.RegisterResponder(
			http.MethodGet, s.client.BaseURL+"/data",
			responder,
		)

		s.storage.EXPECT().OverwriteStore([]*model.StoreData{}).Return(nil)

		err = s.syncerService.Run()

		s.NoError(err)
	})
}

func (s *DataTestSuite) TestLoadToServer() {
	s.Run("Empty items", func() {
		httpmock.RegisterResponder(
			http.MethodPost, s.client.BaseURL+"/data",
			httpmock.NewStringResponder(http.StatusOK, ""),
		)

		err := s.syncerService.loadToServer([]*model.StoreData{})

		s.NoError(err)
	})

	s.Run("Good data items", func() {
		httpmock.RegisterResponder(
			http.MethodPost, s.client.BaseURL+"/data",
			httpmock.NewStringResponder(http.StatusOK, ""),
		)

		err := s.syncerService.loadToServer([]*model.StoreData{
			{UUID: "u1"},
			{UUID: "u2"},
		})

		s.NoError(err)
	})

	s.Run("Error response", func() {
		httpmock.RegisterResponder(
			http.MethodPost, s.client.BaseURL+"/data",
			httpmock.NewErrorResponder(errUnknown),
		)

		err := s.syncerService.loadToServer([]*model.StoreData{
			{UUID: "u1"},
		})

		s.ErrorIs(err, errUnknown)
	})
}

func (s *DataTestSuite) TestUpdateToServer() {
	s.Run("Empty items", func() {
		httpmock.RegisterResponder(
			http.MethodPut, "=~^"+s.client.BaseURL+"/data/\\w+",
			httpmock.NewStringResponder(http.StatusOK, ""),
		)

		err := s.syncerService.updateToServer([]*model.StoreData{})

		s.NoError(err)
	})

	s.Run("Good data items", func() {
		httpmock.RegisterResponder(
			http.MethodPut, "=~^"+s.client.BaseURL+"/data/\\w+",
			httpmock.NewStringResponder(http.StatusOK, ""),
		)

		err := s.syncerService.updateToServer([]*model.StoreData{
			{UUID: "u1"},
			{UUID: "u2"},
		})

		s.NoError(err)
	})

	s.Run("Error response", func() {
		httpmock.RegisterResponder(
			http.MethodPut, "=~^"+s.client.BaseURL+"/data/\\w+",
			httpmock.NewErrorResponder(errUnknown),
		)

		err := s.syncerService.updateToServer([]*model.StoreData{
			{UUID: "u1"},
		})

		s.ErrorIs(err, errUnknown)
	})
}

func (s *DataTestSuite) TestExportFromServer() {
	s.Run("Empty items", func() {
		responder, err := httpmock.NewJsonResponder(http.StatusOK, []struct{}{})
		s.Require().NoError(err)
		httpmock.RegisterResponder(
			http.MethodGet, s.client.BaseURL+"/data",
			responder,
		)

		result, err := s.syncerService.exportFromServer()

		s.Empty(result)
		s.NoError(err)
	})

	s.Run("Good data items", func() {
		responder, err := httpmock.NewJsonResponder(
			http.StatusOK,
			[]*storeData{
				{data: &model.StoreData{UUID: "u1"}},
				{data: &model.StoreData{UUID: "u2"}},
			},
		)
		s.Require().NoError(err)
		httpmock.RegisterResponder(
			http.MethodGet, s.client.BaseURL+"/data",
			responder,
		)

		result, err := s.syncerService.exportFromServer()

		s.Len(result, 2)
		s.NoError(err)
	})

	s.Run("Error response", func() {
		httpmock.RegisterResponder(
			http.MethodGet, s.client.BaseURL+"/data",
			httpmock.NewErrorResponder(errUnknown),
		)

		result, err := s.syncerService.exportFromServer()

		s.Nil(result)
		s.ErrorIs(err, errUnknown)
	})
}

func (s *DataTestSuite) TestRemoveFromServer() {
	s.Run("Empty request", func() {
		err := s.syncerService.removeFromServer([]*storeData{})
		s.NoError(err)
	})

	s.Run("Good data items", func() {
		httpmock.RegisterResponder(
			http.MethodDelete, "=~^"+s.client.BaseURL+"/data/\\w+",
			httpmock.NewStringResponder(http.StatusOK, ""),
		)

		err := s.syncerService.removeFromServer([]*storeData{
			{data: &model.StoreData{UUID: "u1"}},
			{data: &model.StoreData{UUID: "u2"}},
		})

		s.NoError(err)
	})

	s.Run("Error response", func() {
		httpmock.RegisterResponder(
			http.MethodDelete, "=~^"+s.client.BaseURL+"/data/\\w+",
			httpmock.NewErrorResponder(errUnknown),
		)

		err := s.syncerService.removeFromServer([]*storeData{
			{data: &model.StoreData{UUID: "u1"}},
		})

		s.ErrorIs(err, errUnknown)
	})
}

func TestDataTestSuite(t *testing.T) {
	suite.Run(t, new(DataTestSuite))
}
