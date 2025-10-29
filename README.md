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
│ ├── config/
│ ├── dto/
│ │ └── user_dto.go
│ ├── entity/
│ │ └── user.go
│ ├── handler/
│ │ └── user_handler.go
│ ├── repository/
│ │ └── user_repository.go
│ ├── router/
│ │ └── router.go
│ └── usecase/
│ └── user_usecase.go
├── .env
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
└── README.md
```

### Folder Purpose

| Folder / File | Description |
|----------------|-------------|
| `internal/config` | Contains configuration setup (DB connection, environment variables, etc.) |
| `internal/dto` | Data Transfer Objects (for request/response mapping) |
| `internal/entity` | Entity models representing database tables |
| `internal/handler` | HTTP handlers (controllers) that handle requests |
| `internal/repository` | Database access layer (CRUD operations) |
| `internal/router` | Routing setup for all endpoints |
| `internal/usecase` | Business logic layer (service layer) |
| `docs` | Swagger-generated API documentation |
| `main.go` | Application entry point |
| `.env` | Environment variables configuration |
| `Dockerfile` | Container build configuration |
| `go.mod / go.sum` | Go dependencies and version management |

---

## Author
Junico Dwi Chandra  
junicodwi.chandra@gmail.com  
https://junicochandra.com  

## License
MIT License © 2025 — Created with by **Junico Dwi Chandra** using Go, GORM, and Docker.