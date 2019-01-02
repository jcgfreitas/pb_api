package repository

import (
	"time"

	"github.com/jcgfreitas/pb_api/internal/domain"
	"github.com/jinzhu/gorm"
)

type GormRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

func (gr *GormRepository) Close() {
	gr.db.Close()
}

func New(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (gr *GormRepository) NewCoupon(APIc domain.APICoupon) error {
	c := domain.NewCoupon(APIc)
	return gr.db.Create(&c).Error
}

func (gr *GormRepository) GetCouponByID(id uint, c *domain.Coupon) error {
	gr.db.First(c, id)
	if c.ID == 0 {
		return domain.NewCouponNotFoundError()
	}
	return gr.db.Error
}

func (gr *GormRepository) DeleteCoupon(id uint) error {
	var c domain.Coupon
	gr.db.First(&c, id)
	if c.ID == 0 {
		return domain.NewCouponNotFoundError()
	}
	return gr.db.Delete(c).Error
}

func (gr *GormRepository) UpdateCoupon(id uint, APIc domain.APICoupon) error {
	var c domain.Coupon
	gr.db.First(&c, id)
	if c.ID == 0 {
		return domain.NewCouponNotFoundError()
	}

	c = domain.UpdateCoupon(c, APIc)

	return gr.db.Save(&c).Error
}

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

func (gr *GormRepository) QueryBatchingFunction(limit, page uint) func() error {
	return func() error {
		gr.tx = gr.tx.Limit(limit).Offset(limit * (page - 1))
		return gr.tx.Error
	}
}

func (gr *GormRepository) QueryLTExpiryFunction(t time.Time) func() error {
	return func() error {
		gr.tx = gr.tx.Where("expiry < ?", t)
		return gr.tx.Error
	}
}

func (gr *GormRepository) QueryGTExpiryFunction(t time.Time) func() error {
	return func() error {
		gr.tx = gr.tx.Where("expiry > ?", t)
		return gr.tx.Error
	}
}

func (gr *GormRepository) QueryLTCreatedFunction(t time.Time) func() error {
	return func() error {
		gr.tx = gr.tx.Where("created_at < ?", t)
		return gr.tx.Error
	}
}

func (gr *GormRepository) QueryGTCreatedFunction(t time.Time) func() error {
	return func() error {
		gr.tx = gr.tx.Where("created_at > ?", t)
		return gr.tx.Error
	}
}

func (gr *GormRepository) QueryLTValueFunction(v uint) func() error {
	return func() error {
		gr.tx = gr.tx.Where("value < ?", v)
		return gr.tx.Error
	}
}

func (gr *GormRepository) QueryGTValueFunction(v uint) func() error {
	return func() error {
		gr.tx = gr.tx.Where("value > ?", v)
		return gr.tx.Error
	}
}
