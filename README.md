# campaing-consumer-api

#terraform plan

#export AWS_SECRET_ACCESS_KEY=teste
#export AWS_ACCESS_KEY_ID=teste
export app=campaing-consumer-api
export environment=local


aws --endpoint-url=http://localhost:4566 sqs purge-queue --queue-url http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/queue_campaing