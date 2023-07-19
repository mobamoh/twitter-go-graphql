run:
	go run cmd/graphqlserver/*.go
test:
	go test ./... -count=1
integration:
	go test ./... --tags="integration" -count=1
mock:
	mockery --all --keeptree
migrate:
	migrate -source file://postgres/migrations \
			-database postgres://postgres:nimda@127.0.0.1:5432/twitter_go_dev?sslmode=disable up
rollback:
	migrate -source file://postgres/migrations \
			-database postgres://postgres:nimda@127.0.0.1:5432/twitter_go_dev?sslmode=disable down 1
drop:
	migrate -source file://postgres/migrations \
			-database postgres://postgres:nimda@127.0.0.1:5432/twitter_go_dev?sslmode=disable drop
migration:
	@read -p "Enter migration name: " name; \
		migrate create -ext sql -dir postgres/migrations $$name
generate:
	go generate ./...