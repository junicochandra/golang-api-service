# Golang API Service

A modern, clean, and scalable RESTful API designed for performance, maintainability, and ease of deployment.
It follows a modular architecture, promotes best practices for backend development, and provides built-in support for database integration, message queue processing, API documentation, and containerized environments.

---

## Tech Stack
- Language: **Go (1.24+)**
- Framework: **Gin**
- ORM: **GORM**
- Database: **MySQL 8**
- Message Queue: **RabbitMQ**
- Doc: **Swaggo (Swagger)**
- Containerization: **Docker**

## Features
- Built using **Gin** — high-performance HTTP framework for Go  
- **Clean Architecture** for maintainable and modular code  
- **MySQL Integration** using **GORM ORM** (object-relational mapping)  
- **RabbitMQ message queue** for async job processing (TopUp, Payment, etc.)
- **Swagger** auto-generated API documentation  
- **Password hashing** with bcrypt for secure user management  
- **Docker-ready** setup for local or production environments  
- Ready for **scaling into microservices**

## Run with Docker

Easily spin up the API and MySQL database using the [Docker Starterpack](https://github.com/junicochandra/docker-starterpack).

### Start the service
- Swagger : http://localhost:9000/api/doc/index.html


## Project Structure

```bash
├golang-api-service/
│
├── docs/
├── internal/
│   ├── app/
│   │   ├── auth/
│   │   ├── messaging/
│   │   ├── payment/
│   │   └── user/
│   │
│   ├── bootstrap/
│   ├── domain/
│   │   ├── entity/
│   │   └── repository/
│   │
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── payment_handler.go
│   │   ├── profile_handler.go
│   │   └── user_handler.go
│   │
│   ├── infrastructure/
│   │   ├── config/
│   │   ├── repository/
│   │   └── service/
│   │       └── rabbitmq/
│   │            ├── worker/
│   │            ├── rabbitmq.go
│   │            └── setup.go
│   │
│   ├── middleware/
│   └── router/
│
├── .env
├── Dockerfile
├── go.mod
├── go.sum
└── main.go
```

### Folder Purpose

| Folder / File                        | Description                                                                                                                                                                                                             |
| ------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `internal/app`                       | Contains **application logic**, such as use cases and interfaces for each module (e.g., `auth`, `user`).                                                                                                                |
| `internal/app/<module>/dto`          | **Data Transfer Objects (DTOs)** — define request/response structures used between handler and use case.                                                                                                                |
| `internal/domain/entity`             | Core **domain models** representing your business entities (e.g., `User`, `Auth`).                                                                                                                                      |
| `internal/domain/repository`         | Contains **repository interfaces** that define contracts for data access (implemented in `infrastructure/repository`).                                                                                                  |
| `internal/infrastructure/config`     | Configuration setup (e.g., **database connection**, environment variables).                                                                                                                                             |
| `internal/infrastructure/repository` | **Repository implementations** — concrete structs that interact with the database via GORM.                                                                                                                             |
| `internal/infrastructure/service`    | Utility and **shared infrastructure services** such as:<br> • `hash_service.go` — handles password hashing and verification using bcrypt.<br> • `jwt_service.go` — manages JWT token creation, signing, and validation. <br> • `rabbitmq.go` — provides RabbitMQ connection handling, channel management, and message publishing. <br> • `worker/` — contains background consumer workers that listen to specific queues (e.g., top-up queue) and process asynchronous jobs. |
| `internal/middleware`                | Custom **Gin middleware**, e.g., JWT authentication, logging, or request validation.                                                                                                                                    |
| `internal/handler`                   | **HTTP handlers (controllers)** — handle API requests and responses, calling the appropriate use case.                                                                                                                  |
| `internal/router`                    | Central **route definitions**, wiring handlers, middleware, and API groups.                                                                                                                                             |
| `docs`                               | **Swagger-generated API documentation** (auto-created by Swaggo).                                                                                                                                                       |
| `.env`                               | Environment variable configuration file.                                                                                                                                                                                |
| `Dockerfile`                         | Docker container build configuration for deployment.                                                                                                                                                                    |
| `go.mod / go.sum`                    | Go dependencies and module version management.                                                                                                                                                                          |
| `main.go`                            | Application entry point — initializes app, loads configs, and starts the server.                                                                                                                                        |

## Author
Junico Dwi Chandra  
junicodwi.chandra@gmail.com  
https://junicochandra.com  

## License
MIT License © 2025 — Created by **Junico Dwi Chandra**, powered by Go and modern cloud-ready tooling.