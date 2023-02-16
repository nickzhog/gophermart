go-test:
	@go test -count=1 -cover ./... -v
docker-compose-up:
	@docker-compose -f ./docker-compose.yaml up -d