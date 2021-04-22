package coupons

import (
	"context"
	"testing"
)

func TestNewCoupon(t *testing.T) {
	t.Parallel()

	_, err := NewCoupon("1", "2", 0)
	wantErr := "coupon maxAllowed is less than 1"

	if err.Error() != wantErr {
		t.Fatalf("%s", "coupon maxAllowed validation failed")
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("create coupon", func(t *testing.T) {
		ctx := context.Background()

		c, _ := NewClient(ctx, &ClientOpts{
			ProjectID: "dummy-project-id",
		})

		defer c.Close()

		coupon, _ := NewCoupon("testTitle", "testDesc", 10)

		id, _ := c.Create(ctx, &coupon)

		if id == "" {
			t.Fatalf("%s", "coupon creation failed")
		}
	})
}

func TestRedeem(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	c, _ := NewClient(ctx, &ClientOpts{
		ProjectID: "dummy-project-id",
	})

	defer c.Close()

	coupon, _ := NewCoupon("testTitle", "testDesc", 1)

	id, _ := c.Create(ctx, &coupon)

	t.Run("redeem coupon once", func(t *testing.T) {
		err := c.Redeem(ctx, id)

		if err != nil {
			t.Fatalf("%s [id:%s]", "coupon redeem failed", id)
		}
	})

	t.Run("should fail redeem with maxAllowed 1", func(t *testing.T) {
		err := c.Redeem(ctx, id)

		if err.Error() != "coupon purchase not allowed, limit reached" {
			t.Fatalf("%s [id:%s]", "redeem success if maxAllowed reached", id)
		}
	})
}
