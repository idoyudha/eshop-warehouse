mock: ### run mockgen
	mockgen -source ./internal/usecase/interfaces.go -package usecase_test > ./internal/usecase/mocks_test.go
.PHONY: mock

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test