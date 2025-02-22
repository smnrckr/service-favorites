# Cimri Internship Case - Favorites Service

This repository contains the **Favorites Service**, which is a part of a **microservices-based architecture** developed in Go. The system consists of the following core services:

- **User Service**: Handles user-related operations, including managing user profiles, preferences, and interactions.
- **Product Service**: Manages product data and retrieves product details based on product ID.
- **Favorites Service**: Enables users to create, update, delete, and view their favorite lists. Users can also add or remove products from these lists. This service communicates with both the **User Service** and the **Product Service**. 

Within the **Favorites Service**, the operations are separated into two main functionalities:

1. **Favorites Service**: 
   - Users can add products to their favorite list, remove products, and view the products they have favorited.
   - This service communicates with the **User Service** and **Product Service** to manage the products in users' favorites.

2. **Favorite Lists Service**: 
   - Enables users to create and manage multiple favorite lists.
   - Users can organize their favorite products into different lists, making it easier to categorize and view their preferences.

The **Favorites Service** is built with **Go**, **Fiber** (for HTTP routing), and **GORM** (for ORM database interactions). It uses **PostgreSQL** as the database and **Redis** for caching. The service is Dockerized for easy deployment and integrates with other microservices via HTTP clients.

## 📌 Technologies Used

- **Go (Golang)**
- **Fiber** (Web Framework)
- **GORM** (Object Relational Mapping)
- **PostgreSQL**
- **Redis**
- **Docker & Docker Compose**
- **AWS EC2, RDS, S3**
- **Swagger / OpenAPI**

## 📂 Project Structure

```
/cmd
  /main.go          # Entry point for the HTTP server
/internal
  /handler         # Handles HTTP requests and responses for Favorites and Favorite Lists services
  /service         # Contains business logic for managing favorites and lists
  /repository      # Handles database interactions related to favorites and lists (using GORM)
  /models          # Defines data models (using GORM)
  /clients         # Contains HTTP clients to interact with User Service and Product Service
/utils
  /envloader       # Loads environment variables
/pkg
  /redis           # Handles Redis connection
  /s3              # Handles S3 connection
  /postgres        # Handles PostgreSQL connection (with GORM)
```
## The Favorites Service includes two important HTTP clients that communicate with the User Service and Product Service:

### UserClient:
This client is responsible for interacting with the User Service. It allows the Favorites Service to check whether a user exists by querying the User Service with the user ID.
If the user exists, the Favorites Service can proceed to associate favorite products with that user.
### ProductClient:
This client is used for interacting with the Product Service. It enables the Favorites Service to fetch product details based on product IDs.
The ProductClient retrieves information about products that a user may add to their favorites, providing necessary product details such as name, description, price, and more.

### Key Points:

- **Fiber**: The application uses **Fiber** as the web framework to handle routing and HTTP requests efficiently.
- **GORM**: **GORM** is used for ORM-based database interactions with **PostgreSQL**.
- **Tests**: Unit tests for **handlers** and **repository** are written directly within the respective packages (e.g., `handler/favorites_handler_test.go`, `repository/favorites_repository_test.go`).
- **`/utils`**: Contains the **env loader**.
- **`/pkg`**: Contains utilities to handle connections to **Redis**, **S3**, and **PostgreSQL**.
- **Main Function**: The **main** function is located in the `/cmd` directory, where the application is initialized with the environment variables and client connections to the **User Service** and **Product Service**.

## 🚀 Getting Started

### 1️⃣ Install Dependencies

Ensure Go modules are up to date:

```sh
go mod tidy
```

### 2️⃣ Run with Docker

To start the service along with PostgreSQL and Redis:

```sh
docker-compose up --build
```

### 3️⃣ Run Manually (Without Docker)

If you prefer to run the service manually:

```sh
go run cmd/server/main.go
```

Make sure to configure your `.env` file correctly before running the service.

## ✅ Running Unit Tests

This repository includes unit tests for the repository and handler layers. To execute all tests:

```sh
go test ./...
```

To run tests for a specific package:

```sh
go test ./internal/handler
```

Unit tests utilize mock data, so no real database connection is required.

## 📖 API Documentation

API endpoints are documented using Swagger/OpenAPI. Once the service is running, access the API documentation at:

```
http://localhost:8081/swagger
```
