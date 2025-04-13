#!/usr/bin/bash

go vet -v ./...
go test $(go list ./... | grep -v /tests/integration)
go build -v
