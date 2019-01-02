postgres:
	docker-compose up

run:
	go run cmd/pb_api/main.go

test:
	go test -cover ./...

integration:
	go test -cover ./... -tags=integration

get:
	curl -X GET http://localhost:8080/coupons/1 -i

create:
	curl -X POST --data '{"name" : "CouponName","brand" : "CouponBrand","value" : 10,"expiry" : "2020-01-01T23:59:59Z"}' http://localhost:8080/coupons -i

gets:
	curl -X GET http://localhost:8080/coupons?value=10 -i

delete:
	curl -X DELETE http://localhost:8080/coupons/1 -i

update:
	curl -X POST --data '{"value" : 20}' http://localhost:8080/coupons/4 -i

