// Code generated by MockGen. DO NOT EDIT.
// Source: desly/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	db "desly/db/sqlc"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateDesly mocks base method.
func (m *MockStore) CreateDesly(arg0 context.Context, arg1 string) (db.Desly, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDesly", arg0, arg1)
	ret0, _ := ret[0].(db.Desly)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDesly indicates an expected call of CreateDesly.
func (mr *MockStoreMockRecorder) CreateDesly(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDesly", reflect.TypeOf((*MockStore)(nil).CreateDesly), arg0, arg1)
}

// CreateDeslyTx mocks base method.
func (m *MockStore) CreateDeslyTx(arg0 context.Context, arg1 db.CreateDeslyTxParams) (db.Desly, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDeslyTx", arg0, arg1)
	ret0, _ := ret[0].(db.Desly)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDeslyTx indicates an expected call of CreateDeslyTx.
func (mr *MockStoreMockRecorder) CreateDeslyTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDeslyTx", reflect.TypeOf((*MockStore)(nil).CreateDeslyTx), arg0, arg1)
}

// GetDesly mocks base method.
func (m *MockStore) GetDesly(arg0 context.Context, arg1 int32) (db.Desly, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDesly", arg0, arg1)
	ret0, _ := ret[0].(db.Desly)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDesly indicates an expected call of GetDesly.
func (mr *MockStoreMockRecorder) GetDesly(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDesly", reflect.TypeOf((*MockStore)(nil).GetDesly), arg0, arg1)
}

// GetDeslyByDesly mocks base method.
func (m *MockStore) GetDeslyByDesly(arg0 context.Context, arg1 string) (db.Desly, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeslyByDesly", arg0, arg1)
	ret0, _ := ret[0].(db.Desly)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeslyByDesly indicates an expected call of GetDeslyByDesly.
func (mr *MockStoreMockRecorder) GetDeslyByDesly(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeslyByDesly", reflect.TypeOf((*MockStore)(nil).GetDeslyByDesly), arg0, arg1)
}