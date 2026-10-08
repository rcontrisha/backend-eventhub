# EVENTHUB API 

![](https://img.shields.io/badge/Go-1.27.1-grey?style=for-the-badge&logo=go&logoColor=white&labelColor=007ACC)
![](https://img.shields.io/badge/Gin-1.12.0-grey?style=for-the-badge&logo=gin&logoColor=white&labelColor=008ECF)
![](https://img.shields.io/badge/Postgresql-16.15-grey?style=for-the-badge&logo=postgresql&logoColor=white&labelColor=4169E1)
![](https://img.shields.io/badge/redis-8.10.2-grey?style=for-the-badge&logo=redis&logoColor=white&labelColor=FF4438)

A Go-based event management API built with Gin, PostgreSQL, and Redis. It provides a scalable backend for creating, retrieving, and managing event-related data.

## Technologies
- **Go 1.27.1** — A statically typed programming language used to build the API and manage its backend services.
- **Gin 1.12.0** — A lightweight HTTP web framework used for routing and handling API requests.
- **PostgreSQL 16.15** — A relational database used to store and manage event-related data.
- **Redis 8.10.2** — An in-memory data store used for caching and fast data operations.
- **Swagger** — Generating interactive API documentation and testing endpoints from the API's OpenAPI specification.

## Features
- **Create and manage events** — Build a reliable API for creating and maintaining event-related data.
- **Retrieve event data** — Query stored event information through the API.
- **Secure authentication** — Protects API endpoints with authentication middleware for authorized access.
- **Scalable service architecture** — Uses Gin for request handling, PostgreSQL for persistent storage, and Redis for fast data operations.
- **Cached data operations** — Uses Redis to improve response speed and reduce repeated work for data access.
- **Interactive API documentation** — Swagger exposes the OpenAPI specification, making endpoints easy to discover and test.

## Getting Started

### Prerequisites

Install the following tools before running the project with Docker Compose:

- **Docker**
- **Docker Compose**
- **Git**

### Clone the repository

```bash
git clone https://github.com/rcontrisha/backend-eventhub.git
cd backend-eventhub
```

### Configure the environment

Create a `.env` file in the `backend-eventhub` directory and add the variables required by the backend and database services:

```env
PORT=8080
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=eventhub
DB_HOST=database
DB_PORT=5432
REDIS_USER=redis
REDIS_PASSWORD=redis
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=replace-with-a-long-random-secret
```

The Docker Compose configuration overrides the database and Redis hosts and ports so the backend connects to the service names `database` and `redis` on the Docker network.

Create a `.env` file in the `database` directory and configure the PostgreSQL credentials used by the database container:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=eventhub
```

### Build and start the services

From the directory containing the Docker Compose configuration, build the images and start the backend, PostgreSQL, and Redis containers:

```bash
docker compose up -d
```

The `backend` service listens on port `8080`, PostgreSQL is exposed on port `7777`, and Redis is exposed on port `6666`.

### Verify the services

Check that the containers are running:

```bash
docker compose ps
```

Confirm that PostgreSQL is ready:

```bash
docker compose logs database
```

The database container should report that PostgreSQL is accepting connections. The backend will wait for the database health check before starting.

### Run the API

The API is available at `http://localhost:8080`.

To stop the services, run:

```bash
docker compose down
```

To stop the services and remove their local data volumes, run:

```bash
docker compose down -v
```

### API documentation

Swagger documentation is available at:

- `http://localhost:8080/swagger/index.html`

## License
Distributed under the MIT License. See [License](/LICENSE) for more information.