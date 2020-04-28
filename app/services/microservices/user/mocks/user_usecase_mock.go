// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUseCase is a mock of UseCase interface
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockUseCase) Create(user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockUseCaseMockRecorder) Create(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUseCase)(nil).Create), user)
}

// GetByID mocks base method
func (m *MockUseCase) GetByID(uid uint) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", uid)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockUseCaseMockRecorder) GetByID(uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUseCase)(nil).GetByID), uid)
}

// GetByNickname mocks base method
func (m *MockUseCase) GetByNickname(nickname string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNickname", nickname)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByNickname indicates an expected call of GetByNickname
func (mr *MockUseCaseMockRecorder) GetByNickname(nickname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNickname", reflect.TypeOf((*MockUseCase)(nil).GetByNickname), nickname)
}

// CheckPassword mocks base method
func (m *MockUseCase) CheckPassword(uid uint, pass []byte) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPassword", uid, pass)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckPassword indicates an expected call of CheckPassword
func (mr *MockUseCaseMockRecorder) CheckPassword(uid, pass interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPassword", reflect.TypeOf((*MockUseCase)(nil).CheckPassword), uid, pass)
}

// Update mocks base method
func (m *MockUseCase) Update(oldPass []byte, newUser models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", oldPass, newUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockUseCaseMockRecorder) Update(oldPass, newUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUseCase)(nil).Update), oldPass, newUser)
}

// Delete mocks base method
func (m *MockUseCase) Delete(uid uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockUseCaseMockRecorder) Delete(uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUseCase)(nil).Delete), uid)
}

// GetUsersByNicknamePart mocks base method
func (m *MockUseCase) GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersByNicknamePart", nicknamePart, limit)
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersByNicknamePart indicates an expected call of GetUsersByNicknamePart
func (mr *MockUseCaseMockRecorder) GetUsersByNicknamePart(nicknamePart, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersByNicknamePart", reflect.TypeOf((*MockUseCase)(nil).GetUsersByNicknamePart), nicknamePart, limit)
}