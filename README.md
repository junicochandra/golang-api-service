# Golang API Service

A modern, clean, and scalable RESTful API designed for performance, maintainability, and ease of deployment.
It follows a modular architecture, promotes best practices for backend development, and provides built-in support for database integration, API documentation, and containerized environments.

---

## Tech Stack
- Language: **Go (1.24+)**
- Framework: **Gin**
- ORM: **GORM**
- Database: **MySQL 8**
- Doc: **Swaggo (Swagger)**
- Containerization: **Docker**

## Features
- Built using **Gin** — high-performance HTTP framework for Go  
- **Clean Architecture** for maintainable and modular code  
- **MySQL Integration** using **GORM ORM** (object-relational mapping)  
- **Swagger** auto-generated API documentation  
- **Password hashing** with bcrypt for secure user management  
- **Docker-ready** setup for local or production environments  
- Ready for **scaling into microservices**

## Run with Docker

Easily spin up the API and MySQL database using the [Docker Starterpack](https://github.com/junicochandra/docker-starterpack).

### Start the service
- API : http://localhost:9000
- Swagger : http://localhost:9000/api/doc/index.html


## Project Structure

```bash
├── docs/
├── internal/
│   ├── app/                        # Application layer (use cases & interfaces)
│   │   ├── auth/
│   │   │   ├── dto/
│   │   │   │   └── auth_dto.go
│   │   │   ├── auth_interface.go
│   │   │   └── auth_usecase.go
│   │   └── user/
│   │       ├── dto/
│   │       │   └── user_dto.go
│   │       ├── user_interface.go
│   │       └── user_usecase.go
│   │
│   ├── domain/                     # Core business entities & repository interfaces
│   │   ├── entity/
│   │   │   └── user.go
│   │   └── repository/
│   │       └── user_repository_interface.go
│   │
│   ├── handler/                    # HTTP handlers (controllers)
│   │   ├── auth_handler.go
│   │   ├── profile_handler.go
│   │   └── user_handler.go
│   │
│   ├── infrastructure/             # External frameworks and drivers
│   │   ├── config/
│   │   │   └── database/
│   │   │       └── mysql.go
│   │   ├── repository/             # Implementation of repository interfaces
│   │   └── service/
│   │
│   ├── middleware/
│   │
│   └── router/
│       └── router.go
├── .env
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
└── README.md
```

### Folder Purpose

| Folder / File                        | Description                                                                                                            |
| ------------------------------------ | ---------------------------------------------------------------------------------------------------------------------- |
| `internal/app`                       | Contains **application logic**, such as use cases and interfaces for each module (e.g., `auth`, `user`).               |
| `internal/app/<module>/dto`          | **Data Transfer Objects (DTOs)** — define request/response structures used between handler and use case.               |
| `internal/domain/entity`             | Core **domain models** representing your business entities (e.g., `User`, `Auth`).                                     |
| `internal/domain/repository`         | Contains **repository interfaces** that define contracts for data access (implemented in `infrastructure/repository`). |
| `internal/infrastructure/config`     | Configuration setup (e.g., **database connection**, environment variables).                                            |
| `internal/infrastructure/repository` | **Repository implementations** — concrete structs that interact with the database via GORM.                            |
| `internal/infrastructure/middleware` | Custom **Gin middleware** such as JWT authentication.                                                                  |
| `internal/handler`                   | **HTTP handlers (controllers)** — handle API requests and responses, call use cases.                                   |
| `internal/router`                    | Central **route definitions**, wiring handlers, middleware, and API groups.                                            |
| `docs`                               | **Swagger-generated API documentation** (auto-created by Swaggo).                                                      |
| `.env`                               | Environment variable configuration file.                                                                               |
| `Dockerfile`                         | Docker container build configuration for deployment.                                                                   |
| `go.mod / go.sum`                    | Go dependencies and module version management.                                                                         |
| `main.go`                            | Application entry point — initializes app, loads configs, and starts the server.                                       |


---

## Author
Junico Dwi Chandra  
junicodwi.chandra@gmail.com  
https://junicochandra.com  

## License
MIT License © 2025 — Created with by **Junico Dwi Chandra** using Go, GORM, and Docker.