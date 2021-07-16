# Copyright (c) 2021 Nikos Leivadaris
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

build:
	go build ./...

clean:
	go clean

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

vet:
	go vet ./...

lint:
	golangci-lint run --enable-all