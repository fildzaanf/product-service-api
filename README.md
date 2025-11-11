# product-service-api

## 📌 Overview
This project is a simple e-commerce system that includes key features such as user management (sellers and buyers), product catalog, and payment system using Midtrans payment gateway.

## 🚀 Tools and Technologies 
- Go Programming Language
- Echo Framework
- GORM
- MySQL / PostgreSQL
- Docker
- JWT Authentication
- Midtrans Payment Gateway
- Amazon Simple Storage Service (S3)
- Simple Mail Transfer Protocol (SMTP)
- gRPC (Remote Procedure Call)

## 🏛️ System Design and Architecture
- Hexagonal Architecture
- Clean Architecture
- Domain-Driven Design (DDD)
- Command Query Responsibility Segregation (CQRS)

## 📂 Folder Structure
```

product-service-api/
├── cmd/                                        # Entry points for starting the application
│   ├── grpc/                                   # Main gRPC server initialization Module
│   │   └── grpc_server.go                      
│   └── rest/                                   # Main REST server initialization Module
│       └── rest_server.go                      
│
├── docs/                                       # System Documentation (API docs, diagrams, specs, notes)
│
├── internal/                                   # Internal application modules (non-exported Go modules)
│   ├── product/                                # Product domain module
│   │   ├── adapter/                            # Adapters bridging external (HTTP, gRPC, DB, Message Broker, External Service) I/O to core logic (domain + application)
│   │   │   ├── client/                         # Clients to communicate with external/internal services
│   │   │   │   ├── pb/                         # Protobuf-generated files for external/internal service clients
│   │   │   │   │   ├── user_grpc.pb.go         # gRPC client stub for UserService
│   │   │   │   │   ├── user.pb.go              # Protobuf message definitions for User Client
│   │   │   │   │   └── user.proto              # Proto schema for User Client 
│   │   │   │   └── user_client_grpc.go         # Implementation of User gRPC client logic
│   │   │   ├── handler/                        # Inbound Delivery layer (REST & gRPC handlers routing external requests into the application)
│   │   │   │   ├── grpc/                       # gRPC handlers
│   │   │   │   │   ├── pb/                     # Protobuf-generated files for ProductService
│   │   │   │   │   │   ├── product_grpc.pb.go  # gRPC server and client stub for ProductService
│   │   │   │   │   │   ├── product.pb.go       # Protobuf message definitions for Product Domain
│   │   │   │   │   │   ├── product.proto       # Proto schema for Product Domain
│   │   │   │   │   │   └── proto.go            # Mapper for converting protobuf messages requests/responses to domain entities
│   │   │   │   │   ├── command_handler.go      # gRPC command handlers for modifying product data (write operations ex: Create, Update, Delete)
│   │   │   │   │   ├── query_handler.go        # gRPC query handlers for retrieving product data (read operations ex: Get/List)
│   │   │   │   │   └── service.go              # Handles dependency injection and gRPC service registration for the Product Domain
│   │   │   │   └── rest/                       # REST handlers
│   │   │   │       ├── command_handler.go      # REST command handlers for modifying product data (write operations ex: Create, Update, Delete)
│   │   │   │       ├── query_handler.go        # REST query handlers for retrieving product data (read operations ex: Get/List)
│   │   │   │       ├── json.go                 # DTO definitions and mappers for converting JSON messages requests/responses to domain entities
│   │   │   │       └── router.go               # Handles dependency injection and route configuration for the Product Domain
│   │   │   ├── model/                          # ORM model definition for Product using GORM
│   │   │   │   └── gorm_model.go               
│   │   │   └── repository/                     # Outbound Repository layer (Database operations used by the Application)
│   │   │       ├── gorm/                       # Repository implementations using GORM ORM
│   │   │       │   ├── command_repository.go   # GORM-based repository for modifying product data (write operations ex: Create, Update, Delete)
│   │   │       │   └── query_repository.go     # GORM-based repository for retrieving product data (read operations ex: Get/List)
│   │   │       └── postgresql/                 # Repositories using raw SQL/PostgreSQL driver
│   │   │           ├── command_repository.go   # Raw SQL-based repository for modifying product data (write operations ex: Create, Update, Delete)
│   │   │           └── query_repository.go     # Raw SQL-based repository for retrieving product data (read operations ex: Get/List)
│   │   ├── application/                        # Application service layer (business processes)
│   │   │   ├── port/                           # Ports (interfaces) used by Adapters and Application
│   │   │   │   ├── external_interface.go       # Interfaces for external clients 
│   │   │   │   ├── repository_interface.go     # Interface for repositories
│   │   │   │   └── service_interface.go        # Interface for services
│   │   │   └── service/                        # Application service implementations
│   │   │       ├── command_service.go          # Business logic for modifying product data (write operations ex: Create, Update, Delete)
│   │   │       └── query_service.go            # Business logic for retrieving product data (read operations ex: Get/List)
│   │   └── domain/
│   │       └── entity.go                       # Core domain entities and domain rules (DDD)
│
├── infrastructure/                             # System infrastructure & integrations
│   ├── cloud/                                  # Cloud service integrations (S3 buckets, GCS, etc.)
│   ├── config/                                 # Application configuration (YAML/ENV loaders)
│   ├── database/                               # Database initialization, pooling, migration utilities
│   └── email/                                  # Email sending adapters (SMTP, providers, templates)
│
├── pkg/                                        # Shared utility packages (global helpers)
│   ├── constant/                               # Constant definitions 
│   ├── crypto/                                 # Cryptography (hashing, encryption)
│   ├── generator/                              # ID generators, token generators, etc.
│   ├── middleware/                             # HTTP/gRPC middleware (auth, logging, interceptor, etc.)
│   ├── response/                               # Standardized API response formatting
│   └── validator/                              # Validation utilities
│
├── .env                                        # Local environment configuration for development
├── .gitignore                                  # Git ignore rules
├── go.mod                                      # Go module definition
├── go.sum                                      # Go dependency checksums
└── README.md                                   # Main project documentation
```
## 🛠️ Installation & Running the Project
### 1️⃣ Prerequisites
Make sure you have installed:
- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/download/) / [MySQL](https://dev.mysql.com/downloads/)
- [Midtrans](https://midtrans.com/)

### 2️⃣ Clone the Repository
```bash
git clone <repo-url>
cd <repository-root-directory>
```

### 3️⃣ Configure Environment
Create a `.env` file based on `.env.example` and place it in the root directory.

### 4️⃣ Run the Application
```bash
go run cmd/grpc/grpc_server.go
go run cmd/rest/rest_server.go
```






