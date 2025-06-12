# Hospital Management System

A hospital management system with two portals (receptionist and doctor) using Golang with Gin framework and PostgreSQL.

## Features

### Core Features

- Single login API for both portal types (returns JWT with role claim)
- Receptionist portal:
  - Patient CRUD operations (Create, Read, Update, Delete)
  - Patient search functionality
- Doctor portal:
  - View patient details
  - Update medical notes specifically
  - View patient medical history

### Additional Features

- Patient search with filters (name, age range, gender, contact info)
- Role-based access control
- JWT authentication
- Password hashing with bcrypt
- Input validation
- Swagger API documentation

## Tech Stack

- Go 1.22
- Gin Web Framework
- GORM (Go Object Relational Mapper)
- PostgreSQL
- JWT for authentication
- Swagger for API documentation

## Project Structure

The project follows clean architecture principles with the following structure:

- `cmd/server`: Application entry point
- `internal/config`: Configuration code
- `internal/controllers`: HTTP request handlers
- `internal/middleware`: HTTP middleware
- `internal/models`: Database models and DTOs
- `internal/repositories`: Data access layer
- `internal/services`: Business logic
- `internal/utils`: Utility functions
- `migrations`: Database migration scripts
- `tests`: Test files

## Setup

### Prerequisites

- Go 1.22 or higher
- PostgreSQL

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/hospital-project.git
cd hospital-project
```

2. Install dependencies:

```bash
go mod download
```

3. Set up the database:

```bash
# Create a PostgreSQL database
createdb hospital

# Run migrations
# You can use a migration tool like golang-migrate:
# migrate -database "postgresql://postgres:postgres@localhost:5432/hospital?sslmode=disable" -path migrations up
```

4. Configure environment variables:

Create a `.env` file in the root directory with the following content:

```
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=hospital

# JWT configuration
JWT_SECRET_KEY=your_secret_key_here

# Server configuration
PORT=8080
GIN_MODE=debug
```

### Running the Application

```bash
go run cmd/server/main.go
```

The server will start on port 8080 (or the port specified in the `.env` file).

## API Documentation

Swagger documentation is available at:

```
http://localhost:8080/swagger/index.html
```

## Default Users

The system comes with two default users:

1. Receptionist:
   - Username: admin
   - Password: admin123

2. Doctor:
   - Username: doctor
   - Password: doctor123

## API Endpoints

### Authentication

- `POST /api/auth/login`: Login with username and password

### Users

- `POST /api/users`: Create a new user (authenticated)
- `GET /api/users/:id`: Get a user by ID (authenticated)
- `GET /api/users/me`: Get the current authenticated user

### Patients (Receptionist)

- `POST /api/patients`: Create a new patient
- `GET /api/patients`: List all patients
- `GET /api/patients/:id`: Get a patient by ID
- `PUT /api/patients/:id`: Update a patient
- `DELETE /api/patients/:id`: Delete a patient
- `GET /api/patients/search`: Search for patients with filters

### Patients (Doctor)

- `GET /api/patients`: List all patients
- `GET /api/patients/:id`: Get a patient by ID
- `PUT /api/patients/:id/medical-notes`: Update a patient's medical notes

## License

This project is licensed under the MIT License - see the LICENSE file for details.
