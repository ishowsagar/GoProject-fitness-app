# 🎉 Docker + Air Setup Complete!

## ✅ What's Running

Your FitTrack application is now running in Docker with these services:

```
✅ fittrack-app  - Go app with Air (Port 8080)
✅ workoutDB     - PostgreSQL database (Port 5445)
✅ workoutDB_test - Test database (Port 5500)
```

## 🚀 Quick Commands

### Start/Stop Services

```powershell
# Start everything
.\start-docker.ps1
# or
docker-compose up -d

# Stop everything
docker-compose down

# Restart app (picks up code changes)
docker-compose restart app
```

### View Logs

```powershell
# All app logs
docker-compose logs app

# Follow logs (live)
docker-compose logs -f app

# Last 50 lines
docker-compose logs --tail=50 app
```

### Check Status

```powershell
docker-compose ps
```

### Access Services

```
Backend API:     http://localhost:8080
Health Check:    http://localhost:8080/health
Database:        localhost:5445
Test Database:   localhost:5500
```

## 🔥 Air Live Reload

Air is configured and watching your Go files! However, there's a known limitation:

### Windows Docker Desktop File Watching

On Windows, Docker Desktop file change notifications can be delayed or miss changes. This is a Windows/Docker limitation, not an Air issue.

### Solutions:

#### Option 1: Manual Restart (Fast)

When you make changes:

```powershell
docker-compose restart app
```

This reloads your changes instantly!

#### Option 2: Enable Polling in .air.toml

Edit `.air.toml` and change:

```toml
[build]
  poll = true
  poll_interval = 1000  # Check every 1 second
```

Then restart:

```powershell
docker-compose restart app
```

#### Option 3: Use WSL2 (Best Performance)

If you have WSL2 enabled in Docker Desktop:

1. Move your project to WSL2 filesystem
2. Work inside WSL2
3. File watching works perfectly!

#### Option 4: Run Air Locally (Not in Docker)

```powershell
# Install Air locally
go install github.com/air-verse/air@latest

# Run
air

# Or add to PATH and use:
air -c .air.toml
```

This way Air runs natively on Windows with perfect file watching.

## 📊 Development Workflows

### Workflow 1: Docker Everything (Current)

```powershell
# Start
docker-compose up -d

# Make code changes
# Edit any .go file

# Restart to pick up changes
docker-compose restart app

# View logs
docker-compose logs -f app
```

### Workflow 2: Local Air + Docker DB

```powershell
# Start only databases
docker-compose up -d db test_db

# Run app locally with Air
air

# Changes reload automatically!
```

### Workflow 3: All Local (If you prefer)

```powershell
# Start databases
docker-compose up -d db test_db

# Run without Air
go run main.go

# Or with Air
air
```

## 🛠️ Useful Commands

### Rebuild After go.mod Changes

```powershell
docker-compose build app
docker-compose up -d app
```

### Clean Rebuild

```powershell
docker-compose build --no-cache app
docker-compose up -d app
```

### Shell into Container

```powershell
docker-compose exec app sh
```

### View Database

```powershell
docker-compose exec db psql -U postgres
```

### Stop Just the App

```powershell
docker-compose stop app
```

### Remove Version Warning

Edit `docker-compose.yml` and remove the first line:

```yaml
version: "3.8" # <-- Remove this line
```

## 🎯 Recommended Setup

For best Windows development experience, I recommend:

**Option A: Hybrid Approach (Most Practical)**

```powershell
# Keep databases in Docker
docker-compose up -d db test_db

# Run app locally with Air
air
```

This gives you:

- ✅ Perfect live reload
- ✅ Easy database management
- ✅ Fast iteration

**Option B: Full Docker (Clean & Isolated)**

```powershell
# Everything in Docker
docker-compose up -d

# After code changes
docker-compose restart app
```

This gives you:

- ✅ Identical environment everywhere
- ✅ Easy to share/deploy
- ✅ No local dependencies needed

## 🐛 Troubleshooting

### Container Won't Start

```powershell
docker-compose logs app
docker-compose down
docker-compose up -d
```

### Database Connection Error

```powershell
# Check database is running
docker-compose ps

# View database logs
docker-compose logs db

# Restart database
docker-compose restart db
```

### Port Already in Use

```powershell
# Check what's using the port
netstat -ano | findstr :8080

# Use different port in docker-compose.yml
ports:
  - "8081:8080"  # Host:Container
```

### Changes Not Reflecting

```powershell
# Option 1: Restart app
docker-compose restart app

# Option 2: Rebuild
docker-compose build app
docker-compose up -d app

# Option 3: Fresh start
docker-compose down
docker-compose up -d --build
```

## 📚 Files Created

- ✅ `docker-compose.yml` - Service orchestration
- ✅ `Dockerfile` - Multi-stage build
- ✅ `.air.toml` - Air configuration
- ✅ `.dockerignore` - Build optimization
- ✅ `Makefile` - Command shortcuts
- ✅ `start-docker.ps1` - Easy startup
- ✅ `start-docker.bat` - Batch alternative
- ✅ `DOCKER_SETUP.md` - Full documentation

## 🎊 Next Steps

1. **Test your API:**
   - Visit: http://localhost:8080/health
   - Register a user
   - Create workouts

2. **Start Frontend:**

   ```powershell
   cd frontend
   python -m http.server 3000
   ```

   - Visit: http://localhost:3000

3. **Make Changes:**
   - Edit any Go file
   - Run: `docker-compose restart app`
   - Or switch to local Air for auto-reload

4. **Choose Your Workflow:**
   - Full Docker: Easy and isolated
   - Hybrid: Best for development
   - All Local: Traditional approach

## 🔗 Resources

- Docker Compose: https://docs.docker.com/compose/
- Air: https://github.com/air-verse/air
- Make: GNU Make documentation

---

**You're all set! Happy coding! 🚀**

Choose the workflow that works best for you and start building!
