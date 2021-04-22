export FIRESTORE_EMULATOR_HOST=localhost:8080

test:
	@go test -v ./...

doc:
	@godoc -http=:6060
