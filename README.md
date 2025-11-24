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

## Database
````md
-- golang_api.accounts definition
CREATE TABLE `accounts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `account_number` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL,
  `balance` decimal(18,2) DEFAULT '0.00',
  `currency` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT 'IDR',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `account_number` (`account_number`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=99 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- golang_api.transaction_logs definition
CREATE TABLE `transaction_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `transaction_id` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` enum('pending','processing','success','failed') COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `message` text COLLATE utf8mb4_unicode_ci,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `transaction_id` (`transaction_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- golang_api.transactions definition
CREATE TABLE `transactions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `transaction_id` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` enum('transfer','topup','payment') COLLATE utf8mb4_unicode_ci NOT NULL,
  `sender_account_id` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `receiver_account_id` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `amount` decimal(18,2) NOT NULL,
  `status` enum('pending','processing','completed','success','failed') CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'pending',
  `reference` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `payload` text COLLATE utf8mb4_unicode_ci,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `transaction_id` (`transaction_id`),
  KEY `sender_account_id` (`sender_account_id`),
  KEY `receiver_account_id` (`receiver_account_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- golang_api.users definition
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `staff_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `contact` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email_verified_at` datetime(3) DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `remember_token` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `role` tinyint DEFAULT NULL,
  `photo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `office` bigint DEFAULT NULL,
  `last_login_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `deleted_by` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_email` (`email`),
  UNIQUE KEY `uni_users_username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
````  

## Author
Junico Dwi Chandra  
junicodwi.chandra@gmail.com  
https://junicochandra.com  

## License
MIT License © 2025 — Created by **Junico Dwi Chandra**, powered by Go and modern cloud-ready tooling.