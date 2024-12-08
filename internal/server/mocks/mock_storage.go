// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/golovanevvs/metalecoll/internal/server/storage (interfaces: IStorageDB)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/golovanevvs/metalecoll/internal/server/config"
	mapstorage "github.com/golovanevvs/metalecoll/internal/server/mapstorage"
)

// MockIStorageDB is a mock of IStorageDB interface.
type MockIStorageDB struct {
	ctrl     *gomock.Controller
	recorder *MockIStorageDBMockRecorder
}

// MockIStorageDBMockRecorder is the mock recorder for MockIStorageDB.
type MockIStorageDBMockRecorder struct {
	mock *MockIStorageDB
}

// NewMockIStorageDB creates a new mock instance.
func NewMockIStorageDB(ctrl *gomock.Controller) *MockIStorageDB {
	mock := &MockIStorageDB{ctrl: ctrl}
	mock.recorder = &MockIStorageDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIStorageDB) EXPECT() *MockIStorageDBMockRecorder {
	return m.recorder
}

// CloseDB mocks base method.
func (m *MockIStorageDB) CloseDB() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseDB")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseDB indicates an expected call of CloseDB.
func (mr *MockIStorageDBMockRecorder) CloseDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseDB", reflect.TypeOf((*MockIStorageDB)(nil).CloseDB))
}

// GetMetricsFromDB mocks base method.
func (m *MockIStorageDB) GetMetricsFromDB(arg0 context.Context, arg1 *config.Config) (mapstorage.Storage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricsFromDB", arg0, arg1)
	ret0, _ := ret[0].(mapstorage.Storage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetricsFromDB indicates an expected call of GetMetricsFromDB.
func (mr *MockIStorageDBMockRecorder) GetMetricsFromDB(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricsFromDB", reflect.TypeOf((*MockIStorageDB)(nil).GetMetricsFromDB), arg0, arg1)
}

// GetNameDB mocks base method.
func (m *MockIStorageDB) GetNameDB() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNameDB")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetNameDB indicates an expected call of GetNameDB.
func (mr *MockIStorageDBMockRecorder) GetNameDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNameDB", reflect.TypeOf((*MockIStorageDB)(nil).GetNameDB))
}

// Ping mocks base method.
func (m *MockIStorageDB) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockIStorageDBMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockIStorageDB)(nil).Ping))
}

// SaveMetricsToDB mocks base method.
func (m *MockIStorageDB) SaveMetricsToDB(arg0 context.Context, arg1 *config.Config, arg2 mapstorage.Storage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMetricsToDB", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMetricsToDB indicates an expected call of SaveMetricsToDB.
func (mr *MockIStorageDBMockRecorder) SaveMetricsToDB(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMetricsToDB", reflect.TypeOf((*MockIStorageDB)(nil).SaveMetricsToDB), arg0, arg1, arg2)
}