GOPATH_DIR=`go env GOPATH`

test:
	go test -count 2 -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./...
	go test -run=__ -bench=. ./...
	@echo "everything is OK"

ci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH_DIR)/bin v1.50.1
	$(GOPATH_DIR)/bin/golangci-lint run ./...
	@echo "everything is OK"

lint:
	golangci-lint run ./...
	@echo "everything is OK"

.PHONY: ci-lint lint test

