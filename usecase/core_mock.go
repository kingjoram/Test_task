// Code generated by MockGen. DO NOT EDIT.
// Source: core.go

// Package usecase is a generated GoMock package.
package usecase

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockICore is a mock of ICore interface.
type MockICore struct {
	ctrl     *gomock.Controller
	recorder *MockICoreMockRecorder
}

// MockICoreMockRecorder is the mock recorder for MockICore.
type MockICoreMockRecorder struct {
	mock *MockICore
}

// NewMockICore creates a new mock instance.
func NewMockICore(ctrl *gomock.Controller) *MockICore {
	mock := &MockICore{ctrl: ctrl}
	mock.recorder = &MockICoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICore) EXPECT() *MockICoreMockRecorder {
	return m.recorder
}

// GetLong mocks base method.
func (m *MockICore) GetLong(short string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLong", short)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLong indicates an expected call of GetLong.
func (mr *MockICoreMockRecorder) GetLong(short interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLong", reflect.TypeOf((*MockICore)(nil).GetLong), short)
}

// GetShort mocks base method.
func (m *MockICore) GetShort(long string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShort", long)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShort indicates an expected call of GetShort.
func (mr *MockICoreMockRecorder) GetShort(long interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShort", reflect.TypeOf((*MockICore)(nil).GetShort), long)
}
