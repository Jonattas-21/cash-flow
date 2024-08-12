build one:
	@go build ./cmd/cashin-cashout/.
	@echo cashin-cashout built!

build two:
	@go build ./cmd/daily-summary/.
	@echo daily-summary built!

build all:
	@go build ./cmd/cashin-cashout/.
	@go build ./cmd/daily-summary/.
	@echo all services built!

run one: build one
	@go run ./cmd/cashin-cashout/main.go

run two: build two
	@go run ./cmd/daily-summary/main.go

test:
	@echo "Testing..."
	go test -v ./...