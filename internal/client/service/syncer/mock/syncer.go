// Code generated by MockGen. DO NOT EDIT.
// Source: syncer.go

// Package mock_syncer is a generated GoMock package.
package mock_syncer

import (
	reflect "reflect"

	model "github.com/casnerano/seckeep/internal/client/model"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// GetList mocks base method.
func (m *MockStorage) GetList() []*model.StoreData {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList")
	ret0, _ := ret[0].([]*model.StoreData)
	return ret0
}

// GetList indicates an expected call of GetList.
func (mr *MockStorageMockRecorder) GetList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockStorage)(nil).GetList))
}

// Len mocks base method.
func (m *MockStorage) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len.
func (mr *MockStorageMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockStorage)(nil).Len))
}

// OverwriteStore mocks base method.
func (m *MockStorage) OverwriteStore(memStore []*model.StoreData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OverwriteStore", memStore)
	ret0, _ := ret[0].(error)
	return ret0
}

// OverwriteStore indicates an expected call of OverwriteStore.
func (mr *MockStorageMockRecorder) OverwriteStore(memStore interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OverwriteStore", reflect.TypeOf((*MockStorage)(nil).OverwriteStore), memStore)
}
