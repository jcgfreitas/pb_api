package service

import (
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jcgfreitas/pb_api/internal/domain"

	"github.com/golang/mock/gomock"
	"github.com/jcgfreitas/pb_api/mocks"
	"github.com/sirupsen/logrus"

	"testing"
)

const (
	name    = "name"
	brand   = "brand"
	value   = uint(10)
	sValue  = "10"
	secs    = 100000
	sExpiry = "2025-03-01T23:59:59Z"
)

var Expiry = time.Now().Add(secs * time.Second)
var Name = name
var Brand = brand
var Value = value

type TestService struct {
	*Service
	mock *mocks.MockRepository
	ctrl *gomock.Controller
}

func startService(t *testing.T) *TestService {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockRepository(ctrl)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return &TestService{
		NewService(mock, logger),
		mock,
		ctrl,
	}
}

func TestCreateCoupon(t *testing.T) {
	t.Run("success", testCreateCouponSuccess)
	t.Run("invalidName", testCreateCouponInvalidName)
	t.Run("nilName", testCreateCouponNilName)
	t.Run("invalidBrand", testCreateCouponInvalidBrand)
	t.Run("invalidValue", testCreateCouponInvalidValue)
	t.Run("invalidExpiry", testCreateCouponInvalidExpiry)
}

func testCreateCouponSuccess(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	a := domain.APICoupon{
		Name:   &Name,
		Brand:  &Brand,
		Value:  &Value,
		Expiry: &Expiry,
	}

	s.mock.EXPECT().NewCoupon(a).Return(nil)
	assert.Nil(t, s.CreateCoupon(a))
}

func testCreateCouponNilName(t *testing.T) {
	s := startService(t)

	a := domain.APICoupon{
		Brand:  &Brand,
		Value:  &Value,
		Expiry: &Expiry,
	}

	assert.Error(t, s.CreateCoupon(a))
}

func testCreateCouponInvalidName(t *testing.T) {
	s := startService(t)

	n := ""
	a := domain.APICoupon{
		Name:   &n,
		Brand:  &Brand,
		Value:  &Value,
		Expiry: &Expiry,
	}

	assert.Error(t, s.CreateCoupon(a))
}

func testCreateCouponInvalidBrand(t *testing.T) {
	s := startService(t)

	b := ""
	a := domain.APICoupon{
		Name:   &Name,
		Brand:  &b,
		Value:  &Value,
		Expiry: &Expiry,
	}

	assert.Error(t, s.CreateCoupon(a))
}

func testCreateCouponInvalidValue(t *testing.T) {
	s := startService(t)

	v := uint(0)
	a := domain.APICoupon{
		Name:   &Name,
		Brand:  &Brand,
		Value:  &v,
		Expiry: &Expiry,
	}

	assert.Error(t, s.CreateCoupon(a))
}

func testCreateCouponInvalidExpiry(t *testing.T) {
	s := startService(t)

	e := time.Now().Truncate(secs * time.Second)
	a := domain.APICoupon{
		Name:   &Name,
		Brand:  &Brand,
		Value:  &Value,
		Expiry: &e,
	}

	assert.Error(t, s.CreateCoupon(a))
}

func TestGetCoupon(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	var c domain.Coupon
	s.mock.EXPECT().GetCouponByID(uint(1), &c).Return(nil)
	assert.Nil(t, s.GetCoupon(1, &c))
}

func TestDeleteCoupon(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	s.mock.EXPECT().DeleteCoupon(uint(1)).Return(nil)
	assert.Nil(t, s.DeleteCoupon(1))
}

func TestUpdateCoupon(t *testing.T) {
	t.Run("success", testUpdateCouponSuccess)
	t.Run("emptyUpdate", testUpdateCouponEmpty)
	t.Run("invalidName", testUpdateCouponInvalidName)
	t.Run("invalidBrand", testUpdateCouponInvalidBrand)
	t.Run("invalidValue", testUpdateCouponInvalidValue)
	t.Run("invalidExpiry", testUpdateCouponInvalidExpiry)
}

func testUpdateCouponSuccess(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	a := domain.APICoupon{
		Name:   &Name,
		Brand:  &Brand,
		Value:  &Value,
		Expiry: &Expiry,
	}

	s.mock.EXPECT().UpdateCoupon(uint(1), a).Return(nil)
	assert.Nil(t, s.UpdateCoupon(uint(1), a))
}

func testUpdateCouponEmpty(t *testing.T) {
	s := startService(t)

	a := domain.APICoupon{}

	assert.Error(t, s.UpdateCoupon(uint(1), a))
}

func testUpdateCouponInvalidName(t *testing.T) {
	s := startService(t)

	n := ""
	a := domain.APICoupon{
		Name: &n,
	}

	assert.Error(t, s.UpdateCoupon(1, a))
}

func testUpdateCouponInvalidBrand(t *testing.T) {
	s := startService(t)

	b := ""
	a := domain.APICoupon{
		Brand: &b,
	}

	assert.Error(t, s.UpdateCoupon(1, a))
}

func testUpdateCouponInvalidValue(t *testing.T) {
	s := startService(t)

	v := uint(0)
	a := domain.APICoupon{
		Value: &v,
	}

	assert.Error(t, s.UpdateCoupon(1, a))
}

func testUpdateCouponInvalidExpiry(t *testing.T) {
	s := startService(t)

	e := time.Now().Truncate(secs * time.Second)
	a := domain.APICoupon{
		Expiry: &e,
	}

	assert.Error(t, s.UpdateCoupon(1, a))
}

func TestGetCoupons(t *testing.T) {
	t.Run("successLimit", testGetCouponsSuccessLimit)
	t.Run("successPage", testGetCouponsSuccessPage)
	t.Run("successQuery", testGetCouponsSuccessQuery)
	t.Run("successLesserValue", testGetCouponsSuccessLesserValue)
	t.Run("successGreaterValue", testGetCouponsSuccessGreaterValue)
	t.Run("successLesserExpiry", testGetCouponsSuccessLesserExpiry)
	t.Run("successGreaterExpiry", testGetCouponsSuccessGreaterExpiry)
	t.Run("successLesserCreated", testGetCouponsSuccessLesserCreated)
	t.Run("successGreaterCreated", testGetCouponsSuccessGreaterCreated)
	t.Run("invalidLimit", testGetCouponsInvalidLimit)
	t.Run("invalidPage", testGetCouponsInvalidPage)
	t.Run("invalidValue", testGetCouponsInvalidValue)
	t.Run("invalidLesserValue", testGetCouponsInvalidLesserValue)
	t.Run("invalidGreaterValue", testGetCouponsInvalidGreaterValue)
	t.Run("invalidLesserExpiry", testGetCouponsInvalidLesserExpiry)
	t.Run("invalidGreaterExpiry", testGetCouponsInvalidGreaterExpiry)
	t.Run("invalidLesserCreated", testGetCouponsInvalidLesserCreated)
	t.Run("invalidGreaterCreated", testGetCouponsInvalidGreaterCreated)
}

func testGetCouponsSuccessLimit(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	limit := uint(10)
	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryLimit] = []string{"10"}
	query := make(map[string]interface{})

	s.mock.EXPECT().QueryBatchingFunction(limit, defaultPage)
	s.mock.EXPECT().QueryCoupons(&coupons, query, gomock.Any()).Return(nil)

	assert.Nil(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsSuccessPage(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	page := uint(10)
	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryPage] = []string{"10"}
	query := make(map[string]interface{})

	s.mock.EXPECT().QueryBatchingFunction(defaultLimit, page)
	s.mock.EXPECT().QueryCoupons(&coupons, query, gomock.Any()).Return(nil)

	assert.Nil(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsSuccessQuery(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	//input
	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryName] = []string{Name}
	args[queryBrand] = []string{Brand}
	args[queryValue] = []string{sValue}

	//expected in mock
	query := make(map[string]interface{})
	query[queryName] = Name
	query[queryBrand] = Brand
	query[queryValue] = Value

	s.mock.EXPECT().QueryBatchingFunction(defaultLimit, defaultPage)
	s.mock.EXPECT().QueryCoupons(&coupons, query, gomock.Any()).Return(nil)

	assert.Nil(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsSuccessLesserValue(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryLesserValue] = []string{sValue}
	query := make(map[string]interface{})

	s.mock.EXPECT().QueryLTValueFunction(value)
	s.mock.EXPECT().QueryBatchingFunction(defaultLimit, defaultPage)
	s.mock.EXPECT().QueryCoupons(&coupons, query, gomock.Any()).Return(nil)

	assert.Nil(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsSuccessGreaterValue(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryGreaterValue] = []string{sValue}
	query := make(map[string]interface{})

	s.mock.EXPECT().QueryGTValueFunction(value)
	s.mock.EXPECT().QueryBatchingFunction(defaultLimit, defaultPage)
	s.mock.EXPECT().QueryCoupons(&coupons, query, gomock.Any()).Return(nil)

	assert.Nil(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsSuccessLesserExpiry(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryLesserExpiry] = []string{sExpiry}
	query := make(map[string]interface{})

	s.mock.EXPECT().QueryLTExpiryFunction(gomock.Any())
	s.mock.EXPECT().QueryBatchingFunction(defaultLimit, defaultPage)
	s.mock.EXPECT().QueryCoupons(&coupons, query, gomock.Any()).Return(nil)

	assert.Nil(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsSuccessGreaterExpiry(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryGreaterExpiry] = []string{sExpiry}
	query := make(map[string]interface{})

	s.mock.EXPECT().QueryGTExpiryFunction(gomock.Any())
	s.mock.EXPECT().QueryBatchingFunction(defaultLimit, defaultPage)
	s.mock.EXPECT().QueryCoupons(&coupons, query, gomock.Any()).Return(nil)

	assert.Nil(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsSuccessLesserCreated(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryLesserCreated] = []string{sExpiry}
	query := make(map[string]interface{})

	s.mock.EXPECT().QueryLTCreatedFunction(gomock.Any())
	s.mock.EXPECT().QueryBatchingFunction(defaultLimit, defaultPage)
	s.mock.EXPECT().QueryCoupons(&coupons, query, gomock.Any()).Return(nil)

	assert.Nil(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsSuccessGreaterCreated(t *testing.T) {
	s := startService(t)
	defer s.ctrl.Finish()

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryGreaterCreated] = []string{sExpiry}
	query := make(map[string]interface{})

	s.mock.EXPECT().QueryGTCreatedFunction(gomock.Any())
	s.mock.EXPECT().QueryBatchingFunction(defaultLimit, defaultPage)
	s.mock.EXPECT().QueryCoupons(&coupons, query, gomock.Any()).Return(nil)

	assert.Nil(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsInvalidLimit(t *testing.T) {
	s := startService(t)

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryLimit] = []string{name}

	assert.Error(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsInvalidPage(t *testing.T) {
	s := startService(t)

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryPage] = []string{name}

	assert.Error(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsInvalidValue(t *testing.T) {
	s := startService(t)

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryValue] = []string{name}

	assert.Error(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsInvalidLesserValue(t *testing.T) {
	s := startService(t)

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryLesserValue] = []string{name}

	assert.Error(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsInvalidGreaterValue(t *testing.T) {
	s := startService(t)

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryGreaterValue] = []string{name}

	assert.Error(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsInvalidLesserExpiry(t *testing.T) {
	s := startService(t)

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryLesserExpiry] = []string{name}

	assert.Error(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsInvalidGreaterExpiry(t *testing.T) {
	s := startService(t)

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryGreaterExpiry] = []string{name}

	assert.Error(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsInvalidLesserCreated(t *testing.T) {
	s := startService(t)

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryLesserCreated] = []string{name}

	assert.Error(t, s.GetCoupons(&coupons, args))
}

func testGetCouponsInvalidGreaterCreated(t *testing.T) {
	s := startService(t)

	var coupons []domain.Coupon
	args := make(map[string][]string)
	args[queryGreaterCreated] = []string{name}

	assert.Error(t, s.GetCoupons(&coupons, args))
}
