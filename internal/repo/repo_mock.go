// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repo/repo.go

// Package repo is a generated GoMock package.
package repo

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/k0da/tfreg-golang/internal/config"
)

// MockIRepo is a mock of IRepo interface.
type MockIRepo struct {
	ctrl     *gomock.Controller
	recorder *MockIRepoMockRecorder
}

// MockIRepoMockRecorder is the mock recorder for MockIRepo.
type MockIRepoMockRecorder struct {
	mock *MockIRepo
}

// NewMockIRepo creates a new mock instance.
func NewMockIRepo(ctrl *gomock.Controller) *MockIRepo {
	mock := &MockIRepo{ctrl: ctrl}
	mock.recorder = &MockIRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRepo) EXPECT() *MockIRepoMockRecorder {
	return m.recorder
}

// Clone mocks base method.
func (m *MockIRepo) Clone(c config.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clone", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// Clone indicates an expected call of Clone.
func (mr *MockIRepoMockRecorder) Clone(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clone", reflect.TypeOf((*MockIRepo)(nil).Clone), c)
}

// CommitAndPush mocks base method.
func (m *MockIRepo) CommitAndPush(c config.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitAndPush", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitAndPush indicates an expected call of CommitAndPush.
func (mr *MockIRepoMockRecorder) CommitAndPush(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitAndPush", reflect.TypeOf((*MockIRepo)(nil).CommitAndPush), c)
}