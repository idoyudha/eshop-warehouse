mock: ### run mockgen
	mockgen -source ./internal/usecase/interfaces.go -package usecase_test > ./internal/usecase/mocks_test.go
.PHONY: mock

test: ### run test
	@go test -v -race -coverprofile=coverage.out -coverpkg=./... ./internal/...
	@go tool cover -func=coverage.out
.PHONY: test