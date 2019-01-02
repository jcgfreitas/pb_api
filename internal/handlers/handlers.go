//go:generate mockgen -package mocks -destination ../../mocks/service.go github.com/jcgfreitas/pb_api/internal/handlers Service

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jcgfreitas/pb_api/internal/domain"
	"github.com/sirupsen/logrus"
)

const (
	createCouponPath = "/coupons"
	getCouponsPath   = "/coupons"
	getCouponPath    = "/coupons/{id:[0-9]+}"
	deleteCouponPath = "/coupons/{id:[0-9]+}"
	updateCouponPath = "/coupons/{id:[0-9]+}"
)

// Service is the interface used for the API service layer
type Service interface {
	CreateCoupon(APIc domain.APICoupon) error
	GetCoupon(id uint, c *domain.Coupon) error
	DeleteCoupon(id uint) error
	UpdateCoupon(id uint, APIc domain.APICoupon) error
	GetCoupons(coupons *[]domain.Coupon, args map[string][]string) error
}

// Handlers is the structure that holds the API handler functions
type Handlers struct {
	service Service
	logger  *logrus.Logger
}

// NewHandler is the constructor for Handlers
func NewHandlers(service Service, logger *logrus.Logger) *Handlers {
	logger.SetReportCaller(true)
	return &Handlers{service: service, logger: logger}
}

// CreateCouponHandler handles coupon creation requests
func (h *Handlers) CreateCouponHandler(w http.ResponseWriter, r *http.Request) {
	var APIc domain.APICoupon
	if err := json.NewDecoder(r.Body).Decode(&APIc); err != nil {
		h.logger.WithError(err).Debug("failed to decode jason")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCoupon(APIc); err != nil {
		if _, ok := err.(domain.InvalidArgsError); ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.logger.WithError(err).Error("failed to create coupon")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// CreateCouponPath returns the url path associated with the CreateCouponHandler
func (h *Handlers) CreateCouponPath() string {
	return createCouponPath
}

// GetCouponHandler returns the coupon associated with an id
func (h *Handlers) GetCouponHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(w, r)
	if err != nil {
		return
	}

	var c domain.Coupon
	if err = h.service.GetCoupon(id, &c); err != nil {
		if _, ok := err.(domain.CouponNotFoundError); ok {
			h.logger.WithError(err).WithField("id", id).Debug("coupon not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.logger.WithError(err).WithField("id", id).Error("failed to get coupon")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(c)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("failed to Marshal coupon")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetCouponPath returns the url path associated with the GetCouponHandler
func (h *Handlers) GetCouponPath() string {
	return getCouponPath
}

// DeleteCouponHandler deletes a coupon associated with an id
func (h *Handlers) DeleteCouponHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(w, r)
	if err != nil {
		return
	}

	if err = h.service.DeleteCoupon(id); err != nil {
		if _, ok := err.(domain.CouponNotFoundError); ok {
			h.logger.WithError(err).WithField("id", id).Debug("coupon not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.logger.WithError(err).WithField("id", id).Error("failed to delete coupon")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteCouponPath returns the url path associated with the DeleteCouponHandler
func (h *Handlers) DeleteCouponPath() string {
	return deleteCouponPath
}

// UpdateCouponHandler updates an coupon associated with an id
func (h *Handlers) UpdateCouponHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.getID(w, r)
	if err != nil {
		return
	}

	var APIc domain.APICoupon
	if err := json.NewDecoder(r.Body).Decode(&APIc); err != nil {
		h.logger.WithError(err).Debug("failed to decode jason")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.service.UpdateCoupon(id, APIc); err != nil {
		switch err.(type) {
		case domain.CouponNotFoundError:
			h.logger.WithError(err).WithField("id", id).Debug("coupon not found")
			w.WriteHeader(http.StatusNotFound)
			return
		case domain.InvalidArgsError:
			w.WriteHeader(http.StatusBadRequest)
			return
		default:
			h.logger.WithError(err).WithField("id", id).Error("failed to update coupon")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateCouponPath returns the url path associated with the UpdateCouponHandler
func (h *Handlers) UpdateCouponPath() string {
	return updateCouponPath
}

// GetCouponsHandler queries all coupons and filter them accordingly
func (h *Handlers) GetCouponsHandler(w http.ResponseWriter, r *http.Request) {
	var coupons []domain.Coupon
	if err := h.service.GetCoupons(&coupons, r.URL.Query()); err != nil {
		if _, ok := err.(domain.InvalidArgsError); ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h.logger.WithError(err).WithField("query", r.URL.Query()).Error("failed to get coupons")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(coupons)
	if err != nil {
		h.logger.WithError(err).WithField("query", r.URL.Query()).Error("failed to Marshal coupons")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetCouponsPath returns the url path associated with the GetCouponsHandler
func (h *Handlers) GetCouponsPath() string {
	return getCouponsPath
}

func (h *Handlers) getID(w http.ResponseWriter, r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.WithError(err).WithField("id", vars["id"]).Debug("failed to convert id to integer")
		w.WriteHeader(http.StatusBadRequest)
		return 0, err
	}
	return uint(id), nil
}
