package domain

const (
	CouponNotFoundErrorMessage = "coupon not found"
)

// CouponNotFoundError is the error passed when the coupon does not exist in the DB
type CouponNotFoundError struct{}

// Error implements the error interface
func (err CouponNotFoundError) Error() string {
	return CouponNotFoundErrorMessage
}

// NewCouponNotFoundError is the constructor for CouponNotFoundError
func NewCouponNotFoundError() error {
	return CouponNotFoundError{}
}

// InvalidCouponError is the error passed when a user sends the wrong coupon arguments
type InvalidArgsError struct {
	msg string
}

// Error implements the error interface
func (err InvalidArgsError) Error() string {
	return err.msg
}

// NewInvalidCouponError is the constructor for InvalidCouponError
func NewInvalidArgsError(msg string) error {
	return InvalidArgsError{msg: msg}
}
