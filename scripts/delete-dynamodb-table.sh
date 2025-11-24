#!/bin/bash
set -e

# Configuration
REGION="${AWS_REGION:-ap-southeast-2}"
TABLE_NAME="${AWS_DYNAMODB_TABLE_NAME:-workflow_executions}"

echo "⚠️  WARNING: This will permanently delete table: $TABLE_NAME"
read -p "Are you sure? (yes/no): " -r
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "Aborted."
    exit 1
fi

echo "Deleting DynamoDB table: $TABLE_NAME in region: $REGION"
aws dynamodb delete-table --table-name "$TABLE_NAME" --region "$REGION"

echo "Waiting for table to be deleted..."
aws dynamodb wait table-not-exists --table-name "$TABLE_NAME" --region "$REGION"

echo "✅ Table deleted successfully"
