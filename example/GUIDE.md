# Simple Math Workflow Guide

This guide provides instructions on how to use the `simple_math` workflow, a demonstration implementation using the Workflow Engine.

## Overview

The **Simple Math Workflow** demonstrates a sequential workflow that performs basic arithmetic operations and formatting. It showcases:
- Data passing between steps.
- Step configuration (retries, timeouts).
- Workflow orchestration.
- HTTP API integration.

## Workflow Logic

The workflow consists of three sequential steps:

1.  **Add Step** (`add`)
    - **Input**: `Val1` (int), `Val2` (int), `Mult` (int)
    - **Operation**: Calculates `Sum = Val1 + Val2`.
    - **Output**: `Value` (Sum), `Mult` (passed through).

2.  **Multiply Step** (`multiply`)
    - **Input**: `Value` (from Add step), `Factor` (mapped from `Mult`).
    - **Operation**: Calculates `Product = Value * Factor`.
    - **Configuration**: Retries = 5, Backoff = Exponential, Timeout = 5s.
    - **Output**: `Value` (Product).

3.  **Format Step** (`format`)
    - **Input**: `Number` (from Multiply step).
    - **Operation**: Formats string "The final result is {Number}".
    - **Output**: `Message` (string).

## Prerequisites

- **Go 1.21+**
- **AWS DynamoDB**: Required for workflow state persistence. You can use a local instance (e.g., LocalStack or DynamoDB Local).

### Environment Variables

The application expects the following environment variables (based on `config` usage):

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_ENVIRONMENT` | Environment name | `local` |
| `AWS_REGION` | AWS Region | `us-east-1` |
| `AWS_DYNAMODB_TABLE_NAME` | DynamoDB Table Name | `workflows` |
| `SERVER_PORT` | HTTP Server Port | `8080` |

*Note: Actual variable names depend on the `internal/config` implementation. Ensure your `.env` matches the expected config keys.*

## Running the Server

To run the workflow server locally:

```bash
# Navigate to the simple_math main directory
cd simple_math/main

# Run the server
go run main.go
```

The server will start on the configured port (default usually `8080`).

## API Usage

### 1. Start a Workflow

Trigger a new execution of the Simple Math workflow.

**Endpoint**: `POST /api/v1/workflows/simple-math`

**Request Body**:
```json
{
  "val1": 10,
  "val2": 5,
  "mult": 2
}
```

**Example Request**:
```bash
curl -X POST http://localhost:8080/api/v1/workflows/simple-math \
  -H "Content-Type: application/json" \
  -d '{"val1": 10, "val2": 5, "mult": 2}'
```

**Response**:
```json
{
  "runId": "550e8400-e29b-41d4-a716-446655440000",
  "status": "PENDING",
  "message": "Workflow started successfully"
}
```

### 2. Check Workflow Status

Retrieve the current status, step execution details, and final output.

**Endpoint**: `GET /api/v1/workflows/:runId`

**Example Request**:
```bash
curl http://localhost:8080/api/v1/workflows/550e8400-e29b-41d4-a716-446655440000
```

**Response (Completed)**:
```json
{
  "WorkflowRun": {
    "ID": "550e8400-e29b-41d4-a716-446655440000",
    "Status": "COMPLETED",
    ...
  },
  "stepExecutions": [
    { "StepID": "add", "Status": "COMPLETED", ... },
    { "StepID": "multiply", "Status": "COMPLETED", ... },
    { "StepID": "format", "Status": "COMPLETED", ... }
  ],
  "output": {
    "message": "The final result is 30"
  }
}
```

### 3. Cancel Workflow

Cancel a running workflow.

**Endpoint**: `POST /api/v1/workflows/:runId/cancel`

**Example Request**:
```bash
curl -X POST http://localhost:8080/api/v1/workflows/550e8400-e29b-41d4-a716-446655440000/cancel
```

## Troubleshooting

- **Failed to start workflow**: Check DynamoDB connection and table existence.
- **Import errors**: Ensure the `go.mod` file correctly points to the workflow engine module. If the module paths do not match, you may need to use `go work` or update `go.mod`.
