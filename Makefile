build_cashin:
	@go build ./cmd/cashin-cashout/.
	@echo cashin-cashout built!

build_summary:
	@go build ./cmd/daily-summary/.
	@echo daily-summary built!

build:
	@go build ./cmd/cashin-cashout/.
	@go build ./cmd/daily-summary/.
	@echo all services built!

run_cashin: build_cashin
	@go run ./cmd/cashin-cashout/main.go

run_summary: build_summary
	@go run ./cmd/daily-summary/main.go

test:
	@echo "Testing..."
	go test -v ./...