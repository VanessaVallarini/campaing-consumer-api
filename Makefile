.PHONY: build run compose-infra-down compose-infra-up

compose-infra-up:
	docker-compose -f ./build/package/docker/docker-compose.yml --profile infra --profile app up

compose-infra-down:
	docker-compose -f ./build/package/docker/docker-compose.yml --profile infra --profile app down

terraform-init:
	terraform ./build/package/docker/localstack init

terraform-plan:
	terraform ./build/package/docker/localstack plan

terraform-apply:
	terraform ./build/package/docker/localstack apply