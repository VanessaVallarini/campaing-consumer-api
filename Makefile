.PHONY: build run compose-infra-down compose-infra-up

compose-infra-up:
	docker-compose -f ./build/package/docker/docker-compose.yml --profile infra up -d

compose-infra-down:
	docker-compose -f ./build/package/docker/docker-compose.yml --profile infra --profile app down

