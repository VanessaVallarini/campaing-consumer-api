.PHONY: build run compose-infra-down compose-infra-up

compose-infra-up:
	docker-compose -f ./build/package/docker/docker-compose.yml --profile infra --profile app up

compose-infra-down:
	docker-compose -f ./build/package/docker/docker-compose.yml --profile infra --profile app down

.PHONY: build run terraform-init terraform-apply

terraform-init:
	cd ./build/package/docker/localstack && terraform init

terraform-apply:
	cd ./build/package/docker/localstack && terraform apply

config-local:
	./config.sh campaing-consumer-api local

config-dev:
	./config.sh campaing-consumer-api dev

config-prod:
	./config.sh campaing-consumer-api prod

run-api:
	go run ./cmd/campaing-consumer-api

run-api-local: config-local run-api

run-api-dev: config-dev run-api

run-api-prod: config-prod run-api