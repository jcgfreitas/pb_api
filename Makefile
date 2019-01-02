postgres:
	docker-compose up

run:
	go run cmd/pb_api/main.go

test:
	go test -cover ./...

integration:
	go test -cover ./... -tags=integration
