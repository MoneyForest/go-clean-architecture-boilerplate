#!/bin/sh

echo "SQS setup start!"
echo "Creating SQS Queue..."

aws  --endpoint-url=http://aws:4566  sqs create-queue --queue-name sample_queue
echo 'queue created!'

echo "SQS setup Done!"
