package coupons

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/yeexel/go-firepon/helpers"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	c *firestore.Client
}

type ClientOpts struct {
	ProjectID string
}

type Coupon interface{}

type coupon struct {
	Title       string      `firestore:"title"`
	Description string      `firestore:"description"`
	MaxAllowed  int         `firestore:"max_allowed"`
	CreatedAt   time.Time   `firestore:"created_at"`
	Purchases   []time.Time `firestore:"purchases"`
}

const (
	defaultIdLen   = 10
	defaultColName = "coupons"
)

func NewCoupon(title string, description string, maxAllowed int) (Coupon, error) {
	if title == "" {
		return nil, errors.New("coupon title is missing")
	}

	if description == "" {
		return nil, errors.New("coupon description is missing")
	}

	if maxAllowed < 1 {
		return nil, errors.New("coupon maxAllowed is less than 1")
	}

	return &coupon{
		Title:       title,
		Description: description,
		MaxAllowed:  maxAllowed,
		CreatedAt:   time.Now(),
		Purchases:   []time.Time{},
	}, nil
}

// NewClient creates a client with Firestore integration to support storage of coupons.
func NewClient(ctx context.Context, opts *ClientOpts) (*Client, error) {
	firestoreClient, err := firestore.NewClient(ctx, opts.ProjectID)

	if err != nil {
		return nil, errors.New("cannot init firestore client")
	}

	return &Client{
		c: firestoreClient,
	}, nil
}

// Create creates a new coupon in Firestore and returns ID.
func (c *Client) Create(ctx context.Context, coupon *Coupon) (string, error) {
	nanoid, _ := gonanoid.New(defaultIdLen)
	couponDocRef := c.c.Doc(helpers.GetDocPath(defaultColName, nanoid))

	_, err := couponDocRef.Create(ctx, &coupon)

	if err != nil {
		return "", errors.New("cannot create coupon")
	}

	return nanoid, nil
}

// GetOne returns coupon from Firestore in JSON format.
func (c *Client) GetOne(ctx context.Context, cID string) (string, error) {
	couponDocRef := c.c.Doc(helpers.GetDocPath(defaultColName, cID))

	couponDocSnap, err := couponDocRef.Get(ctx)

	if status.Code(err) == codes.NotFound {
		return "", errors.New("coupon not found")
	}

	json, err := c.toJSON(couponDocSnap)

	if err != nil {
		return "", errors.New("cannot convert coupon to JSON representation")
	}

	return json, nil
}

// GetAll returns all coupons available in Firestore.
func (c *Client) GetAll(ctx context.Context) (string, error) {
	col := c.c.Collection(defaultColName)
	iter := col.Documents(ctx)
	defer iter.Stop()

	var strSlice []string

	for {
		docSnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", errors.New("cannot iterate over collection")
		}
		json, err := c.toJSON(docSnap)
		if err != nil {
			return "", errors.New("cannot convert coupon to JSON representation")
		}
		strSlice = append(strSlice, json)
	}

	list := fmt.Sprintf("[%s]", strings.Join(strSlice[:], ","))

	return list, nil
}

// Redeem redeems given coupon.
// If max number of redeems reached - throws error.
func (c *Client) Redeem(ctx context.Context, cID string) error {
	var coupon *coupon

	couponDocRef := c.c.Doc(helpers.GetDocPath(defaultColName, cID))

	couponDocSnap, err := couponDocRef.Get(ctx)

	if status.Code(err) == codes.NotFound {
		return errors.New("coupon not found")
	}

	couponDocSnap.DataTo(&coupon)

	if len(coupon.Purchases) == coupon.MaxAllowed {
		return errors.New("coupon purchase not allowed, limit reached")
	}

	_, err = couponDocRef.Update(ctx, []firestore.Update{{Path: "purchases", Value: firestore.ArrayUnion(time.Now())}})

	if err != nil {
		return errors.New("cannot redeem coupon")
	}

	return nil
}

// toJSON helps to marshall Firestore document snaphot into JSON.
// Also appends firestore ID to the response.
func (c *Client) toJSON(docSnap *firestore.DocumentSnapshot) (string, error) {
	var jsonData []byte

	docData := docSnap.Data()
	docData["id"] = docSnap.Ref.ID

	jsonData, err := json.Marshal(docData)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// Closes closes Firestore client.
func (c *Client) Close() error {
	return c.c.Close()
}
