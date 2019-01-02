//go:generate mockgen -package mocks -destination ../../mocks/repository.go github.com/jcgfreitas/pb_api/internal/service Repository

package service

import (
	"strconv"
	"time"

	"github.com/jcgfreitas/pb_api/internal/domain"
	"github.com/sirupsen/logrus"
)

const (
	maxLimit            = uint(1000)
	defaultLimit        = uint(200)
	defaultPage         = uint(1)
	queryName           = "name"
	queryBrand          = "brand"
	queryValue          = "value"
	queryLimit          = "limit"
	queryPage           = "page"
	queryLesserExpiry   = "le"
	queryGreaterExpiry  = "ge"
	queryLesserCreated  = "lc"
	queryGreaterCreated = "gc"
	queryLesserValue    = "lv"
	queryGreaterValue   = "gv"
)

type Repository interface {
	NewCoupon(APIc domain.APICoupon) error
	GetCouponByID(id uint, c *domain.Coupon) error
	DeleteCoupon(id uint) error
	UpdateCoupon(id uint, APIc domain.APICoupon) error
	QueryCoupons(coupons *[]domain.Coupon, query map[string]interface{}, functions ...func() error) error
	QueryBatchingFunction(limit, page uint) func() error
	QueryLTExpiryFunction(t time.Time) func() error
	QueryGTExpiryFunction(t time.Time) func() error
	QueryLTCreatedFunction(t time.Time) func() error
	QueryGTCreatedFunction(t time.Time) func() error
	QueryLTValueFunction(v uint) func() error
	QueryGTValueFunction(v uint) func() error
}

type Service struct {
	repo   Repository
	logger *logrus.Logger
}

func NewService(repository Repository, logger *logrus.Logger) *Service {
	logger.SetReportCaller(true)
	return &Service{repo: repository, logger: logger}
}

func (s *Service) CreateCoupon(APIc domain.APICoupon) error {
	if err := createCouponValidation(APIc); err != nil {
		s.logger.WithError(err).Debug("failed to create Coupon")
		return err
	}
	return s.repo.NewCoupon(APIc)
}

func (s *Service) GetCoupon(id uint, c *domain.Coupon) error {
	return s.repo.GetCouponByID(id, c)
}

func (s *Service) DeleteCoupon(id uint) error {
	return s.repo.DeleteCoupon(id)
}

func (s *Service) UpdateCoupon(id uint, APIc domain.APICoupon) error {
	if err := updateCouponValidation(APIc); err != nil {
		s.logger.WithError(err).Debug("failed to update Coupon")
		return err
	}
	return s.repo.UpdateCoupon(id, APIc)
}

func (s *Service) GetCoupons(coupons *[]domain.Coupon, args map[string][]string) error {
	var funcs []func() error
	query := make(map[string]interface{})
	limit := defaultLimit
	page := defaultPage

	for k, v := range args {
		switch k {
		case queryLimit:
			l64, err := strconv.ParseUint(v[0], 10, 32)
			if err != nil {
				s.logger.WithError(err).WithField("value", v[0]).Debug("failed to parse limit")
				return domain.NewInvalidArgsError("failed to parse limit value:" + v[0])
			}
			if l64 == 0 || l64 > uint64(maxLimit) {
				s.logger.WithField("value", v[0]).Debug("invalid limit")
				return domain.NewInvalidArgsError("invalid limit value:" + v[0])
			}
			limit = uint(l64)
		case queryPage:
			p64, err := strconv.ParseUint(v[0], 10, 32)
			if err != nil {
				s.logger.WithError(err).WithField("value", v[0]).Debug("failed to parse page")
				return domain.NewInvalidArgsError("failed to parse page value:" + v[0])
			}
			if p64 == 0 {
				s.logger.WithField("value", v[0]).Debug("invalid page value")
				return domain.NewInvalidArgsError("invalid page value" + v[0])
			}
			page = uint(p64)
		case queryName:
			query[k] = v[0]
		case queryBrand:
			query[k] = v[0]
		case queryValue:
			v64, err := strconv.ParseUint(v[0], 10, 32)
			if err != nil {
				s.logger.WithError(err).WithField("value", v[0]).Debug("failed to parse value")
				return domain.NewInvalidArgsError("failed to parse value:" + v[0])
			}
			query[k] = uint(v64)
		case queryLesserValue:
			lv64, err := strconv.ParseUint(v[0], 10, 32)
			if err != nil {
				s.logger.WithError(err).WithField("value", v[0]).Debug("failed to parse LesserValueLimit")
				return domain.NewInvalidArgsError("failed to parse LesserValueLimit:" + v[0])
			}
			funcs = append(funcs, s.repo.QueryLTValueFunction(uint(lv64)))
		case queryGreaterValue:
			gv64, err := strconv.ParseUint(v[0], 10, 32)
			if err != nil {
				s.logger.WithError(err).WithField("value", v[0]).Debug("failed to parse GreaterValueLimit")
				return domain.NewInvalidArgsError("failed to parse GreaterValueLimit:" + v[0])
			}
			funcs = append(funcs, s.repo.QueryGTValueFunction(uint(gv64)))
		case queryLesserExpiry:
			le, err := time.Parse(time.RFC3339, v[0])
			if err != nil {
				s.logger.WithError(err).WithField("value", v[0]).Debug("failed to parse LesserExpiryLimit")
				return domain.NewInvalidArgsError("failed to parse LesserExpiryLimit:" + v[0])
			}
			funcs = append(funcs, s.repo.QueryLTExpiryFunction(le))
		case queryGreaterExpiry:
			ge, err := time.Parse(time.RFC3339, v[0])
			if err != nil {
				s.logger.WithError(err).WithField("value", v[0]).Debug("failed to parse GreaterExpiryLimit")
				return domain.NewInvalidArgsError("failed to parse GreaterExpiryLimit:" + v[0])
			}
			funcs = append(funcs, s.repo.QueryGTExpiryFunction(ge))
		case queryLesserCreated:
			lc, err := time.Parse(time.RFC3339, v[0])
			if err != nil {
				s.logger.WithError(err).WithField("value", v[0]).Debug("failed to parse LesserCreatedLimit")
				return domain.NewInvalidArgsError("failed to parse LesserCreatedLimit:" + v[0])
			}
			funcs = append(funcs, s.repo.QueryLTCreatedFunction(lc))
		case queryGreaterCreated:
			gc, err := time.Parse(time.RFC3339, v[0])
			if err != nil {
				s.logger.WithError(err).WithField("value", v[0]).Debug("failed to parse GreaterCreatedLimit")
				return domain.NewInvalidArgsError("failed to parse GreaterCreatedLimit:" + v[0])
			}
			funcs = append(funcs, s.repo.QueryGTCreatedFunction(gc))
		}
	}
	funcs = append(funcs, s.repo.QueryBatchingFunction(limit, page))
	return s.repo.QueryCoupons(coupons, query, funcs...)
}

func createCouponValidation(APIc domain.APICoupon) error {
	if APIc.Name == nil || APIc.Brand == nil || APIc.Value == nil || APIc.Expiry == nil {
		return domain.NewInvalidArgsError("coupon fields must not be nil")
	}
	if *APIc.Name == "" {
		return domain.NewInvalidArgsError("coupon name cannot be empty")
	}
	if *APIc.Brand == "" {
		return domain.NewInvalidArgsError("coupon brand cannot be empty")
	}
	if *APIc.Value == 0 {
		return domain.NewInvalidArgsError("coupon value must be bigger than 0")
	}
	if APIc.Expiry.Before(time.Now()) {
		return domain.NewInvalidArgsError("coupon expiry must be after now")
	}
	return nil
}

func updateCouponValidation(APIc domain.APICoupon) error {
	if APIc.Name == nil && APIc.Brand == nil && APIc.Value == nil && APIc.Expiry == nil {
		return domain.NewInvalidArgsError("coupons fields must not be empty")
	}
	if APIc.Name != nil && *APIc.Name == "" {
		return domain.NewInvalidArgsError("coupon name cannot be empty")
	}
	if APIc.Brand != nil && *APIc.Brand == "" {
		return domain.NewInvalidArgsError("coupon brand cannot be empty")
	}
	if APIc.Value != nil && *APIc.Value == 0 {
		return domain.NewInvalidArgsError("coupon value must be bigger than 0")
	}
	if APIc.Expiry != nil && APIc.Expiry.Before(time.Now()) {
		return domain.NewInvalidArgsError("coupon expiry must be after now")
	}
	return nil
}
