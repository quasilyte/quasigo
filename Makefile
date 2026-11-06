GOPATH_DIR=`go env GOPATH`

test:
	go test -count 2 -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...
	go test -run=__ -bench=. ./...
	@echo "everything is OK"

ci-lint:
	# Built from source, so it can handle the Go version that is used to build it.
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0
	$(GOPATH_DIR)/bin/golangci-lint run ./...
	@echo "everything is OK"

lint:
	golangci-lint run ./...
	@echo "everything is OK"

.PHONY: ci-lint lint test
