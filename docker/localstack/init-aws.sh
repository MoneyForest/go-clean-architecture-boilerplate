#!/bin/sh

echo "SQS setup start!"
echo "Creating SQS Queue..."

aws  --endpoint-url=http://localstack:4566  sqs create-queue --queue-name sample_queue

echo 'queue created!'

aws --endpoint-url=http://localstack:4566 sqs list-queues

echo "SQS setup Done!"
