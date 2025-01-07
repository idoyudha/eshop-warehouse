# Warehouse Service
Part of [eshop](https://github.com/idoyudha/eshop) Microservices Architecture.

## Overview
This service handles user warehouse operations as usual create, read, update, and delete. Also handle product stock in each warehouse and it's movement. Movement will divided by two part, one is movement in (between warehouse) and other one is movement out (warehouse to customer). We use redis as database for warehouse ranking by zipcode and postgres for transactional database.

## Architecture
```
eshop-auth
├── .github/
│   └── workflows/
├── cmd/
│   └── app/
├── config/
├── internal/   
│   ├── app/
│   ├── controller/
│   │   ├── http/
│   │   |   └── v1/
│   │   └── kafka
│   │       └── v1/
│   ├── entity/
│   ├── usecase/
│   │   └── repo/
│   └── utils/
├── migrations/
└── pkg/
    ├── httpserver/
    ├── kafka/
    ├── logger/
    ├── postgresql/
    └── redis/
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