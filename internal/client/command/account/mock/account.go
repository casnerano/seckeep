// Code generated by MockGen. DO NOT EDIT.
// Source: account.go

// Package mock_account is a generated GoMock package.
package mock_account

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAccountService is a mock of AccountService interface.
type MockAccountService struct {
	ctrl     *gomock.Controller
	recorder *MockAccountServiceMockRecorder
}

// MockAccountServiceMockRecorder is the mock recorder for MockAccountService.
type MockAccountServiceMockRecorder struct {
	mock *MockAccountService
}

// NewMockAccountService creates a new mock instance.
func NewMockAccountService(ctrl *gomock.Controller) *MockAccountService {
	mock := &MockAccountService{ctrl: ctrl}
	mock.recorder = &MockAccountServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountService) EXPECT() *MockAccountServiceMockRecorder {
	return m.recorder
}

// SignIn mocks base method.
func (m *MockAccountService) SignIn(login, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", login, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAccountServiceMockRecorder) SignIn(login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAccountService)(nil).SignIn), login, password)
}

// SignUp mocks base method.
func (m *MockAccountService) SignUp(login, password, fullName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", login, password, fullName)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAccountServiceMockRecorder) SignUp(login, password, fullName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAccountService)(nil).SignUp), login, password, fullName)
}