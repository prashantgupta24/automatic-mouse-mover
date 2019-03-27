start:
	go run cmd/main.go

test:
	go test -v -race -failfast ./...