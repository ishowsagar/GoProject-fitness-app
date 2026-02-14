# 🏋️ FitTrack API

A production-ready RESTful API for workout tracking built with Go, PostgreSQL, and Docker. Features JWT authentication, complete CRUD operations, and a robust database design.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-12.4-316192?style=flat-square&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

## 🚀 Features

- **JWT Authentication** - Secure token-based authentication system
- **User Management** - Registration with email validation and password hashing (bcrypt)
- **Workout CRUD** - Full create, read, update, delete operations for workouts
- **Authorization** - Users can only modify their own workouts
- **Database Migrations** - Automated schema versioning with Goose
- **Docker Support** - Complete containerization with Docker Compose
- **Live Reload** - Air integration for hot reload during development
- **RESTful Design** - Clean, intuitive API endpoints
- **Error Handling** - Comprehensive error responses with proper HTTP status codes
- **Logging** - Request logging and error tracking

## 📋 Table of Contents

- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [API Endpoints](#-api-endpoints)
- [Getting Started](#-getting-started)
- [Database Schema](#-database-schema)
- [Authentication Flow](#-authentication-flow)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Contributing](#-contributing)

## 🛠️ Tech Stack

### Backend

- **Go 1.25** - Primary programming language
- **Chi Router** - Lightweight, fast HTTP router
- **PostgreSQL** - Relational database
- **pgx** - PostgreSQL driver and toolkit
- **Goose** - Database migration tool
- **JWT** - JSON Web Tokens for authentication
- **Bcrypt** - Password hashing

### DevOps

- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **Air** - Live reload for Go applications

## 📁 Project Structure

```
GoProject#4/
├── internal/
│   ├── api/                    # HTTP handlers
│   │   ├── user_handler.go     # User registration
│   │   ├── workout_handler.go  # Workout CRUD operations
│   │   └── token_handler.go    # Authentication
│   ├── store/                  # Database layer
│   │   ├── database.go         # DB connection
│   │   ├── user_store.go       # User queries
│   │   ├── workout_store.go    # Workout queries
│   │   └── tokens.go           # Token management
│   ├── middleware/             # HTTP middleware
│   │   └── middleware.go       # Authentication middleware
│   ├── routes/                 # Route definitions
│   │   └── routes.go           # API routes
│   ├── tokens/                 # JWT logic
│   │   └── tokens.go           # Token generation
│   ├── utils/                  # Helper functions
│   │   └── utils.go            # JSON helpers
│   └── app/                    # Application setup
│       └── app.go              # App initialization
├── migrations/                 # SQL migrations
│   ├── 00001_users.sql
│   ├── 00002_workouts.sql
│   ├── 00003_workout_entries.sql
│   ├── 00004_tokens.sql
│   └── 00005_user_id_alter.sql
├── docker-compose.yml          # Docker services
├── Dockerfile                  # Container build
├── .air.toml                   # Live reload config
└── main.go                     # Entry point
```

## 🔌 API Endpoints

### Public Endpoints

| Method | Endpoint                 | Description            | Request Body                                      |
| ------ | ------------------------ | ---------------------- | ------------------------------------------------- |
| `GET`  | `/health`                | Health check           | -                                                 |
| `POST` | `/users`                 | Register new user      | `username`, `email`, `password`, `bio` (optional) |
| `POST` | `/tokens/authentication` | Login / Get auth token | `username`, `password`                            |

### Protected Endpoints (Require Authentication)

| Method   | Endpoint         | Description          | Request Body                                                  |
| -------- | ---------------- | -------------------- | ------------------------------------------------------------- |
| `GET`    | `/workouts/{id}` | Get specific workout | -                                                             |
| `POST`   | `/workouts`      | Create new workout   | `title`, `description`, `duration_minutes`, `calories_burned` |
| `PUT`    | `/workouts/{id}` | Update workout       | Same as POST (all fields optional)                            |
| `DELETE` | `/workouts/{id}` | Delete workout       | -                                                             |

### Example Requests

#### Register User

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepass123",
    "bio": "Fitness enthusiast"
  }'
```

#### Login

```bash
curl -X POST http://localhost:8080/tokens/authentication \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepass123"
  }'
```

#### Create Workout

```bash
curl -X POST http://localhost:8080/workouts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "Morning Run",
    "description": "5K run in the park",
    "duration_minutes": 30,
    "calories_burned": 250
  }'
```

#### Get Workout

```bash
curl -X GET http://localhost:8080/workouts/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

#### Update Workout

```bash
curl -X PUT http://localhost:8080/workouts/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "Evening Run",
    "duration_minutes": 45
  }'
```

#### Delete Workout

```bash
curl -X DELETE http://localhost:8080/workouts/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 🚀 Getting Started

### Prerequisites

- **Docker** & **Docker Compose** (recommended)
- OR **Go 1.25+** and **PostgreSQL** (for local development)

### Quick Start with Docker (Recommended)

1. **Clone the repository**

   ```bash
   git clone <your-repo-url>
   cd GoProject#4
   ```

2. **Start all services**

   ```bash
   docker-compose up -d
   ```

3. **Verify everything is running**

   ```bash
   docker-compose ps
   ```

4. **Test the API**
   ```bash
   curl http://localhost:8080/health
   ```

That's it! Your API is running at `http://localhost:8080` 🎉

### Manual Setup (Without Docker)

1. **Install Go 1.25+**

   ```bash
   go version
   ```

2. **Install PostgreSQL**
   - Create database: `fitness_db`

3. **Set environment variables** (optional, defaults to localhost:5445)

   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=fitness_db
   ```

4. **Install dependencies**

   ```bash
   go mod download
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on port 8080, and migrations will run automatically.

## 🗄️ Database Schema

### Users Table

```sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash BYTEA NOT NULL,
  bio TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Workouts Table

```sql
CREATE TABLE workouts (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  duration_minutes INTEGER NOT NULL,
  calories_burned INTEGER,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Tokens Table

```sql
CREATE TABLE tokens (
  hash BYTEA PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  expiry TIMESTAMP WITH TIME ZONE NOT NULL,
  scope TEXT NOT NULL
);
```

### Entity Relationship Diagram

```
┌─────────────┐         ┌──────────────┐         ┌──────────────┐
│    Users    │         │   Workouts   │         │    Tokens    │
├─────────────┤         ├──────────────┤         ├──────────────┤
│ id (PK)     │1      ∞ │ id (PK)      │         │ hash (PK)    │
│ username    ├─────────┤ user_id (FK) │    ┌────┤ user_id (FK) │
│ email       │         │ title        │    │    │ expiry       │
│ pass_hash   │1      ∞ │ description  │    │    │ scope        │
│ bio         ├─────────┤ duration     │    │    └──────────────┘
│ created_at  │         │ calories     │    │
└─────────────┘         │ created_at   │    │
                        │ updated_at   │    │
                        └──────────────┘    │
                                ↑           │
                                └───────────┘
```

## 🔐 Authentication Flow

1. **User Registration**
   - Client sends username, email, password
   - Server validates input
   - Password is hashed with bcrypt (cost factor 12)
   - User is stored in database
   - Returns user data (without password)

2. **User Login**
   - Client sends username and password
   - Server looks up user by username
   - Password is compared with stored hash using bcrypt
   - If valid, JWT token is generated (expires in 24 hours)
   - Token is returned to client

3. **Accessing Protected Routes**
   - Client includes token in `Authorization: Bearer <token>` header
   - Server validates token
   - If valid, user is fetched from database
   - User is added to request context
   - Handler processes request with authenticated user

4. **Authorization Checks**
   - For UPDATE/DELETE operations, server verifies ownership
   - Queries `workout.user_id` and compares with authenticated user
   - Returns 403 Forbidden if user doesn't own the resource

## 💻 Development

### With Docker (Recommended)

The project includes Air for live reload in development mode:

```bash
# Start development environment
docker-compose up -d

# View logs (watch for auto-reload)
docker-compose logs -f app

# Make changes to any .go file, then restart
docker-compose restart app
```

### Local Development with Air

For the best development experience on Windows:

```bash
# Install Air
go install github.com/air-verse/air@latest

# Start only databases
docker-compose up -d db test_db

# Run with Air (auto-reload on file changes)
air
```

Now any changes to `.go` files will automatically rebuild and restart the app!

### Project Commands

```bash
# Build the application
go build -o bin/fittrack main.go

# Run tests
go test -v ./...

# Run specific test
go test -v ./internal/store -run TestCreateWorkout

# Format code
go fmt ./...

# Lint code
go vet ./...

# Clean build artifacts
rm -rf bin/ tmp/
```

### Docker Commands

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f app

# Restart app
docker-compose restart app

# Rebuild after code changes
docker-compose build app
docker-compose up -d app

# Clean rebuild
docker-compose down
docker-compose build --no-cache
docker-compose up -d

# Access database
docker-compose exec db psql -U postgres
```

## 🧪 Testing

### Run All Tests

```bash
go test -v ./...
```

### Run Specific Package Tests

```bash
go test -v ./internal/store
```

### Test Coverage

```bash
go test -cover ./...
```

### Integration Tests

The project includes integration tests for the database layer. Make sure the test database is running:

```bash
docker-compose up -d test_db
go test -v ./internal/store
```

## 📦 Deployment

### Production Docker Build

```bash
# Build production image
docker build --target production -t fittrack:prod .

# Run production container
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PORT=5432 \
  -e DB_USER=your-user \
  -e DB_PASSWORD=your-secure-password \
  -e DB_NAME=fitness_db \
  fittrack:prod
```

### Environment Variables

| Variable      | Default     | Description       |
| ------------- | ----------- | ----------------- |
| `DB_HOST`     | `localhost` | Database host     |
| `DB_PORT`     | `5445`      | Database port     |
| `DB_USER`     | `postgres`  | Database user     |
| `DB_PASSWORD` | `postgres`  | Database password |
| `DB_NAME`     | `postgres`  | Database name     |

### Deployment Checklist

- [ ] Set strong database password
- [ ] Configure environment variables
- [ ] Set up SSL/TLS for database connection
- [ ] Configure reverse proxy (e.g., Nginx)
- [ ] Set up HTTPS with Let's Encrypt
- [ ] Configure CORS for your domain
- [ ] Set up monitoring and logging
- [ ] Configure database backups
- [ ] Set up CI/CD pipeline

## 🏗️ Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│           HTTP Layer (Chi)              │
│    (routes, middleware, handlers)       │
├─────────────────────────────────────────┤
│         Business Logic Layer            │
│   (validation, authorization, auth)     │
├─────────────────────────────────────────┤
│          Database Layer (Store)         │
│      (queries, transactions, CRUD)      │
├─────────────────────────────────────────┤
│         PostgreSQL Database             │
│     (data persistence, constraints)     │
└─────────────────────────────────────────┘
```

### Request Flow

```
Client Request
    ↓
Chi Router
    ↓
Middleware (CORS, Logging)
    ↓
Authentication Middleware
    ↓
Handler (validate input)
    ↓
Store (database operations)
    ↓
Database (PostgreSQL)
    ↓
Response (JSON)
```

## 🔒 Security Features

- ✅ **Password Hashing** - Bcrypt with cost factor 12
- ✅ **JWT Tokens** - Secure authentication tokens with expiry
- ✅ **SQL Injection Protection** - Parameterized queries
- ✅ **Authorization Checks** - Resource ownership verification
- ✅ **Input Validation** - Server-side validation for all inputs
- ✅ **Error Handling** - No sensitive data in error responses
- ✅ **HTTPS Ready** - Configure with reverse proxy

## 📊 Performance

- **Database Connection Pooling** - Efficient connection management
- **Indexed Queries** - Primary and foreign keys for fast lookups
- **Stateless Authentication** - JWT tokens for horizontal scaling
- **Lightweight Router** - Chi router with minimal overhead

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go conventions and idioms
- Add comments for exported functions
- Write tests for new features
- Format code with `go fmt`
- Check with `go vet` and `golint`

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

**Your Name**

- GitHub: [@yourusername](https://github.com/yourusername)
- LinkedIn: [Your LinkedIn](https://linkedin.com/in/yourprofile)

## 🙏 Acknowledgments

- [Chi Router](https://github.com/go-chi/chi) - Lightweight HTTP router
- [Goose](https://github.com/pressly/goose) - Database migrations
- [Air](https://github.com/air-verse/air) - Live reload for Go
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation

## 📚 Learn More

- [Go Documentation](https://golang.org/doc/)
- [Chi Router Documentation](https://go-chi.io/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [REST API Best Practices](https://restfulapi.net/)

---

**Built with ❤️ using Go and PostgreSQL**

If you found this project helpful, please give it a ⭐!
