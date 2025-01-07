# Warehouse Service
Part of [eshop](https://github.com/idoyudha/eshop) Microservices Architecture.

## Overview
This service handles user warehouse operations as usual create, read, update, and delete. Also handle product stock in each warehouse and it's movement. Movement will divided by two part, one is movement in (between warehouse) and other one is movement out (warehouse to customer). We use redis as database for warehouse ranking by zipcode and postgres for transactional database.

## Architecture
```
eshop-auth
├── .github/
│   └── workflows/     # github workflows to automatically test, build, and push
├── cmd/
│   └── app/            # configuration and log initialization
├── config/             # configuration
├── internal/   
│   ├── app/            # one run function in the `app.go`
│   ├── controller/     # serve handler layer
│   │   ├── http/
│   │   |   └── v1/     # rest http
│   │   └── kafka
│   │       └── v1/     # kafka subscriber
│   ├── entity/         # entities of business logic (models) can be used in any layer
│   ├── usecase/        # business logic
│   │   └── repo/       # abstract storage (database) that business logic works with
│   └── utils/          # helpers function
├── migrations/         # sql migration
└── pkg/
    ├── httpserver/     # http server initialization
    ├── kafka/          # kafka initialization
    ├── logger/         # logger initialization
    ├── postgresql/     # postgresql initialization
    └── redis/          # redis initialization
```

## Tech Stack
- Backend: Go
- Authorization: AWS Cognito
- Database: PostgreSQL and Redis
- CI/CD: Github Actions
- Message Broker: Apache Kafka
- Container: Docker

## API Documentation
tbd