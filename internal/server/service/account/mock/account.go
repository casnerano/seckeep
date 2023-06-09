// Code generated by MockGen. DO NOT EDIT.
// Source: account.go

// Package mock_account is a generated GoMock package.
package mock_account

import (
	reflect "reflect"
	time "time"

	jwtoken "github.com/casnerano/seckeep/pkg/jwtoken"
	gomock "github.com/golang/mock/gomock"
)

// MockJWT is a mock of JWT interface.
type MockJWT struct {
	ctrl     *gomock.Controller
	recorder *MockJWTMockRecorder
}

// MockJWTMockRecorder is the mock recorder for MockJWT.
type MockJWTMockRecorder struct {
	mock *MockJWT
}

// NewMockJWT creates a new mock instance.
func NewMockJWT(ctrl *gomock.Controller) *MockJWT {
	mock := &MockJWT{ctrl: ctrl}
	mock.recorder = &MockJWTMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWT) EXPECT() *MockJWTMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockJWT) Create(payload jwtoken.Payload, ttl time.Duration, secret []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", payload, ttl, secret)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockJWTMockRecorder) Create(payload, ttl, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockJWT)(nil).Create), payload, ttl, secret)
}
