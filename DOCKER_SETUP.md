# 🐳 Docker + Air Setup Guide

Complete Docker Compose setup with Air live reload for FitTrack application.

## 📦 What's Included

### Files Created

- **docker-compose.yml** - Multi-service Docker setup
- **Dockerfile** - Multi-stage build (development + production)
- **.air.toml** - Air configuration for live reload
- **.dockerignore** - Docker build optimization
- **Makefile** - Convenient command shortcuts
- **start-docker.ps1** - Windows PowerShell startup script
- **start-docker.bat** - Windows batch startup script

### Services Configured

1. **app** - Go application with Air live reload
2. **db** - PostgreSQL database (port 5445 → 5432)
3. **test_db** - Test database (port 5500 → 5432)

## 🚀 Quick Start

### Option 1: PowerShell Script (Easiest)

```powershell
.\start-docker.ps1
```

### Option 2: Batch Script

```cmd
start-docker.bat
```

### Option 3: Docker Compose Commands

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop all services
docker-compose down
```

### Option 4: Using Makefile

```bash
# Start everything
make start

# View logs
make logs

# Restart app
make restart

# Stop everything
make stop
```

## 📋 Available Commands

### Makefile Commands

```bash
make help          # Show all available commands
make build         # Build Go binary
make run           # Run locally (no Docker)
make dev           # Run with Air locally
make test          # Run tests
make clean         # Clean build artifacts

# Docker commands
make docker-build  # Build Docker images
make docker-up     # Start all services
make docker-down   # Stop all services
make docker-logs   # View application logs
make docker-restart    # Restart app service
make docker-rebuild    # Rebuild everything

# Database commands
make db-up         # Start only database
make db-down       # Stop database
make db-logs       # View database logs

# Combined
make start         # Start everything
make stop          # Stop everything
make restart       # Restart app
make logs          # View logs
```

### Docker Compose Commands

```bash
# Start services
docker-compose up              # Start in foreground
docker-compose up -d           # Start in background (detached)
docker-compose up --build      # Rebuild and start

# Stop services
docker-compose stop            # Stop services
docker-compose down            # Stop and remove containers
docker-compose down -v         # Stop and remove volumes (⚠️ deletes data)

# View logs
docker-compose logs            # All logs
docker-compose logs -f app     # Follow app logs
docker-compose logs -f db      # Follow database logs

# Restart services
docker-compose restart         # Restart all
docker-compose restart app     # Restart only app

# Rebuild
docker-compose build           # Build images
docker-compose build --no-cache    # Clean rebuild

# Check status
docker-compose ps              # List running services
docker-compose top             # Show running processes

# Execute commands
docker-compose exec app sh     # Shell into app container
docker-compose exec db psql -U postgres    # Access database
```

## 🔧 Development Workflow

### 1. Start Development Environment

```bash
# Using script (recommended)
.\start-docker.ps1

# Or using Docker Compose
docker-compose up -d

# Or using Makefile
make start
```

### 2. Edit Your Code

- Make changes to any `.go` file
- Air automatically detects changes
- Application rebuilds and restarts
- Changes reflect immediately

### 3. View Logs

```bash
# Real-time logs
docker-compose logs -f app

# Or
make logs
```

### 4. Access Services

- **Backend API**: http://localhost:8080
- **Database**: localhost:5445 (from host) or localhost:5432 (from containers)
- **Test DB**: localhost:5500

### 5. Stop Services

```bash
docker-compose down

# Or
make stop
```

## 🏗️ Architecture

### Service Communication

```
┌─────────────────────────────────────────┐
│           Host Machine                   │
│  ┌─────────────────────────────────┐   │
│  │     fittrack-network             │   │
│  │                                  │   │
│  │  ┌──────────┐    ┌───────────┐ │   │
│  │  │   app    │───→│    db     │ │   │
│  │  │  :8080   │    │  :5432    │ │   │
│  │  └──────────┘    └───────────┘ │   │
│  │                                  │   │
│  │          ┌───────────┐          │   │
│  │          │  test_db  │          │   │
│  │          │   :5432   │          │   │
│  │          └───────────┘          │   │
│  └─────────────────────────────────┘   │
│         ↓           ↓          ↓        │
│      :8080       :5445      :5500       │
└─────────────────────────────────────────┘
```

### Volume Mounts

- **App Code**: `.:/app` (live sync for hot reload)
- **Database Data**: `./database/postgres-data:/var/lib/postgresql/data`
- **Test DB Data**: `./database/postgres-test-data:/var/lib/postgresql/data`

## 🔥 Air (Live Reload)

### How It Works

1. Air watches your Go files for changes
2. When you save a file, Air:
   - Compiles the code
   - Restarts the application
   - Shows build errors if any
3. All happens automatically in seconds!

### Configuration (.air.toml)

```toml
[build]
  cmd = "go build -o ./tmp/main ."
  bin = "./tmp/main"
  exclude_dir = ["frontend", "database", "tmp"]
  include_ext = ["go", "html", "tmpl"]
```

### View Build Errors

Check `build-errors.log` file or watch logs:

```bash
docker-compose logs -f app
```

## 🗄️ Database Configuration

### Connection Details

#### From Host Machine (Your Windows PC)

```
Host:     localhost
Port:     5445
User:     postgres
Password: postgres
Database: postgres
```

#### From Docker Containers

```
Host:     db
Port:     5432
User:     postgres
Password: postgres
Database: postgres
```

### Environment Variables

The app automatically detects environment and connects appropriately:

- **Docker**: Uses `DB_HOST=db` and `DB_PORT=5432`
- **Local**: Uses `DB_HOST=localhost` and `DB_PORT=5445`

### Connect with psql

```bash
# From host
docker-compose exec db psql -U postgres

# From container
docker exec -it workoutDB psql -U postgres
```

### Backup Database

```bash
# Backup
docker-compose exec db pg_dump -U postgres postgres > backup.sql

# Restore
docker-compose exec -T db psql -U postgres postgres < backup.sql
```

## 🐛 Troubleshooting

### Port Already in Use

```bash
# Find process using port
netstat -ano | findstr :8080

# Stop specific container
docker-compose stop app

# Or use different port in docker-compose.yml
ports:
  - "8081:8080"
```

### Database Connection Failed

```bash
# Check if database is running
docker-compose ps

# View database logs
docker-compose logs db

# Restart database
docker-compose restart db

# Wait for database to be ready
docker-compose up -d db
timeout 10  # Wait 10 seconds
docker-compose up -d app
```

### Air Not Reloading

```bash
# Check if Air is running
docker-compose logs app | grep "Air"

# Restart app service
docker-compose restart app

# Rebuild with no cache
docker-compose build --no-cache app
docker-compose up -d app
```

### Build Errors

```bash
# View detailed logs
docker-compose logs -f app

# Check build-errors.log
docker-compose exec app cat build-errors.log

# Shell into container for debugging
docker-compose exec app sh
go build -v
```

### Container Won't Start

```bash
# Check container status
docker-compose ps

# View all logs
docker-compose logs

# Remove and recreate
docker-compose down
docker-compose up -d --force-recreate
```

### Clean Slate Reset

```bash
# Stop and remove everything (⚠️ deletes database data)
docker-compose down -v

# Remove build artifacts
rm -rf tmp/ database/postgres-data/ database/postgres-test-data/

# Rebuild everything
docker-compose build --no-cache
docker-compose up -d
```

## 📊 Monitoring

### View Resource Usage

```bash
# Real-time stats
docker stats

# Specific service
docker stats fittrack-app
```

### Check Disk Usage

```bash
docker system df
```

### Clean Up Unused Resources

```bash
# Remove unused containers, networks, images
docker system prune

# Remove volumes too (⚠️ deletes data)
docker system prune -a --volumes
```

## 🚀 Production Deployment

### Build Production Image

```dockerfile
# Build production stage
docker build --target production -t fittrack:prod .
```

### Run Production Container

```bash
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=your-prod-db.com \
  -e DB_PORT=5432 \
  -e DB_USER=produser \
  -e DB_PASSWORD=prodpass \
  -e DB_NAME=proddb \
  fittrack:prod
```

### Docker Compose for Production

Create `docker-compose.prod.yml`:

```yaml
version: "3.8"
services:
  app:
    build:
      context: .
      target: production
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=postgres
    restart: always
```

Run with:

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## 📚 Learn More

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Air GitHub Repository](https://github.com/cosmtrek/air)
- [Multi-stage Docker Builds](https://docs.docker.com/build/building/multi-stage/)
- [Docker Networking](https://docs.docker.com/network/)

## ✅ Checklist

- [ ] Docker Desktop installed and running
- [ ] Python installed (for frontend server)
- [ ] `docker-compose.yml` configured
- [ ] `.air.toml` created
- [ ] Database connection updated with env vars
- [ ] Run `.\start-docker.ps1`
- [ ] Access http://localhost:8080/health
- [ ] Edit a Go file and watch it reload
- [ ] Frontend running on http://localhost:3000

---

**Now you have a professional development environment with hot reload! 🔥**

Edit your code → Air detects changes → Auto rebuild → Instant feedback 🚀
