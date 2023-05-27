package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/casnerano/seckeep/internal/pkg"
	mock_handler "github.com/casnerano/seckeep/internal/server/http/handler/mock"
	"github.com/casnerano/seckeep/internal/server/http/middleware"
	"github.com/casnerano/seckeep/internal/server/model"
	dataService "github.com/casnerano/seckeep/internal/server/service/data"
	"github.com/casnerano/seckeep/pkg/log"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type DataHandlerTestSuite struct {
	suite.Suite
	handler     *Data
	dataService *mock_handler.MockDataService
}

func (s *DataHandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.dataService = mock_handler.NewMockDataService(ctrl)

	s.handler = NewData(s.dataService, log.NewStub())
}

func (s *DataHandlerTestSuite) TestCreateHandler() {
	userUUID := "ba3cfc2c-f7fd-11ed-b67e-0242ac120002"

	rd := model.DataCreateRequest{
		Type:      pkg.DataTypeText,
		Value:     []byte(""),
		Version:   time.Now(),
		CreatedAt: time.Now(),
	}

	data := pkg.Data{
		UserUUID:  userUUID,
		Type:      rd.Type,
		Value:     rd.Value,
		Version:   rd.Version,
		CreatedAt: rd.CreatedAt,
	}

	resultData := data
	resultData.UUID = "9b92672a-f7fe-11ed-b67e-0242ac120002"

	request := httptest.NewRequest(http.MethodPost, "/api/data", nil)
	ctx := context.WithValue(request.Context(), middleware.CtxUserUUIDKey, userUUID)
	requestWithUserUUIDCtx := request.WithContext(ctx)

	s.Run("Correct data with user uuid", func() {
		s.dataService.EXPECT().Create(gomock.Any(), data).Return(&resultData, nil)

		result, status := s.handler.Create(rd, httptest.NewRecorder(), requestWithUserUUIDCtx)

		s.Equal(http.StatusOK, status)
		s.Equal(result, &resultData)
	})

	s.Run("Without user uuid", func() {
		result, status := s.handler.Create(rd, httptest.NewRecorder(), request)
		s.Nil(result)
		s.Equal(http.StatusUnauthorized, status)
	})

	s.Run("Has unknown error", func() {
		s.dataService.EXPECT().Create(gomock.Any(), data).Return(nil, errors.New("unknown error"))

		result, status := s.handler.Create(rd, httptest.NewRecorder(), requestWithUserUUIDCtx)
		s.Require().Nil(result)
		s.Equal(http.StatusInternalServerError, status)
	})
}

func (s *DataHandlerTestSuite) TestUpdateHandler() {
	uuid := "9b92672a-f7fe-11ed-b67e-0242ac120002"
	userUUID := "ba3cfc2c-f7fd-11ed-b67e-0242ac120002"

	rd := model.DataUpdateRequest{
		Value:   []byte(""),
		Version: time.Now(),
	}

	data := pkg.Data{
		UserUUID:  userUUID,
		Type:      pkg.DataTypeText,
		Value:     rd.Value,
		Version:   rd.Version,
		CreatedAt: time.Now(),
	}

	request := httptest.NewRequest(http.MethodPut, "/api/data/"+uuid, nil)
	ctx := context.WithValue(request.Context(), middleware.CtxUserUUIDKey, userUUID)
	requestWithUserUUIDCtx := request.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("uuid", uuid)

	ctx = context.WithValue(requestWithUserUUIDCtx.Context(), chi.RouteCtxKey, chiCtx)
	requestWithDataAndUserCtx := request.WithContext(ctx)

	s.Run("Correct data with user uuid", func() {
		s.dataService.EXPECT().Update(gomock.Any(), userUUID, uuid, rd.Value, rd.Version).Return(&data, nil)
		result, status := s.handler.Update(rd, httptest.NewRecorder(), requestWithDataAndUserCtx)

		s.Equal(result, &data)
		s.Equal(http.StatusOK, status)
	})

	s.Run("Without user uuid", func() {
		result, status := s.handler.Update(rd, httptest.NewRecorder(), request)

		s.Nil(result)
		s.Equal(http.StatusUnauthorized, status)
	})

	s.Run("Has unknown error", func() {
		s.dataService.EXPECT().Update(gomock.Any(), userUUID, uuid, rd.Value, rd.Version).Return(nil, errors.New("unknown error"))
		result, status := s.handler.Update(rd, httptest.NewRecorder(), requestWithDataAndUserCtx)

		s.Nil(result)
		s.Equal(http.StatusInternalServerError, status)
	})

	s.Run("Without data uuid", func() {
		result, status := s.handler.Update(rd, httptest.NewRecorder(), requestWithUserUUIDCtx)

		s.Nil(result)
		s.Equal(http.StatusBadRequest, status)
	})
}

func (s *DataHandlerTestSuite) TestGetHandler() {
	uuid := "9b92672a-f7fe-11ed-b67e-0242ac120002"
	userUUID := "ba3cfc2c-f7fd-11ed-b67e-0242ac120002"

	data := pkg.Data{
		UserUUID:  userUUID,
		Type:      pkg.DataTypeText,
		Value:     []byte(""),
		Version:   time.Now(),
		CreatedAt: time.Now(),
	}

	request := httptest.NewRequest(http.MethodGet, "/api/data/"+uuid, nil)
	ctx := context.WithValue(request.Context(), middleware.CtxUserUUIDKey, userUUID)
	requestWithUserUUIDCtx := request.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("uuid", uuid)

	ctx = context.WithValue(requestWithUserUUIDCtx.Context(), chi.RouteCtxKey, chiCtx)
	requestWithDataAndUserCtx := request.WithContext(ctx)

	s.Run("Existing data with user uuid", func() {
		s.dataService.EXPECT().FindByUUID(gomock.Any(), userUUID, uuid).Return(&data, nil)
		w := httptest.NewRecorder()
		result, status := s.handler.Get(w, requestWithDataAndUserCtx)

		s.Equal(result, &data)
		s.Equal(http.StatusOK, status)
	})

	s.Run("Without user uuid", func() {
		result, status := s.handler.Get(httptest.NewRecorder(), request)

		s.Nil(result)
		s.Equal(http.StatusUnauthorized, status)
	})

	s.Run("Has unknown error", func() {
		s.dataService.EXPECT().FindByUUID(gomock.Any(), userUUID, uuid).Return(nil, errors.New("unknown error"))
		result, status := s.handler.Get(httptest.NewRecorder(), requestWithDataAndUserCtx)

		s.Nil(result)
		s.Equal(http.StatusInternalServerError, status)
	})

	s.Run("Without data uuid", func() {
		result, status := s.handler.Get(httptest.NewRecorder(), requestWithUserUUIDCtx)

		s.Nil(result)
		s.Equal(http.StatusBadRequest, status)
	})

	s.Run("Non-existing data", func() {
		s.dataService.EXPECT().FindByUUID(gomock.Any(), userUUID, uuid).Return(nil, dataService.ErrNotFound)
		result, status := s.handler.Get(httptest.NewRecorder(), requestWithDataAndUserCtx)

		s.Nil(result)
		s.Equal(http.StatusNotFound, status)
	})
}

func (s *DataHandlerTestSuite) TestGetListHandler() {
	userUUID := "ba3cfc2c-f7fd-11ed-b67e-0242ac120002"

	dataList := []*pkg.Data{
		{
			UserUUID:  userUUID,
			Type:      pkg.DataTypeText,
			Value:     []byte("1"),
			Version:   time.Now(),
			CreatedAt: time.Now(),
		},
		{
			UserUUID:  userUUID,
			Type:      pkg.DataTypeCard,
			Value:     []byte("2"),
			Version:   time.Now(),
			CreatedAt: time.Now(),
		},
	}

	request := httptest.NewRequest(http.MethodGet, "/api/data", nil)
	ctx := context.WithValue(request.Context(), middleware.CtxUserUUIDKey, userUUID)
	requestWithUserUUIDCtx := request.WithContext(ctx)

	s.Run("Existing data with user uuid", func() {
		s.dataService.EXPECT().FindByUserUUID(gomock.Any(), userUUID).Return(dataList, nil)
		w := httptest.NewRecorder()
		result, status := s.handler.GetList(w, requestWithUserUUIDCtx)

		s.Equal(result, dataList)
		s.Equal(http.StatusOK, status)
	})

	s.Run("Without user uuid", func() {
		result, status := s.handler.GetList(httptest.NewRecorder(), request)

		s.Nil(result)
		s.Equal(http.StatusUnauthorized, status)
	})

	s.Run("Has unknown error", func() {
		s.dataService.EXPECT().FindByUserUUID(gomock.Any(), userUUID).Return(nil, errors.New("unknown error"))
		result, status := s.handler.GetList(httptest.NewRecorder(), requestWithUserUUIDCtx)

		s.Nil(result)
		s.Equal(http.StatusInternalServerError, status)
	})
}

func (s *DataHandlerTestSuite) TestDeleteHandler() {
	uuid := "9b92672a-f7fe-11ed-b67e-0242ac120002"
	userUUID := "ba3cfc2c-f7fd-11ed-b67e-0242ac120002"

	request := httptest.NewRequest(http.MethodDelete, "/api/data/"+uuid, nil)
	ctx := context.WithValue(request.Context(), middleware.CtxUserUUIDKey, userUUID)
	requestWithUserUUIDCtx := request.WithContext(ctx)

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("uuid", uuid)

	ctx = context.WithValue(requestWithUserUUIDCtx.Context(), chi.RouteCtxKey, chiCtx)
	requestWithDataAndUserCtx := request.WithContext(ctx)

	s.Run("Correct data with user uuid", func() {
		s.dataService.EXPECT().Delete(gomock.Any(), userUUID, uuid).Return(nil)
		result, status := s.handler.Delete(httptest.NewRecorder(), requestWithDataAndUserCtx)

		s.Nil(result)
		s.Equal(http.StatusOK, status)
	})

	s.Run("Without user uuid", func() {
		result, status := s.handler.Delete(httptest.NewRecorder(), request)

		s.Nil(result)
		s.Equal(http.StatusUnauthorized, status)
	})

	s.Run("Has unknown error", func() {
		s.dataService.EXPECT().Delete(gomock.Any(), userUUID, uuid).Return(errors.New("unknown error"))
		result, status := s.handler.Delete(httptest.NewRecorder(), requestWithDataAndUserCtx)

		s.Nil(result)
		s.Equal(http.StatusInternalServerError, status)
	})

	s.Run("Without data uuid", func() {
		result, status := s.handler.Delete(httptest.NewRecorder(), requestWithUserUUIDCtx)

		s.Nil(result)
		s.Equal(http.StatusBadRequest, status)
	})

	s.Run("Non-existing data", func() {
		s.dataService.EXPECT().Delete(gomock.Any(), userUUID, uuid).Return(dataService.ErrNotFound)
		result, status := s.handler.Delete(httptest.NewRecorder(), requestWithDataAndUserCtx)

		s.Nil(result)
		s.Equal(http.StatusNotFound, status)
	})
}

func TestDataHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DataHandlerTestSuite))
}
