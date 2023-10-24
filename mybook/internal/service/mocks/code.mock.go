// Code generated by MockGen. DO NOT EDIT.
// Source: mybook/internal/service/code.go
//
// Generated by this command:
//
//	mockgen -source=mybook/internal/service/code.go -package=svcmocks -destination=mybook/internal/service/mocks/code.mock.go
//
// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockCodeServicePackage is a mock of CodeServicePackage interface.
type MockCodeServicePackage struct {
	ctrl     *gomock.Controller
	recorder *MockCodeServicePackageMockRecorder
}

// MockCodeServicePackageMockRecorder is the mock recorder for MockCodeServicePackage.
type MockCodeServicePackageMockRecorder struct {
	mock *MockCodeServicePackage
}

// NewMockCodeServicePackage creates a new mock instance.
func NewMockCodeServicePackage(ctrl *gomock.Controller) *MockCodeServicePackage {
	mock := &MockCodeServicePackage{ctrl: ctrl}
	mock.recorder = &MockCodeServicePackageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCodeServicePackage) EXPECT() *MockCodeServicePackageMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockCodeServicePackage) Send(ctx context.Context, biz, phone string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", ctx, biz, phone)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockCodeServicePackageMockRecorder) Send(ctx, biz, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockCodeServicePackage)(nil).Send), ctx, biz, phone)
}

// Verify mocks base method.
func (m *MockCodeServicePackage) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", ctx, biz, phone, inputCode)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Verify indicates an expected call of Verify.
func (mr *MockCodeServicePackageMockRecorder) Verify(ctx, biz, phone, inputCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockCodeServicePackage)(nil).Verify), ctx, biz, phone, inputCode)
}