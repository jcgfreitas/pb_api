package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/jcgfreitas/pb_api/internal/domain"
	"github.com/jcgfreitas/pb_api/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type TestHandlers struct {
	*Handlers
	mock *mocks.MockService
	ctrl *gomock.Controller
	w    *httptest.ResponseRecorder
}

const (
	name  = "name"
	brand = "brand"
	value = uint(10)
	secs  = 100000
)

var Expiry = time.Now().Add(secs * time.Second)
var Name = name
var Brand = brand
var Value = value

func startHandlers(t *testing.T) *TestHandlers {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockService(ctrl)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return &TestHandlers{
		Handlers: NewHandlers(mock, logger),
		mock:     mock,
		ctrl:     ctrl,
		w:        httptest.NewRecorder(),
	}
}

func TestCreateCouponHandler(t *testing.T) {
	t.Run("success", testCreateCouponSuccess)
	t.Run("failedDecoding", testCreateCouponFailedDecoding)
	t.Run("invalidArgs", testCreateCouponInvalidArgs)
	t.Run("error", testCreateCouponError)
}

func testCreateCouponSuccess(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	body := marshalAPICoupon(t)
	r, err := http.NewRequest("POST", "/coupons", body)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	h.mock.EXPECT().CreateCoupon(gomock.Any()).Return(nil)

	h.CreateCouponHandler(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusCreated)
}

func testCreateCouponFailedDecoding(t *testing.T) {
	h := startHandlers(t)

	var b []byte
	body := bytes.NewReader(b)
	r, err := http.NewRequest("POST", "/coupons", body)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	h.CreateCouponHandler(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusBadRequest)

}

func testCreateCouponInvalidArgs(t *testing.T) {
	h := startHandlers(t)

	body := marshalAPICoupon(t)
	r, err := http.NewRequest("POST", "/coupons", body)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	h.mock.EXPECT().CreateCoupon(gomock.Any()).Return(domain.NewInvalidArgsError(""))

	h.CreateCouponHandler(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusBadRequest)
}

func testCreateCouponError(t *testing.T) {
	h := startHandlers(t)

	body := marshalAPICoupon(t)
	r, err := http.NewRequest("POST", "/coupons", body)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	h.mock.EXPECT().CreateCoupon(gomock.Any()).Return(errors.New(""))

	h.CreateCouponHandler(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusInternalServerError)
}

func TestGetCouponHandler(t *testing.T) {
	t.Run("badID", testGetCouponBadID)
	t.Run("success", testGetCouponSuccess)
	t.Run("notFound", testGetCouponNotFound)
	t.Run("serviceError", testGetCouponServiceError)
}

func testGetCouponSuccess(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("POST", "/coupons/4", nil)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(getCouponPath, h.GetCouponHandler).Methods("POST")

	h.mock.EXPECT().GetCoupon(uint(4), gomock.Any()).Return(nil)

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusOK)
}

func testGetCouponBadID(t *testing.T) {
	h := startHandlers(t)

	r, err := http.NewRequest("POST", "/coupons/{id}", nil)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	h.GetCouponHandler(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusBadRequest)
}

func testGetCouponNotFound(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("POST", "/coupons/4", nil)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(getCouponPath, h.GetCouponHandler).Methods("POST")

	h.mock.EXPECT().GetCoupon(uint(4), gomock.Any()).Return(domain.NewCouponNotFoundError())

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusNotFound)
}

func testGetCouponServiceError(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("POST", "/coupons/4", nil)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(getCouponPath, h.GetCouponHandler).Methods("POST")

	h.mock.EXPECT().GetCoupon(uint(4), gomock.Any()).Return(domain.NewInvalidArgsError(""))

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusInternalServerError)
}

func TestDeleteCouponHandler(t *testing.T) {
	t.Run("success", testDeleteCouponSuccess)
	t.Run("badID", testDeleteCouponBadID)
	t.Run("notFound", testDeleteCouponNotFound)
	t.Run("serviceError", testDeleteCouponServiceError)
}

func testDeleteCouponSuccess(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("DELETE", "/coupons/4", nil)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.DeleteCouponPath(), h.DeleteCouponHandler).Methods("DELETE")

	h.mock.EXPECT().DeleteCoupon(uint(4)).Return(nil)

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusOK)
}

func testDeleteCouponBadID(t *testing.T) {
	h := startHandlers(t)

	r, err := http.NewRequest("POST", "/coupons/{id}", nil)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	h.DeleteCouponHandler(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusBadRequest)
}

func testDeleteCouponNotFound(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("DELETE", "/coupons/4", nil)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.DeleteCouponPath(), h.DeleteCouponHandler).Methods("DELETE")

	h.mock.EXPECT().DeleteCoupon(uint(4)).Return(domain.NewCouponNotFoundError())

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusNotFound)
}

func testDeleteCouponServiceError(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("DELETE", "/coupons/4", nil)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.DeleteCouponPath(), h.DeleteCouponHandler).Methods("DELETE")

	h.mock.EXPECT().DeleteCoupon(uint(4)).Return(errors.New(""))

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusInternalServerError)
}

func TestUpdateCouponHandler(t *testing.T) {
	t.Run("decodeFails", testUpdateCouponDecodeFailure)
	t.Run("success", testUpdateCouponSuccess)
	t.Run("notFound", testUpdateCouponNotFound)
	t.Run("serviceError", testUpdateCouponServiceError)
	t.Run("badRequest", testUpdateCouponBadRequest)
}

func testUpdateCouponDecodeFailure(t *testing.T) {
	h := startHandlers(t)

	r, err := http.NewRequest("POST", "/coupons/4", http.NoBody)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.UpdateCouponPath(), h.UpdateCouponHandler).Methods("POST")

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusBadRequest)
}

func testUpdateCouponSuccess(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("POST", "/coupons/4", marshalAPICoupon(t))
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.UpdateCouponPath(), h.UpdateCouponHandler).Methods("POST")

	h.mock.EXPECT().UpdateCoupon(uint(4), gomock.Any()).Return(nil)

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusOK)
}

func marshalAPICoupon(t *testing.T) io.Reader {
	a := domain.APICoupon{
		Name:   &Name,
		Brand:  &Brand,
		Value:  &Value,
		Expiry: &Expiry,
	}

	data, err := json.Marshal(a)
	if err != nil {
		t.Fatal("failed to marshal the APIcoupon")
	}
	return bytes.NewReader(data)
}

func testUpdateCouponNotFound(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("POST", "/coupons/4", marshalAPICoupon(t))
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.UpdateCouponPath(), h.UpdateCouponHandler).Methods("POST")

	h.mock.EXPECT().UpdateCoupon(uint(4), gomock.Any()).Return(domain.NewCouponNotFoundError())

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusNotFound)
}

func testUpdateCouponBadRequest(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("POST", "/coupons/4", marshalAPICoupon(t))
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.UpdateCouponPath(), h.UpdateCouponHandler).Methods("POST")

	h.mock.EXPECT().UpdateCoupon(uint(4), gomock.Any()).Return(domain.NewInvalidArgsError(""))

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusBadRequest)
}

func testUpdateCouponServiceError(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("POST", "/coupons/4", marshalAPICoupon(t))
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.UpdateCouponPath(), h.UpdateCouponHandler).Methods("POST")

	h.mock.EXPECT().UpdateCoupon(uint(4), gomock.Any()).Return(errors.New(""))

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusInternalServerError)
}

func TestGetCouponsHandler(t *testing.T) {
	t.Run("success", testGetCouponsSuccess)
	t.Run("invalidArgs", testGetCouponsInvalidArgs)
	t.Run("serviceError", testGetCouponsServiceError)
}

func testGetCouponsSuccess(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("GET", "/coupons", http.NoBody)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.GetCouponsPath(), h.GetCouponsHandler).Methods("GET")

	h.mock.EXPECT().GetCoupons(gomock.Any(), gomock.Any()).Return(nil)

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusOK)
}

func testGetCouponsInvalidArgs(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("GET", "/coupons", http.NoBody)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.GetCouponsPath(), h.GetCouponsHandler).Methods("GET")

	h.mock.EXPECT().GetCoupons(gomock.Any(), gomock.Any()).Return(domain.NewInvalidArgsError(""))

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusBadRequest)
}

func testGetCouponsServiceError(t *testing.T) {
	h := startHandlers(t)
	defer h.ctrl.Finish()

	r, err := http.NewRequest("GET", "/coupons", http.NoBody)
	if err != nil {
		t.Fatal("failed to create http request")
	}

	router := mux.NewRouter()
	router.HandleFunc(h.GetCouponsPath(), h.GetCouponsHandler).Methods("GET")

	h.mock.EXPECT().GetCoupons(gomock.Any(), gomock.Any()).Return(errors.New(""))

	router.ServeHTTP(h.w, r)
	assert.Equal(t, h.w.Code, http.StatusInternalServerError)
}
