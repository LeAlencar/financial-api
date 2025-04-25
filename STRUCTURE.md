# Project Structure

This document describes the organization of our microservices architecture within a monorepo structure.

## Directory Structure

```
.
├── services/                  # All microservices
│   ├── s1-generator/         # Service 1: Message Generator
│   │   ├── src/             # Source code
│   │   └── tests/           # Unit and integration tests
│   │
│   ├── s2-processor/        # Service 2: Processing Services
│   │   ├── user-service/    # User data processor
│   │   ├── quotation-service/ # Currency quotation processor
│   │   └── transaction-service/ # Transaction processor
│   │
│   └── s3-validator/        # Service 3: Log and Validation
│       ├── src/
│       └── tests/
│
├── shared/                   # Shared components
│   ├── utils/               # Common utilities
│   ├── models/              # Shared data models
│   └── config/              # Shared configurations
│
├── deployments/             # Deployment configurations
├── docker-compose.yml       # Service orchestration
└── README.md               # Project documentation
```

## Services Description

### S1 - Message Generator

- Generates mock data for users, quotations, and transactions
- Publishes messages to RabbitMQ

### S2 - Processing Services

- **User Service**: Handles user data with PostgreSQL
- **Quotation Service**: Manages currency quotations with MongoDB
- **Transaction Service**: Processes transactions with Cassandra

### S3 - Validator

- Logs all messages
- Validates message integrity
- Provides audit trail

## Development Workflow

1. Each service is independently deployable
2. Shared code is in the `shared/` directory
3. Use `docker-compose up` to run the entire system
4. For local development:
   - Run dependencies via docker-compose
   - Run specific service locally

## Building and Running

### Run Everything

```bash
docker-compose up
```

### Run Specific Service

```bash
docker-compose up s1-generator
```

### Local Development

```bash
# Run dependencies
docker-compose up rabbitmq postgres mongodb cassandra

# Run service locally
cd services/s1-generator
go run cmd/main.go
```

## Communication

- All inter-service communication happens through RabbitMQ
- Each service has its own message queues
- Use shared models for message formats

## Database Access

- PostgreSQL: User data (via Supabase)
- MongoDB: Currency quotations
- Cassandra: Transaction history

## Testing

Each service has its own tests directory:

- Unit tests
- Integration tests
- E2E tests (in the future)

## Deployment

The `deployments/` directory contains:

- Kubernetes manifests
- Environment configurations
- CI/CD pipelines
