// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/happyxhw/iself/pkg/oauth2x (interfaces: Oauth2x)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	http "net/http"
	reflect "reflect"

	model "github.com/happyxhw/iself/model"
	gomock "github.com/golang/mock/gomock"
	oauth2 "golang.org/x/oauth2"
)

// MockOauth2x is a mock of Oauth2x interface.
type MockOauth2x struct {
	ctrl     *gomock.Controller
	recorder *MockOauth2xMockRecorder
}

// MockOauth2xMockRecorder is the mock recorder for MockOauth2x.
type MockOauth2xMockRecorder struct {
	mock *MockOauth2x
}

// NewMockOauth2x creates a new mock instance.
func NewMockOauth2x(ctrl *gomock.Controller) *MockOauth2x {
	mock := &MockOauth2x{ctrl: ctrl}
	mock.recorder = &MockOauth2xMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOauth2x) EXPECT() *MockOauth2xMockRecorder {
	return m.recorder
}

// Client mocks base method.
func (m *MockOauth2x) Client(arg0 context.Context, arg1 *oauth2.Token) *http.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Client", arg0, arg1)
	ret0, _ := ret[0].(*http.Client)
	return ret0
}

// Client indicates an expected call of Client.
func (mr *MockOauth2xMockRecorder) Client(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Client", reflect.TypeOf((*MockOauth2x)(nil).Client), arg0, arg1)
}

// Exchange mocks base method.
func (m *MockOauth2x) Exchange(arg0 context.Context, arg1 string) (*oauth2.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exchange", arg0, arg1)
	ret0, _ := ret[0].(*oauth2.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exchange indicates an expected call of Exchange.
func (mr *MockOauth2xMockRecorder) Exchange(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exchange", reflect.TypeOf((*MockOauth2x)(nil).Exchange), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockOauth2x) GetUser(arg0 context.Context, arg1 *oauth2.Token) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockOauth2xMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockOauth2x)(nil).GetUser), arg0, arg1)
}

// Refresh mocks base method.
func (m *MockOauth2x) Refresh(arg0 context.Context, arg1 *oauth2.Token) (*oauth2.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", arg0, arg1)
	ret0, _ := ret[0].(*oauth2.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh.
func (mr *MockOauth2xMockRecorder) Refresh(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockOauth2x)(nil).Refresh), arg0, arg1)
}
