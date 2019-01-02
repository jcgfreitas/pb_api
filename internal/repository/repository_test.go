// +build integration

package repository

import (
	"testing"
	"time"

	"github.com/jcgfreitas/pb_api/internal/domain"
	"github.com/jcgfreitas/pb_api/pkg/gormdb/postgres"
	"github.com/stretchr/testify/assert"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	dbName   = "postgres"
	password = "password1"
	name     = "name"
	brand    = "brand"
	value    = uint(10)
	secs     = int64(1000000000)
)

var Time = time.Unix(secs, 0)
var Name = "new" + name
var Brand = "new" + brand
var Value = 2 * value

// TestNewCoupon tests default values insertion and insertion with non existing tables (tests the error)
func TestNewCoupon(t *testing.T) {
	var testCase = []domain.APICoupon{
		{
			Name:   &Name,
			Brand:  &Brand,
			Value:  &Value,
			Expiry: &Time,
		},
		{
			Name:   &Name,
			Value:  &Value,
			Expiry: &Time,
		},
		{
			Name:   &Name,
			Brand:  &Brand,
			Expiry: &Time,
		},
		{
			Name:  &Name,
			Brand: &Brand,
			Value: &Value,
		},
		{
			Name: &Name,
		},
		{
			Brand: &Brand,
		},
		{
			Value: &Value,
		},
		{
			Expiry: &Time,
		},
		{},
	}

	repo := startDB(t)
	defer repo.Close()

	// insert coupons
	for _, c := range testCase {
		if err := repo.NewCoupon(c); err != nil {
			t.Error(err)
		}
	}

	// check row count
	var count int
	repo.db.Model(&domain.Coupon{}).Count(&count)
	if count != len(testCase) {
		t.Error("failed to insert all Coupons")
	}

	// compare expected with result3
	for i, e := range testCase {
		var c domain.Coupon
		repo.db.Find(&c, i+1)
		if e.Name != nil {
			assert.Equal(t, *e.Name, c.Name)
		}
		if e.Brand != nil {
			assert.Equal(t, *e.Brand, c.Brand)
		}
		if e.Value != nil {
			assert.Equal(t, *e.Value, c.Value)
		}
		if e.Expiry != nil {
			assert.Equal(t, e.Expiry.Unix(), c.Expiry.Unix())
		}
	}

	repo.db.DropTableIfExists(&domain.Coupon{})
	if err := repo.NewCoupon(testCase[0]); err == nil {
		t.Error("should error when inserting into an nonexisting table")
	}
}

// TestGetCouponByID test if it can get a record and if the record does not exist it returns a coupon.ID == 0
func TestGetCouponByID(t *testing.T) {
	t.Run("recordExist", testRecordExist)
	t.Run("recordDoesNotExist", testNoRecord)
}

func testRecordExist(t *testing.T) {
	repo := singleRecordDB(t)
	defer repo.Close()
	var rCoupon domain.Coupon

	repo.GetCouponByID(1, &rCoupon)

	assert.Equal(t, rCoupon.Name, name)
	assert.Equal(t, rCoupon.Brand, brand)
	assert.Equal(t, rCoupon.Value, value)
	assert.Equal(t, rCoupon.Expiry.Unix(), secs)
}

func testNoRecord(t *testing.T) {
	repo := singleRecordDB(t)
	defer repo.Close()
	var rCoupon = domain.Coupon{}

	assert.Error(t, repo.GetCouponByID(2, &rCoupon), domain.CouponNotFoundErrorMessage)
	assert.Equal(t, rCoupon.ID, uint(0))
}

func TestDelete(t *testing.T) {
	t.Run("normalDeletion", testDelete)
	t.Run("doesNotExist", testDeleteDoesNotExist)
}

func testDelete(t *testing.T) {
	repo := singleRecordDB(t)
	defer repo.Close()

	assert.Nil(t, repo.DeleteCoupon(1))
}

func testDeleteDoesNotExist(t *testing.T) {
	repo := startDB(t)
	defer repo.Close()

	assert.Equal(t, repo.DeleteCoupon(1), domain.NewCouponNotFoundError())
}

func TestUpdateCoupon(t *testing.T) {
	t.Run("emptyUpdate", testUpdateEmptyCoupon)
	t.Run("updateName", testUpdateName)
	t.Run("updateBrand", testUpdateBrand)
	t.Run("updateValue", testUpdateValue)
	t.Run("updateExpiry", testUpdateExpiry)
	t.Run("updateAllFields", testFullUpdate)
	t.Run("updateCouponNotFound", testUpdateCouponNotFound)
}

func testUpdateEmptyCoupon(t *testing.T) {
	repo := singleRecordDB(t)
	defer repo.Close()
	a := domain.APICoupon{}

	assert.Nil(t, repo.UpdateCoupon(1, a))
}

func testUpdateName(t *testing.T) {
	repo := singleRecordDB(t)
	defer repo.Close()
	a := domain.APICoupon{
		Name: &Name,
	}

	assert.Nil(t, repo.UpdateCoupon(1, a))

	var c domain.Coupon
	repo.db.Find(&c, 1)
	assert.Equal(t, *a.Name, c.Name)
}

func testUpdateBrand(t *testing.T) {
	repo := singleRecordDB(t)
	defer repo.Close()
	a := domain.APICoupon{
		Brand: &Brand,
	}

	assert.Nil(t, repo.UpdateCoupon(1, a))

	var c domain.Coupon
	repo.db.Find(&c, 1)
	assert.Equal(t, *a.Brand, c.Brand)
}

func testUpdateValue(t *testing.T) {
	repo := singleRecordDB(t)
	defer repo.Close()
	a := domain.APICoupon{
		Value: &Value,
	}

	assert.Nil(t, repo.UpdateCoupon(1, a))

	var c domain.Coupon
	repo.db.Find(&c, 1)
	assert.Equal(t, *a.Value, c.Value)
}

func testUpdateExpiry(t *testing.T) {
	repo := singleRecordDB(t)
	defer repo.Close()
	a := domain.APICoupon{
		Expiry: &Time,
	}

	assert.Nil(t, repo.UpdateCoupon(1, a))

	var c domain.Coupon
	repo.db.Find(&c, 1)
	assert.Equal(t, a.Expiry.Unix(), c.Expiry.Unix())
}

func testFullUpdate(t *testing.T) {
	repo := singleRecordDB(t)
	defer repo.Close()
	a := domain.APICoupon{
		Name:   &Name,
		Brand:  &Brand,
		Value:  &Value,
		Expiry: &Time,
	}

	assert.Nil(t, repo.UpdateCoupon(1, a))

	var c domain.Coupon
	repo.db.Find(&c, 1)
	assert.Equal(t, *a.Name, c.Name)
	assert.Equal(t, *a.Brand, c.Brand)
	assert.Equal(t, *a.Value, c.Value)
	assert.Equal(t, a.Expiry.Unix(), c.Expiry.Unix())
}

func testUpdateCouponNotFound(t *testing.T) {
	repo := startDB(t)
	defer repo.Close()
	a := domain.APICoupon{}

	assert.Equal(t, repo.UpdateCoupon(1, a), domain.NewCouponNotFoundError())
}

// TestQueryCoupons tests the function QueryCoupons
func TestQueryCoupons(t *testing.T) {
	t.Run("nilQuery", testNilQuery)
	t.Run("emptyQuery", testEmptyQuery)
	t.Run("name1Query", testName1Query)
	t.Run("name2Query", testName2Query)
	t.Run("brand1Query", testBrand1Query)
	t.Run("brand2Query", testBrand2Query)
	t.Run("value1Query", testValue1Query)
	t.Run("value2Query", testValue2Query)
	t.Run("batching1", testBatching1)
	t.Run("batching2", testBatching2)
	t.Run("batching3", testBatching3)
	t.Run("lesserThanExpiry", testLTExpiry)
	t.Run("greaterThanExpiry", testGTExpiry)
	t.Run("limitedExpiry", testLimitedExpiry)
	t.Run("lesserThanValue", testLTValue)
	t.Run("greaterThanValue", testGTValue)
	t.Run("limitedValue", testLimitedValue)
	t.Run("lesserThanCreated", testLTCreated)
	t.Run("greaterThanCreated", testGTCreated)
	t.Run("limitedCreated", testLimitedCreated)
}

func testNilQuery(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon

	repo.QueryCoupons(&Coupons, nil)

	assert.Equal(t, len(Coupons), 4)
}

func testEmptyQuery(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	query := make(map[string]interface{})

	repo.QueryCoupons(&Coupons, query)

	assert.Equal(t, len(Coupons), 4)
}

func testName1Query(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	query := make(map[string]interface{})
	query["name"] = name + "1"

	repo.QueryCoupons(&Coupons, query)

	assert.Equal(t, len(Coupons), 2)
	for _, c := range Coupons {
		assert.Equal(t, c.Name, query["name"])
	}
}

func testName2Query(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	query := make(map[string]interface{})
	query["name"] = name + "2"

	repo.QueryCoupons(&Coupons, query)

	assert.Equal(t, len(Coupons), 2)
	for _, c := range Coupons {
		assert.Equal(t, c.Name, query["name"])
	}
}

func testBrand1Query(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	query := make(map[string]interface{})
	query["brand"] = brand + "1"

	repo.QueryCoupons(&Coupons, query)

	assert.Equal(t, len(Coupons), 2)
	for _, c := range Coupons {
		assert.Equal(t, c.Brand, query["brand"])
	}
}

func testBrand2Query(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	query := make(map[string]interface{})
	query["brand"] = brand + "2"

	repo.QueryCoupons(&Coupons, query)

	assert.Equal(t, len(Coupons), 2)
	for _, c := range Coupons {
		assert.Equal(t, c.Brand, query["brand"])
	}
}

func testValue1Query(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	query := make(map[string]interface{})
	query["value"] = value

	repo.QueryCoupons(&Coupons, query)

	assert.Equal(t, len(Coupons), 2)
	for _, c := range Coupons {
		assert.Equal(t, c.Value, query["value"])
	}
}

func testValue2Query(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	query := make(map[string]interface{})
	query["value"] = value * 2

	repo.QueryCoupons(&Coupons, query)

	assert.Equal(t, len(Coupons), 2)
	for _, c := range Coupons {
		assert.Equal(t, c.Value, query["value"])
	}
}

func testBatching1(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	batch := repo.QueryBatchingFunction(1, 1)

	repo.QueryCoupons(&Coupons, nil, batch)

	assert.Equal(t, len(Coupons), 1)
	assert.Equal(t, Coupons[0].ID, uint(1))

}

func testBatching2(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	batch := repo.QueryBatchingFunction(2, 1)

	repo.QueryCoupons(&Coupons, nil, batch)

	assert.Equal(t, len(Coupons), 2)
	assert.Equal(t, Coupons[0].ID, uint(1))
	assert.Equal(t, Coupons[1].ID, uint(2))
}

func testBatching3(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	batch := repo.QueryBatchingFunction(1, 4)

	repo.QueryCoupons(&Coupons, nil, batch)

	assert.Equal(t, len(Coupons), 1)
	assert.Equal(t, Coupons[0].ID, uint(4))
}

func testLTExpiry(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	ltExpiry := repo.QueryLTExpiryFunction(time.Unix(secs+1, 0))

	repo.QueryCoupons(&Coupons, nil, ltExpiry)

	assert.Equal(t, len(Coupons), 2)
	assert.Equal(t, Coupons[0].Expiry.Unix(), int64(secs))
	assert.Equal(t, Coupons[1].Expiry.Unix(), int64(secs))
}

func testGTExpiry(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	gtExpiry := repo.QueryGTExpiryFunction(time.Unix(secs+1, 0))

	repo.QueryCoupons(&Coupons, nil, gtExpiry)

	assert.Equal(t, len(Coupons), 2)
	assert.Equal(t, Coupons[0].Expiry.Unix(), int64(secs*2))
	assert.Equal(t, Coupons[1].Expiry.Unix(), int64(secs*2))
}

func testLimitedExpiry(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	ltExpiry := repo.QueryLTExpiryFunction(time.Unix(secs+1, 0))
	gtExpiry := repo.QueryGTExpiryFunction(time.Unix(secs-1, 0))

	repo.QueryCoupons(&Coupons, nil, gtExpiry, ltExpiry)

	assert.Equal(t, len(Coupons), 2)
	assert.Equal(t, Coupons[0].Expiry.Unix(), int64(secs))
	assert.Equal(t, Coupons[1].Expiry.Unix(), int64(secs))
}

func testLTValue(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	ltValue := repo.QueryLTValueFunction(value + 1)

	repo.QueryCoupons(&Coupons, nil, ltValue)

	assert.Equal(t, len(Coupons), 2)
	assert.Equal(t, Coupons[0].Value, value)
	assert.Equal(t, Coupons[1].Value, value)
}

func testGTValue(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	gtValue := repo.QueryGTValueFunction(value + 1)

	repo.QueryCoupons(&Coupons, nil, gtValue)

	assert.Equal(t, len(Coupons), 2)
	assert.Equal(t, Coupons[0].Value, value*2)
	assert.Equal(t, Coupons[1].Value, value*2)
}

func testLimitedValue(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	gtValue := repo.QueryGTValueFunction(value - 1)
	ltValue := repo.QueryLTValueFunction(value + 1)

	repo.QueryCoupons(&Coupons, nil, gtValue, ltValue)

	assert.Equal(t, len(Coupons), 2)
	assert.Equal(t, Coupons[0].Value, value)
	assert.Equal(t, Coupons[1].Value, value)
}

func testLTCreated(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	ltCreated := repo.QueryLTCreatedFunction(time.Now().Add(time.Duration(10 * time.Second)))

	repo.QueryCoupons(&Coupons, nil, ltCreated)

	assert.Equal(t, len(Coupons), 4)
}

func testGTCreated(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	gtCreated := repo.QueryGTCreatedFunction(time.Now().Truncate(time.Duration(10 * time.Second)))

	repo.QueryCoupons(&Coupons, nil, gtCreated)

	assert.Equal(t, len(Coupons), 4)
}

func testLimitedCreated(t *testing.T) {
	repo := multipleRecordDB(t)
	defer repo.Close()
	var Coupons []domain.Coupon
	gtCreated := repo.QueryGTCreatedFunction(time.Now().Truncate(time.Duration(10 * time.Second)))
	ltCreated := repo.QueryLTCreatedFunction(time.Now().Add(time.Duration(10 * time.Second)))

	repo.QueryCoupons(&Coupons, nil, gtCreated, ltCreated)

	assert.Equal(t, len(Coupons), 4)
}

func singleRecordDB(t *testing.T) *GormRepository {
	var coupon = domain.Coupon{
		Name:   name,
		Brand:  brand,
		Value:  value,
		Expiry: time.Unix(secs, 0),
	}

	repo := startDB(t)

	// add coupon
	if err := repo.db.Create(&coupon).Error; err != nil {
		t.Fatal("failed to create coupon:", err)
	}

	return repo
}

func multipleRecordDB(t *testing.T) *GormRepository {
	repo := startDB(t)

	var testBed = []domain.Coupon{
		{
			Name:   name + "1",
			Brand:  brand + "1",
			Value:  value,
			Expiry: time.Unix(secs, 0),
		},
		{
			Name:   name + "2",
			Brand:  brand + "1",
			Value:  value * 2,
			Expiry: time.Unix(2*secs, 0),
		},
		{
			Name:   name + "1",
			Brand:  brand + "2",
			Value:  value,
			Expiry: time.Unix(secs, 0),
		},
		{
			Name:   name + "2",
			Brand:  brand + "2",
			Value:  value * 2,
			Expiry: time.Unix(2*secs, 0),
		},
	}

	// insert testbed into db
	for _, c := range testBed {
		if err := repo.db.Create(&c).Error; err != nil {
			t.Fatal("failed to create coupon:", err)
		}
	}

	return repo
}

func startDB(t *testing.T) *GormRepository {
	// open db with GORM
	db, err := postgres.Open(host, port, user, dbName, password)
	if err != nil {
		t.Fatal(err)
	}

	// drop and create table
	db.DropTableIfExists(&domain.Coupon{})
	db.AutoMigrate(&domain.Coupon{})

	return New(db)
}
