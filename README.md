# Daily Worker Roster

A RESTful API for managing daily worker shifts

## Features
- User registration and login (JWT-based authentication)
- Admin and worker roles
- CRUD operations for users and shifts
- Shift request, approval, and assignment workflows
- API documentation with Swagger UI
- Containerized with Docker/Podman and MySQL

## Getting Started


### Build and Run with Docker Compose
```sh
docker-compose up --build
```
- The API will be available at `http://localhost:8080`
- MySQL will be available at `localhost:3306`

### 3. API Documentation
Visit: [http://localhost:8080/swagger/index.html]

### 4. Generate Swagger Docs
Install swag CLI
```sh
go install github.com/swaggo/swag/cmd/swag@latest
```
Generate docs:
```sh
make swag
```
