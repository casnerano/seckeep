// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"
	time "time"

	model "github.com/casnerano/seckeep/internal/pkg/model"
	model0 "github.com/casnerano/seckeep/internal/server/model"
	gomock "github.com/golang/mock/gomock"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockUser) Add(ctx context.Context, login, password, fullName string) (*model0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, login, password, fullName)
	ret0, _ := ret[0].(*model0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockUserMockRecorder) Add(ctx, login, password, fullName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockUser)(nil).Add), ctx, login, password, fullName)
}

// FindByLogin mocks base method.
func (m *MockUser) FindByLogin(ctx context.Context, login string) (*model0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByLogin", ctx, login)
	ret0, _ := ret[0].(*model0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByLogin indicates an expected call of FindByLogin.
func (mr *MockUserMockRecorder) FindByLogin(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByLogin", reflect.TypeOf((*MockUser)(nil).FindByLogin), ctx, login)
}

// FindByUUID mocks base method.
func (m *MockUser) FindByUUID(ctx context.Context, uuid string) (*model0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUUID", ctx, uuid)
	ret0, _ := ret[0].(*model0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUUID indicates an expected call of FindByUUID.
func (mr *MockUserMockRecorder) FindByUUID(ctx, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUUID", reflect.TypeOf((*MockUser)(nil).FindByUUID), ctx, uuid)
}

// MockData is a mock of Data interface.
type MockData struct {
	ctrl     *gomock.Controller
	recorder *MockDataMockRecorder
}

// MockDataMockRecorder is the mock recorder for MockData.
type MockDataMockRecorder struct {
	mock *MockData
}

// NewMockData creates a new mock instance.
func NewMockData(ctrl *gomock.Controller) *MockData {
	mock := &MockData{ctrl: ctrl}
	mock.recorder = &MockDataMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockData) EXPECT() *MockDataMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockData) Add(ctx context.Context, data model.Data) (*model.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, data)
	ret0, _ := ret[0].(*model.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockDataMockRecorder) Add(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockData)(nil).Add), ctx, data)
}

// Delete mocks base method.
func (m *MockData) Delete(ctx context.Context, userUUID, uuid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, userUUID, uuid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDataMockRecorder) Delete(ctx, userUUID, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockData)(nil).Delete), ctx, userUUID, uuid)
}

// FindByUUID mocks base method.
func (m *MockData) FindByUUID(ctx context.Context, userUUID, uuid string) (*model.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUUID", ctx, userUUID, uuid)
	ret0, _ := ret[0].(*model.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUUID indicates an expected call of FindByUUID.
func (mr *MockDataMockRecorder) FindByUUID(ctx, userUUID, uuid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUUID", reflect.TypeOf((*MockData)(nil).FindByUUID), ctx, userUUID, uuid)
}

// FindByUserUUID mocks base method.
func (m *MockData) FindByUserUUID(ctx context.Context, userUUID string) ([]*model.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserUUID", ctx, userUUID)
	ret0, _ := ret[0].([]*model.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserUUID indicates an expected call of FindByUserUUID.
func (mr *MockDataMockRecorder) FindByUserUUID(ctx, userUUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserUUID", reflect.TypeOf((*MockData)(nil).FindByUserUUID), ctx, userUUID)
}

// Update mocks base method.
func (m *MockData) Update(ctx context.Context, userUUID, uuid string, value []byte, version time.Time) (*model.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, userUUID, uuid, value, version)
	ret0, _ := ret[0].(*model.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockDataMockRecorder) Update(ctx, userUUID, uuid, value, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockData)(nil).Update), ctx, userUUID, uuid, value, version)
}
