# campaing-consumer-api
Application responsible for listening to campaign creation, update and deletion queues

## How to run locally
- Install [Go](https://go.dev/)
- Install docker [Rancher Desktop](https://ifood.atlassian.net/wiki/spaces/IL/pages/3049586786/Migrando+do+Docker+Desktop+para+o+Rancher+Desktop+no+Mac) / [Colima](https://ifood.atlassian.net/wiki/spaces/EN/pages/2971992107/Instala+o+do+Docker+no+MacBook+M1+e+Intel)

### Assemble Configuration File
This section is just to explain how to run configuration if necessary.

Run `./config.sh {deploymentName} {environment}`, this command will merge, chart values of deployment with override environment file and the secret file
The `deploymentName` is the name of the folder in k8s folder (./k8s/{deploymentName}), are available:
- campaing-consumer-api
To simplify in makefile we have the following commands:

| Command             | Environment                 |
|---------------------|-----------------------------|
| make config-local   | campaing-consumer-api local |
| make config-dev     | campaing-consumer-api dev   |
| make config-prod    | campaing-consumer-api prod  |

The `environment` are available:
- local
- dev
- prod

### Prepare AWS and Start Docker
```
make compose-infra-up
```
```
make terraform-init
```
```
make terraform-apply
```

### Launch
```shell
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "cwd": "${workspaceFolder}",
            "program": "${cwd}/cmd/campaing-consumer-api/main.go"
        }
    ]
}
```

### Starting App No Debug
Once this is done, open another terminal and start the docker compose with dependencies using `docker compose up`, wait for the containers to start and run `make docker-create-topics` to create all the necessary topics
and use to start the app we can use:
| Command               | Environment           |
|-----------------------|-----------------------|
| make run-api-local    | campaing-consumer-api local |
| make run-api-dev      | campaing-consumer-api dev   |
| make run-api-prod     | campaing-consumer-api prod  |

### SQS Comands
#### Clear SQS
```
aws --endpoint-url=http://localhost:4566 sqs purge-queue --queue-url http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/queue_campaing
```
#### Send message create campaing
```
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/queue_campaing --message-body '{"user_id": "c3eeb9b0-051c-4803-b0a4-f6060bcb40d9", "slug_id": "f43e580b-ffb2-490d-aea1-b2f0435d624b", "merchant_id": "2ed8b772-1714-46de-98ab-c2653bb03d78", "lat": 45.6085, "long": -73.5493, "action": "C"}'
```
#### Send message updated campaing
- Use the campaign id created in the previous step
```
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/queue_campaing --message-body '{"id":"5b228816-54a0-4e2b-8e62-a00ba148c973","user_id":"c3eeb9b0-051c-4803-b0a4-f6060bcb40d9","slug_id":"f43e580b-ffb2-490d-aea1-b2f0435d624b","merchant_id":"2ed8b772-1714-46de-98ab-c2653bb03d78","active":true,"lat":45.6085,"long":-73.5493,"clicks":15,"impressions":50,"action":"U"}'
```
#### Send message delete campaing
- Use the campaign id created in the previous step
```
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/queue_campaing --message-body '{"id":"5b228816-54a0-4e2b-8e62-a00ba148c973","action":"D"}'
```