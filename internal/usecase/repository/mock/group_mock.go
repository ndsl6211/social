// Code generated by MockGen. DO NOT EDIT.
// Source: mashu.example/internal/usecase/repository (interfaces: GroupRepo)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	entity "mashu.example/internal/entity"
)

// MockGroupRepo is a mock of GroupRepo interface.
type MockGroupRepo struct {
	ctrl     *gomock.Controller
	recorder *MockGroupRepoMockRecorder
}

// MockGroupRepoMockRecorder is the mock recorder for MockGroupRepo.
type MockGroupRepoMockRecorder struct {
	mock *MockGroupRepo
}

// NewMockGroupRepo creates a new mock instance.
func NewMockGroupRepo(ctrl *gomock.Controller) *MockGroupRepo {
	mock := &MockGroupRepo{ctrl: ctrl}
	mock.recorder = &MockGroupRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupRepo) EXPECT() *MockGroupRepoMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockGroupRepo) Delete(arg0 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockGroupRepoMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockGroupRepo)(nil).Delete), arg0)
}

// GetGroupById mocks base method.
func (m *MockGroupRepo) GetGroupById(arg0 uuid.UUID) (*entity.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupById", arg0)
	ret0, _ := ret[0].(*entity.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupById indicates an expected call of GetGroupById.
func (mr *MockGroupRepoMockRecorder) GetGroupById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupById", reflect.TypeOf((*MockGroupRepo)(nil).GetGroupById), arg0)
}

// GetGroupByName mocks base method.
func (m *MockGroupRepo) GetGroupByName(arg0 string) (*entity.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupByName", arg0)
	ret0, _ := ret[0].(*entity.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupByName indicates an expected call of GetGroupByName.
func (mr *MockGroupRepoMockRecorder) GetGroupByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupByName", reflect.TypeOf((*MockGroupRepo)(nil).GetGroupByName), arg0)
}

// Save mocks base method.
func (m *MockGroupRepo) Save(arg0 *entity.Group) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockGroupRepoMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockGroupRepo)(nil).Save), arg0)
}
