# Go Commerce : E-Commerce Platform [ Product Service ]

## 📝 Project Overview
Go Commerce is an e-commerce system that provides user account management (buyer and seller roles), product management, and secure payment processing through the payment gateway

## 🎯 Problem Statement & Solution

#### Problem Statement
Many users and sellers face challenges managing their online sales and purchases due to difficulties in tracking products, handling orders, and processing payments securely. Manual processes or fragmented systems can lead to errors, delayed transactions, and poor user experience.

#### Solution
Go Commerce provides a comprehensive e-commerce platform that centralizes user account management, product management, and payment processing. The platform allows users to:

* Create and manage buyer or seller accounts efficiently
* Add, update, and manage product listings
* Process payments securely through integrated payment gateways
* Track orders and transactions seamlessly in one system

By centralizing these processes, Go Commerce improves operational efficiency, reduces errors, and enhances the overall online shopping experience for both buyers and sellers.

## 📚 Documentation
* [Go Commerce with REST API](https://github.com/fildzaanf/go-commerce-api)

## 🚀 Tools and Technologies 
* Go Programming Language
* Echo Go Framework
* GORM for Object Relational Mapping
* MySQL / PostgreSQL for Relational Database
* JSON Web Token (JWT) for Authentication
* Docker for Containerization
* Midtrans Payment Gateway integrated with Webhooks, SMTP, and GoMail for real-time payment notifications
* Amazon Web Services (AWS)
  * Amazon Simple Storage Service (S3)
* GRPC for efficient, low-latency, and strongly-typed API communication

## 🏛️ System Design and Architecture

* Clean Architecture
* Hexagonal Architecture
* Domain-Driven Design (DDD)
* Command Query Responsibility Segregation (CQRS)
* Microservices Architecture
* REST API
* GRPC
* Webhook

## ✨ Features

#### User Management

| Feature                   | Description                                                        |
| ------------------------- | ------------------------------------------------------------------ |
| User Registration & Login | Allows users to register and log in to access the platform         |
| Profile                   | Provides functionality to retrieve user profile information by ID  |

#### Product Management

| Feature           | Description                                                                    |
| ----------------- | ------------------------------------------------------------------------------ |
| Create Product    | Enables adding new products to the platform                                    |
| Update Product    | Allows updating existing product details by product ID                         |
| Delete Product    | Supports removing products from the platform by product ID                     |
| Retrieve Product  | Provides access to a single product by ID or a list of all available products  |

#### Payment Management

| Feature             | Description                                                                          |
| ------------------- | ------------------------------------------------------------------------------------ |
| Create Payment      | Allows users to create new payments for products or services                         |
| Retrieve Payment    | Provides access to all payments or details of a specific payment by ID               |
| Integration Payment | Supports real-time payment updates via Midtrans Webhook integration                  |
| Integration Email   | Integrated Midtrans Payment Gateway using Webhooks with SMTP and Go Mail to automate | 
|                     | event-driven email notifications based on real-time payment status updates           |

## 📡 gRPC Services

#### Products

| RPC Method     | RPC Type | Description             |
|---------------|----------|-------------------------|
| CreateProduct | Unary    | Create new product      |
| UpdateProduct | Unary    | Update product data     |
| DeleteProduct | Unary    | Delete product          |
| GetProductByID | Unary    | Retrieve product by ID      |
| GetAllProducts | Unary    | Retrieve all products       |




## 📡 API Endpoints

#### Products

| Method | Endpoint      | Description                    |
| ------ | ------------- | ------------------------------ |
| POST   | /products     | Create a new product           |
| PUT    | /products/:id | Update product details by ID   |
| DELETE | /products/:id | Delete product by ID           |
| GET    | /products/:id | Retrieve product details by ID |
| GET    | /products     | Retrieve all products          |


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

