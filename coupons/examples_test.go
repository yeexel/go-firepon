package coupons

import "context"

func ExampleNewClient() {
	ctx := context.Background()
	client, err := NewClient(ctx, &ClientOpts{
		ProjectID: "project-id",
	})
	if err != nil {
		// TODO: Handle error
	}
	defer client.Close() // Close client when done
	_ = client           // TODO: Use client
}

func ExampleClient_Create() {
	ctx := context.Background()
	client, err := NewClient(ctx, &ClientOpts{
		ProjectID: "project-id",
	})
	if err != nil {
		// TODO: Handle error
	}
	defer client.Close() // Close client when done

	// Create coupon struct and specifiy name, description, and maxAllowed values
	coupon, err := NewCoupon("coupon_name", "coupon_desc", 2)

	// Persist coupon in Firestore collection
	couponId, err := client.Create(ctx, &coupon)

	_ = couponId // TODO: Use coupon ID later
}

func ExampleClient_GetOne() {
	ctx := context.Background()
	client, err := NewClient(ctx, &ClientOpts{
		ProjectID: "project-id",
	})
	if err != nil {
		// TODO: Handle error
	}
	defer client.Close() // Close client when done

	// Persist coupon in Firestore collection
	json, err := client.GetOne(ctx, "coupon-id")

	_ = json // {"created_at":"2021-04-21T13:44:04.447064Z","description":"str2","id":"2uIvxCZpEA","max_allowed":200,"purchases":[],"title":"str"}
}

func ExampleClient_GetAll() {
	ctx := context.Background()
	client, err := NewClient(ctx, &ClientOpts{
		ProjectID: "project-id",
	})
	if err != nil {
		// TODO: Handle error
	}
	defer client.Close() // Close client when done

	// Persist coupon in Firestore collection
	json, err := client.GetAll(ctx)

	_ = json // [{"created_at":"2021-04-21T13:44:04.447064Z","description":"str2","id":"2uIvxCZpEA","max_allowed":200,"purchases":[],"title":"str"},{"created_at":"2021-04-21T15:04:03.414614Z","description":"str2qdq","id":"B9qqp2Dqv7","max_allowed":25,"purchases":[],"title":"strw3"},{"created_at":"2021-04-22T09:01:31.670615Z","description":"desc","id":"VSciUPTPGI","max_allowed":2,"purchases":["2021-04-22T09:01:57.613793Z","2021-04-22T09:02:07.323676Z"],"title":"test_coupon"}]
}

func ExampleClient_Redeem() {
	ctx := context.Background()
	client, err := NewClient(ctx, &ClientOpts{
		ProjectID: "project-id",
	})
	if err != nil {
		// TODO: Handle error
	}
	defer client.Close() // Close client when done

	// Redeem coupon and add new entry to "purchases" path.
	// If redeem limit reached - throw error.
	err = client.Redeem(ctx, "VSciUPTPGI")

	if err != nil {
		// TODO: handle error
	}
}
