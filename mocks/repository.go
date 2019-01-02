// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jcgfreitas/pb_api/internal/service (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	domain "github.com/jcgfreitas/pb_api/internal/domain"
	reflect "reflect"
	time "time"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// DeleteCoupon mocks base method
func (m *MockRepository) DeleteCoupon(arg0 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCoupon", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCoupon indicates an expected call of DeleteCoupon
func (mr *MockRepositoryMockRecorder) DeleteCoupon(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCoupon", reflect.TypeOf((*MockRepository)(nil).DeleteCoupon), arg0)
}

// GetCouponByID mocks base method
func (m *MockRepository) GetCouponByID(arg0 uint, arg1 *domain.Coupon) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCouponByID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetCouponByID indicates an expected call of GetCouponByID
func (mr *MockRepositoryMockRecorder) GetCouponByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCouponByID", reflect.TypeOf((*MockRepository)(nil).GetCouponByID), arg0, arg1)
}

// NewCoupon mocks base method
func (m *MockRepository) NewCoupon(arg0 domain.APICoupon) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewCoupon", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewCoupon indicates an expected call of NewCoupon
func (mr *MockRepositoryMockRecorder) NewCoupon(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewCoupon", reflect.TypeOf((*MockRepository)(nil).NewCoupon), arg0)
}

// QueryBatchingFunction mocks base method
func (m *MockRepository) QueryBatchingFunction(arg0, arg1 uint) func() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryBatchingFunction", arg0, arg1)
	ret0, _ := ret[0].(func() error)
	return ret0
}

// QueryBatchingFunction indicates an expected call of QueryBatchingFunction
func (mr *MockRepositoryMockRecorder) QueryBatchingFunction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryBatchingFunction", reflect.TypeOf((*MockRepository)(nil).QueryBatchingFunction), arg0, arg1)
}

// QueryCoupons mocks base method
func (m *MockRepository) QueryCoupons(arg0 *[]domain.Coupon, arg1 map[string]interface{}, arg2 ...func() error) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryCoupons", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueryCoupons indicates an expected call of QueryCoupons
func (mr *MockRepositoryMockRecorder) QueryCoupons(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryCoupons", reflect.TypeOf((*MockRepository)(nil).QueryCoupons), varargs...)
}

// QueryGTCreatedFunction mocks base method
func (m *MockRepository) QueryGTCreatedFunction(arg0 time.Time) func() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryGTCreatedFunction", arg0)
	ret0, _ := ret[0].(func() error)
	return ret0
}

// QueryGTCreatedFunction indicates an expected call of QueryGTCreatedFunction
func (mr *MockRepositoryMockRecorder) QueryGTCreatedFunction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryGTCreatedFunction", reflect.TypeOf((*MockRepository)(nil).QueryGTCreatedFunction), arg0)
}

// QueryGTExpiryFunction mocks base method
func (m *MockRepository) QueryGTExpiryFunction(arg0 time.Time) func() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryGTExpiryFunction", arg0)
	ret0, _ := ret[0].(func() error)
	return ret0
}

// QueryGTExpiryFunction indicates an expected call of QueryGTExpiryFunction
func (mr *MockRepositoryMockRecorder) QueryGTExpiryFunction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryGTExpiryFunction", reflect.TypeOf((*MockRepository)(nil).QueryGTExpiryFunction), arg0)
}

// QueryGTValueFunction mocks base method
func (m *MockRepository) QueryGTValueFunction(arg0 uint) func() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryGTValueFunction", arg0)
	ret0, _ := ret[0].(func() error)
	return ret0
}

// QueryGTValueFunction indicates an expected call of QueryGTValueFunction
func (mr *MockRepositoryMockRecorder) QueryGTValueFunction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryGTValueFunction", reflect.TypeOf((*MockRepository)(nil).QueryGTValueFunction), arg0)
}

// QueryLTCreatedFunction mocks base method
func (m *MockRepository) QueryLTCreatedFunction(arg0 time.Time) func() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryLTCreatedFunction", arg0)
	ret0, _ := ret[0].(func() error)
	return ret0
}

// QueryLTCreatedFunction indicates an expected call of QueryLTCreatedFunction
func (mr *MockRepositoryMockRecorder) QueryLTCreatedFunction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryLTCreatedFunction", reflect.TypeOf((*MockRepository)(nil).QueryLTCreatedFunction), arg0)
}

// QueryLTExpiryFunction mocks base method
func (m *MockRepository) QueryLTExpiryFunction(arg0 time.Time) func() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryLTExpiryFunction", arg0)
	ret0, _ := ret[0].(func() error)
	return ret0
}

// QueryLTExpiryFunction indicates an expected call of QueryLTExpiryFunction
func (mr *MockRepositoryMockRecorder) QueryLTExpiryFunction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryLTExpiryFunction", reflect.TypeOf((*MockRepository)(nil).QueryLTExpiryFunction), arg0)
}

// QueryLTValueFunction mocks base method
func (m *MockRepository) QueryLTValueFunction(arg0 uint) func() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryLTValueFunction", arg0)
	ret0, _ := ret[0].(func() error)
	return ret0
}

// QueryLTValueFunction indicates an expected call of QueryLTValueFunction
func (mr *MockRepositoryMockRecorder) QueryLTValueFunction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryLTValueFunction", reflect.TypeOf((*MockRepository)(nil).QueryLTValueFunction), arg0)
}

// UpdateCoupon mocks base method
func (m *MockRepository) UpdateCoupon(arg0 uint, arg1 domain.APICoupon) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCoupon", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCoupon indicates an expected call of UpdateCoupon
func (mr *MockRepositoryMockRecorder) UpdateCoupon(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCoupon", reflect.TypeOf((*MockRepository)(nil).UpdateCoupon), arg0, arg1)
}
