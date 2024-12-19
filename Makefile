build:
	@go build ./cmd/cashin-cashout/.
	@echo all services built!

run: build
	@go run ./cmd/cashin-cashout/main.go
	@echo services is running!
test:
	@echo "Testing..."
	go test -v ./...