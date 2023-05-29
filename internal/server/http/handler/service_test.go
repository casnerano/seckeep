package handler

import (
	"net/http"
	"testing"

	"github.com/casnerano/seckeep/pkg/log"
	"github.com/stretchr/testify/suite"
)

type ServiceHandlerTestSuite struct {
	suite.Suite
	service *Service
}

func (s *ServiceHandlerTestSuite) SetupSuite() {
	s.service = NewService(log.NewStub())
}

func (s *ServiceHandlerTestSuite) TestPingHandler() {
	s.HTTPStatusCode(s.service.Ping(), http.MethodGet, "/api/ping", nil, http.StatusOK)
}

func TestServiceHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceHandlerTestSuite))
}
