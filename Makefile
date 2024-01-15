.PHONY: build run compose-infra-down compose-infra-up

compose-infra-up:
	docker-compose -f ./build/package/docker/docker-compose.yml --profile infra --profile app up

compose-infra-down:
	docker-compose -f ./build/package/docker/docker-compose.yml --profile infra --profile app down

#terraform init

#terraform plan

#terraform apply

#export AWS_SECRET_ACCESS_KEY=teste
#export AWS_ACCESS_KEY_ID=teste