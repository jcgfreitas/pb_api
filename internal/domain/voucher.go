package domain

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Coupon is the base structure representing coupons which are stored in our database
type Coupon struct {
	gorm.Model
	Name   string    `json:"name"`
	Brand  string    `json:"brand"`
	Value  uint      `json:"value"`
	Expiry time.Time `json:"expiry"`
}

type APICoupon struct {
	Name   *string    `json:"name"`
	Brand  *string    `json:"brand"`
	Value  *uint      `json:"value"`
	Expiry *time.Time `json:"expiry"`
}

// NewCoupon instantiates a Coupon from a APICoupon struct
func NewCoupon(APIc APICoupon) Coupon {
	var c Coupon
	if APIc.Name != nil {
		c.Name = *APIc.Name
	}
	if APIc.Brand != nil {
		c.Brand = *APIc.Brand
	}
	if APIc.Value != nil {
		c.Value = *APIc.Value
	}
	if APIc.Expiry != nil {
		c.Expiry = *APIc.Expiry
	}
	return c
}

func UpdateCoupon(c Coupon, APIc APICoupon) Coupon {
	if APIc.Name != nil {
		c.Name = *APIc.Name
	}
	if APIc.Brand != nil {
		c.Brand = *APIc.Brand
	}
	if APIc.Value != nil {
		c.Value = *APIc.Value
	}
	if APIc.Expiry != nil {
		c.Expiry = *APIc.Expiry
	}
	return c
}
