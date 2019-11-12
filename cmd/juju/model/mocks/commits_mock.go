// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/cmd/juju/model (interfaces: CommitsCommandAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	model "github.com/juju/juju/core/model"
	reflect "reflect"
)

// MockCommitsCommandAPI is a mock of CommitsCommandAPI interface
type MockCommitsCommandAPI struct {
	ctrl     *gomock.Controller
	recorder *MockCommitsCommandAPIMockRecorder
}

// MockCommitsCommandAPIMockRecorder is the mock recorder for MockCommitsCommandAPI
type MockCommitsCommandAPIMockRecorder struct {
	mock *MockCommitsCommandAPI
}

// NewMockCommitsCommandAPI creates a new mock instance
func NewMockCommitsCommandAPI(ctrl *gomock.Controller) *MockCommitsCommandAPI {
	mock := &MockCommitsCommandAPI{ctrl: ctrl}
	mock.recorder = &MockCommitsCommandAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommitsCommandAPI) EXPECT() *MockCommitsCommandAPIMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockCommitsCommandAPI) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockCommitsCommandAPIMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCommitsCommandAPI)(nil).Close))
}

// ListCommits mocks base method
func (m *MockCommitsCommandAPI) ListCommits() ([]model.GenerationCommit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCommits")
	ret0, _ := ret[0].([]model.GenerationCommit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListCommits indicates an expected call of ListCommits
func (mr *MockCommitsCommandAPIMockRecorder) ListCommits() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCommits", reflect.TypeOf((*MockCommitsCommandAPI)(nil).ListCommits))
}
