provider "aws" {
    version = "~> 2.0"
    region = "us-east-1"
    access_key = "teste"
    secret_key = "teste"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    skip_requesting_account_id  = true
    endpoints {
        sns = "http://localhost:4566"
        sqs = "http://localhost:4566"
  }
}

resource "aws_sns_topic" "topic_campaing" {
  name = "topic_campaing"
}

resource "aws_sqs_queue" "queue_campaing" {
  name = "queue_campaing"
}

resource "aws_sns_topic_subscription" "client_campaing_providers_sqs_target" {
  topic_arn = aws_sns_topic.topic_campaing.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.queue_campaing.arn
  filter_policy = "${jsonencode({
    action = ["campaing"]
  })}"
}