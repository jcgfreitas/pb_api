package repository

import (
	"time"

	"github.com/jcgfreitas/pb_api/internal/domain"
	"github.com/jinzhu/gorm"
)

// GormRepository handles the flow of control from the service upper layer to the database
type GormRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

// Close closes the underlying db
func (gr *GormRepository) Close() {
	gr.db.Close()
}

// Reset drops coupons table rows
func (gr *GormRepository) Reset() {
	gr.db.AutoMigrate(&domain.Coupon{})
}

// New is the GormRepository constructor
func New(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// NewCoupon creates a new coupon record in the db
func (gr *GormRepository) NewCoupon(APIc domain.APICoupon) error {
	c := domain.NewCoupon(APIc)
	return gr.db.Create(&c).Error
}

// GetCouponByID gets a coupon from the db according to the ID
// If there is no record with the given ID a CouponNotFoundError is returned
func (gr *GormRepository) GetCouponByID(id uint, c *domain.Coupon) error {
	gr.db.First(c, id)
	if c.ID == 0 {
		return domain.NewCouponNotFoundError()
	}
	return gr.db.Error
}

// DeleteCoupon deletes the coupon record with the given ID
// If there is no record with the given ID a CouponNotFoundError is returned
func (gr *GormRepository) DeleteCoupon(id uint) error {
	var c domain.Coupon
	gr.db.First(&c, id)
	if c.ID == 0 {
		return domain.NewCouponNotFoundError()
	}
	return gr.db.Delete(c).Error
}

// UpdateCoupon updates a coupon record with a given ID
// Only the name, brand, value and expiry can be changed
// If there is no record with the given ID a CouponNotFoundError is returned
func (gr *GormRepository) UpdateCoupon(id uint, APIc domain.APICoupon) error {
	var c domain.Coupon
	gr.db.First(&c, id)
	if c.ID == 0 {
		return domain.NewCouponNotFoundError()
	}

	c = domain.UpdateCoupon(c, APIc)

	return gr.db.Save(&c).Error
}

// QueryCoupons queries the db for coupon records according to the query and the variadic functions
//
// Query can be used like a "WHERE {key} = {value}" in sql
// The possible key-values are:
// name: string
// brand: string
// value; uint
//
// the functions that can be used to limit our query are the ones generated through this package with a signature:
// "Query...(...) func error"
func (gr *GormRepository) QueryCoupons(coupons *[]domain.Coupon, query map[string]interface{}, functions ...func() error) error {
	gr.tx = gr.db.Where(query)

	for _, c := range functions {
		if err := c(); err != nil {
			return err
		}
	}

	if err := gr.tx.Find(coupons).Error; err != nil {
		return err
	}
	return nil
}

// QueryBatchingFunction changes the amount of coupons and current page of the query
func (gr *GormRepository) QueryBatchingFunction(limit, page uint) func() error {
	return func() error {
		gr.tx = gr.tx.Limit(limit).Offset(limit * (page - 1))
		return gr.tx.Error
	}
}

// QueryLTExpiryFunction limits the query with a "WHERE expiry < ?"
func (gr *GormRepository) QueryLTExpiryFunction(t time.Time) func() error {
	return func() error {
		gr.tx = gr.tx.Where("expiry < ?", t)
		return gr.tx.Error
	}
}

// QueryGTExpiryFunction limits the query with a "WHERE expiry > ?"
func (gr *GormRepository) QueryGTExpiryFunction(t time.Time) func() error {
	return func() error {
		gr.tx = gr.tx.Where("expiry > ?", t)
		return gr.tx.Error
	}
}

// QueryLTCreatedFunction limits the query with a "WHERE created_at < ?"
func (gr *GormRepository) QueryLTCreatedFunction(t time.Time) func() error {
	return func() error {
		gr.tx = gr.tx.Where("created_at < ?", t)
		return gr.tx.Error
	}
}

// QueryGTCreatedFunction limits the query with "WHERE created_at > ?"
func (gr *GormRepository) QueryGTCreatedFunction(t time.Time) func() error {
	return func() error {
		gr.tx = gr.tx.Where("created_at > ?", t)
		return gr.tx.Error
	}
}

// QueryLTValueFunction limits the query with "WHERE value < ?"
func (gr *GormRepository) QueryLTValueFunction(v uint) func() error {
	return func() error {
		gr.tx = gr.tx.Where("value < ?", v)
		return gr.tx.Error
	}
}

// QueryGTValueFunction limts the query with "WHERE value > ?"
func (gr *GormRepository) QueryGTValueFunction(v uint) func() error {
	return func() error {
		gr.tx = gr.tx.Where("value > ?", v)
		return gr.tx.Error
	}
}
