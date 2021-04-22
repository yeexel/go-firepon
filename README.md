# Go Coupons

[![CircleCI](https://circleci.com/gh/yeexel/go-firepon.svg?style=shield)](https://circleci.com/gh/yeexel/go-firepon)

Coupons API client with Firestore integration written in Go.

### Installation

Run command:

`go get github.com/yeexel/go-firepon@v0.1.0`

### Example

```
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
```

### Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/yeexel/go-firepon.svg)](https://pkg.go.dev/github.com/yeexel/go-firepon)

### Development instructions

Run commands:

`docker-compose up -d`

`make test`

### Production setup

Set `GOOGLE_APPLICATION_CREDENTIALS` environment variable which points to your `firebase.json` file, etc.

More info on Firebase Docs page:
[https://firebase.google.com/docs/admin/setup/#initialize-sdk](https://firebase.google.com/docs/admin/setup/#initialize-sdk)
