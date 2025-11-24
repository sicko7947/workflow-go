#!/bin/bash
set -e

# Configuration
REGION="${AWS_REGION:-ap-southeast-2}"
TABLE_NAME="${AWS_DYNAMODB_TABLE_NAME:-workflow_executions}"

echo "Creating DynamoDB table: $TABLE_NAME in region: $REGION"

# Create table with generic PK/SK and GSIs for Single Table Design
aws dynamodb create-table \
  --table-name "$TABLE_NAME" \
  --region "$REGION" \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
    AttributeName=GSI1PK,AttributeType=S \
    AttributeName=GSI1SK,AttributeType=S \
    AttributeName=GSI2PK,AttributeType=S \
    AttributeName=GSI2SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  --global-secondary-indexes \
    "[
      {
        \"IndexName\": \"GSI1\",
        \"KeySchema\": [
          {\"AttributeName\": \"GSI1PK\", \"KeyType\": \"HASH\"},
          {\"AttributeName\": \"GSI1SK\", \"KeyType\": \"RANGE\"}
        ],
        \"Projection\": {\"ProjectionType\": \"ALL\"}
      },
      {
        \"IndexName\": \"GSI2\",
        \"KeySchema\": [
          {\"AttributeName\": \"GSI2PK\", \"KeyType\": \"HASH\"},
          {\"AttributeName\": \"GSI2SK\", \"KeyType\": \"RANGE\"}
        ],
        \"Projection\": {\"ProjectionType\": \"ALL\"}
      }
    ]" \
  --tags \
    Key=Project,Value=WorkflowGo \
    Key=Environment,Value=dev \
    Key=Type,Value=WorkflowEngine

echo "Waiting for table to become active..."
aws dynamodb wait table-exists --table-name "$TABLE_NAME" --region "$REGION"

echo "Enabling TTL on 'ttl' attribute..."
aws dynamodb update-time-to-live \
  --table-name "$TABLE_NAME" \
  --region "$REGION" \
  --time-to-live-specification "Enabled=true, AttributeName=ttl"

echo "✅ DynamoDB table '$TABLE_NAME' created successfully!"
echo "Table ARN:"
aws dynamodb describe-table --table-name "$TABLE_NAME" --region "$REGION" --query 'Table.TableArn' --output text
