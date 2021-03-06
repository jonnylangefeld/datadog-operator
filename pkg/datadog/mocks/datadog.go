// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jonnylangefeld/datadog-operator/pkg/datadog (interfaces: ClientInterface)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	datadog "github.com/zorkian/go-datadog-api"
)

// MockClientInterface is a mock of ClientInterface interface
type MockClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientInterfaceMockRecorder
}

// MockClientInterfaceMockRecorder is the mock recorder for MockClientInterface
type MockClientInterfaceMockRecorder struct {
	mock *MockClientInterface
}

// NewMockClientInterface creates a new mock instance
func NewMockClientInterface(ctrl *gomock.Controller) *MockClientInterface {
	mock := &MockClientInterface{ctrl: ctrl}
	mock.recorder = &MockClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClientInterface) EXPECT() *MockClientInterfaceMockRecorder {
	return m.recorder
}

// CreateMonitor mocks base method
func (m *MockClientInterface) CreateMonitor(arg0 *datadog.Monitor) (*datadog.Monitor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMonitor", arg0)
	ret0, _ := ret[0].(*datadog.Monitor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMonitor indicates an expected call of CreateMonitor
func (mr *MockClientInterfaceMockRecorder) CreateMonitor(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMonitor", reflect.TypeOf((*MockClientInterface)(nil).CreateMonitor), arg0)
}

// DeleteMonitor mocks base method
func (m *MockClientInterface) DeleteMonitor(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMonitor", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMonitor indicates an expected call of DeleteMonitor
func (mr *MockClientInterfaceMockRecorder) DeleteMonitor(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMonitor", reflect.TypeOf((*MockClientInterface)(nil).DeleteMonitor), arg0)
}

// UpdateMonitor mocks base method
func (m *MockClientInterface) UpdateMonitor(arg0 *datadog.Monitor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMonitor", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMonitor indicates an expected call of UpdateMonitor
func (mr *MockClientInterfaceMockRecorder) UpdateMonitor(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMonitor", reflect.TypeOf((*MockClientInterface)(nil).UpdateMonitor), arg0)
}
