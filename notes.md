# Workout Tracker API - Full Application Workflow

## 🚀 Application Startup Flow (main.go)

1. **Parse Flags** → Read port number from command line (default: 8080)
2. **Initialize Application** → Call `app.NewApplication()` to bootstrap everything
3. **Setup Server** → Create HTTP server with timeouts and configuration
4. **Start Listening** → Server begins accepting requests on specified port

---

## 🔧 Application Initialization (internal/app/app.go)

### Order of Operations:

```
1. Database Connection (store.Open)
   └─> Connects to PostgreSQL on port 5445
   └─> Returns connection pool (*sql.DB)

2. Run Migrations (store.Migratefs)
   └─> Applies SQL migrations using goose
   └─> Creates/updates tables: users, workouts, workout_entries, tokens

3. Initialize Logger
   └─> Creates centralized logger for error tracking

4. Create Store Layer (Database Operations)
   ├─> WorkoutStore → CRUD for workouts and entries
   ├─> UserStore → User registration and authentication
   └─> TokenStore → Token generation and validation

5. Create Handler Layer (HTTP Request Handlers)
   ├─> WorkoutHandler → Handles /workouts endpoints
   ├─> UserHandler → Handles /users registration
   └─> TokenHandler → Handles /tokens/authentication (login)

6. Initialize Middleware
   └─> UserMiddleware → Authenticates requests via Bearer token

7. Wire Everything Together
   └─> Return Application struct with all dependencies
```

---

## 🛣️ Route Setup (internal/routes/routes.go)

### Public Routes (No Authentication):

- `GET /health` → Health check endpoint
- `POST /users` → User registration
- `POST /tokens/authentication` → Login / Get auth token

### Protected Routes (Requires Authentication):

**Middleware Chain: Authenticate → RequireUser → Handler**

- `GET /workouts/{id}` → Fetch single workout
- `POST /workouts` → Create new workout
- `PUT /workouts/{id}` → Update workout (owner only)
- `DELETE /workouts/{id}` → Delete workout (owner only)

---

## 🔐 Authentication Flow

### 1️⃣ User Registration (`POST /users`)

```
Client Request (JSON)
   ↓
UserHandler.HandleRegisterUser
   ↓
Validate input (username, email, password format)
   ↓
Hash password using bcrypt (cost: 12)
   ↓
Save to database → users table
   ↓
Return user object (password hash excluded)
```

### 2️⃣ Login (`POST /tokens/authentication`)

```
Client Request (username + password)
   ↓
TokenHandler.HandleCreateToken
   ↓
Lookup user by username
   ↓
Compare plaintext password with bcrypt hash
   ↓
Generate secure random token (32 bytes → base32 encoded)
   ↓
Hash token with SHA-256 for database storage
   ↓
Save token hash to database (expires in 24 hours)
   ↓
Return plaintext token to client (only time they see it!)
```

### 3️⃣ Using Token (Protected Routes)

```
Client sends: Authorization: Bearer <TOKEN>
   ↓
Authenticate Middleware
   ├─> Extract token from header
   ├─> Hash incoming token with SHA-256
   ├─> Lookup in database (check expiry)
   ├─> Fetch associated user
   └─> Inject user into request context
   ↓
RequireUser Middleware
   ├─> Extract user from context
   ├─> Check if anonymous
   └─> Block if not authenticated
   ↓
Handler executes with authenticated user
```

---

## 📝 Complete CRUD Workflow Example (Workouts)

### CREATE Workout (`POST /workouts`)

```
1. Authentication Layer
   ├─> Validate Bearer token
   └─> Get authenticated user from context

2. Request Processing
   ├─> Decode JSON body into Workout struct
   ├─> Attach user.ID to workout.UserID (ownership)
   └─> Validate workout data

3. Database Transaction
   ├─> Begin transaction
   ├─> Insert into workouts table → get workout ID
   ├─> Loop through entries array
   │   └─> Insert each exercise into workout_entries table
   ├─> Commit transaction (all or nothing)
   └─> Rollback if any error occurs

4. Response
   └─> Return created workout with all entries (HTTP 201)
```

### READ Workout (`GET /workouts/{id}`)

```
1. Extract ID from URL path parameter
   └─> utils.ReadIDParam(r) converts string to int64

2. Fetch from Database
   ├─> Query workouts table by ID
   ├─> Query workout_entries table (ordered by order_index)
   └─> Combine into single Workout object

3. Return Response
   └─> Send workout JSON (HTTP 200)
```

### UPDATE Workout (`PUT /workouts/{id}`)

```
1. Authentication & Authorization
   ├─> Validate Bearer token
   ├─> Get authenticated user
   ├─> Fetch workout owner from database
   └─> Verify current user == workout owner (HTTP 403 if not)

2. Partial Update Logic
   ├─> Fetch existing workout
   ├─> Only update fields that are present in request
   │   (uses pointers: nil = no change, value = update)
   └─> Replace all entries with new ones

3. Database Transaction
   ├─> Begin transaction
   ├─> Update workouts table
   ├─> Delete all old entries (CASCADE)
   ├─> Insert new entries
   └─> Commit transaction

4. Response
   └─> Return updated workout (HTTP 200)
```

### DELETE Workout (`DELETE /workouts/{id}`)

```
1. Authentication & Authorization
   ├─> Validate Bearer token
   ├─> Get authenticated user
   ├─> Fetch workout owner
   └─> Verify ownership

2. Delete from Database
   ├─> DELETE FROM workouts WHERE id = ?
   └─> CASCADE automatically deletes entries

3. Response
   └─> HTTP 204 No Content (success, no body)
```

---

## 🏗️ Architecture Pattern (Layered Architecture)

```
┌─────────────────────────────────────────────┐
│           CLIENT (Postman/Browser)          │
└─────────────────┬───────────────────────────┘
                  │ HTTP Request
                  ↓
┌─────────────────────────────────────────────┐
│              ROUTES LAYER                   │
│  • Route matching (/workouts, /users)       │
│  • Middleware attachment                    │
└─────────────────┬───────────────────────────┘
                  ↓
┌─────────────────────────────────────────────┐
│           MIDDLEWARE LAYER                  │
│  • Authenticate (token validation)          │
│  • RequireUser (authorization)              │
│  • Set user in context                      │
└─────────────────┬───────────────────────────┘
                  ↓
┌─────────────────────────────────────────────┐
│            HANDLER LAYER (API)              │
│  • Parse request (JSON decode)              │
│  • Validate input                           │
│  • Business logic                           │
│  • Call store methods                       │
│  • Format response (JSON encode)            │
└─────────────────┬───────────────────────────┘
                  ↓
┌─────────────────────────────────────────────┐
│             STORE LAYER                     │
│  • Interface contracts (loose coupling)     │
│  • SQL queries                              │
│  • Transaction management                   │
│  • Data mapping (DB ↔ Structs)             │
└─────────────────┬───────────────────────────┘
                  ↓
┌─────────────────────────────────────────────┐
│          DATABASE (PostgreSQL)              │
│  Tables: users, workouts, entries, tokens   │
└─────────────────────────────────────────────┘
```

---

## 🔑 Key Design Decisions

### 1. **Interface-Based Store Layer**

- Allows swapping PostgreSQL for MySQL/MongoDB without changing handlers
- Makes testing easier (mock implementations)
- Loose coupling between layers

### 2. **Transaction Management**

- Ensures data consistency (workout + entries saved together)
- Rollback on any error (atomicity)

### 3. **Token Security**

- Plaintext token only sent once (at creation)
- SHA-256 hash stored in database
- Can't reverse engineer from database

### 4. **Password Security**

- Bcrypt hashing (cost factor 12)
- Salt automatically included
- Slow algorithm (prevents brute force)

### 5. **Ownership Verification**

- Users can only modify their own workouts
- `user_id` foreign key links workouts to owners
- Authorization checks before update/delete

### 6. **Partial Updates**

- Pointer fields (*string, *int) allow nil values
- Only update fields present in request
- Flexible API for clients

### 7. **Context for User Passing**

- Middleware injects user into request context
- Available to all downstream handlers
- Type-safe with custom context key

---

## 📊 Database Schema

```sql
users
  ├─ id (primary key)
  ├─ username (unique)
  ├─ email
  ├─ password_hash
  ├─ bio
  └─ timestamps

workouts
  ├─ id (primary key)
  ├─ user_id (foreign key → users.id)
  ├─ title
  ├─ description
  ├─ duration_minutes
  └─ calories_burned

workout_entries
  ├─ id (primary key)
  ├─ workout_id (foreign key → workouts.id, CASCADE)
  ├─ exercise_name
  ├─ sets, reps, duration_seconds, weight
  ├─ notes
  └─ order_index

tokens
  ├─ hash (primary key, SHA-256)
  ├─ user_id (foreign key → users.id, CASCADE)
  ├─ expiry (timestamp)
  └─ scope (authentication, password-reset, etc.)
```

---

## 🎯 Request → Response Journey

**Example: Creating a workout**

```
1. Client sends POST request to http://localhost:8080/workouts
   Headers: Authorization: Bearer ABC123...
   Body: { title, description, entries: [...] }

2. Server receives request → chi router matches /workouts

3. Authenticate middleware runs
   → Extracts "ABC123..." from header
   → Hashes it with SHA-256
   → Queries tokens table
   → Finds user_id = 5, expiry still valid
   → Queries users table for full user object
   → Injects user into request context

4. RequireUser middleware runs
   → Checks if user.IsAnonymousUser()
   → User is authenticated ✓
   → Allows request to continue

5. WorkoutHandler.HandleCreateWorkout executes
   → Gets user from context (user_id = 5)
   → Decodes JSON body into Workout struct
   → Sets workout.UserID = 5 (ownership)
   → Calls workoutStore.CreateWorkout(workout)

6. PostgresWorkoutStore.CreateWorkout runs
   → Begins database transaction
   → INSERT INTO workouts (...) VALUES (...) RETURNING id
   → Loops through each entry
   → INSERT INTO workout_entries (...) VALUES (...)
   → Commits transaction
   → Returns complete workout object

7. Handler formats response
   → Wraps in Envelope{"workout": workout}
   → JSON encodes with indentation
   → Sets Content-Type: application/json
   → Writes HTTP 201 Created

8. Client receives response
   {
     "workout": {
       "id": 10,
       "user_id": 5,
       "title": "Morning Cardio",
       ...
       "entries": [...]
     }
   }
```

---

## 🛠️ Development Notes

### Port Configuration

- **Database**: Port 5445 (not 5432)
  - Windows reserves 5345-5444 range
  - Changed to avoid "bind: address already in use" errors
  - Updated in: database.go connection string
  - Updated in: docker-compose.yml port mappings
  - Test DB on port 5500

### Air Hot Reload

- Watches Go files for changes
- Automatically rebuilds and restarts server
- Configuration in .air.toml (if exists)

### Database Migrations

- Managed by goose
- Located in migrations/ folder
- Run automatically on app startup
- Manual command: `goose -dir migrations postgres "connection_string" up`

---

## 🐛 Common Issues Fixed

1. **Invalid Token Error**
   - Bug: `GetUserToken` returned `nil,err` instead of `user,nil`
   - Bug: Query used `token` table instead of `tokens` (typo)
   - Fix: Return actual user object and fix table name

2. **Capitalization Issues**
   - Go exports: Capital = public, lowercase = private
   - Bug: `isAnonymousUser()` was private
   - Fix: Changed to `IsAnonymousUser()` for cross-package access

3. **Missing Return Statements**
   - Bug: Middleware continued after writing error response
   - Fix: Added `return` after error responses

4. **Port Conflicts**
   - Windows reserves certain port ranges
   - Solution: Use ports outside reserved ranges (5445, 5500)

---

## ✅ Testing the API

```bash
# 1. Register User
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "email": "john@example.com", "password": "pass123"}'

# 2. Login (Get Token)
curl -X POST http://localhost:8080/tokens/authentication \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "pass123"}'

# Response: {"auth_token": {"token": "ABC123XYZ...", "expiry": "..."}}

# 3. Create Workout (Use Token)
curl -X POST http://localhost:8080/workouts \
  -H "Authorization: Bearer ABC123XYZ..." \
  -H "Content-Type: application/json" \
  -d '{"title": "Leg Day", "duration_minutes": 60, "entries": [...]}'

# 4. Get Workout
curl -X GET http://localhost:8080/workouts/1 \
  -H "Authorization: Bearer ABC123XYZ..."

# 5. Update Workout
curl -X PUT http://localhost:8080/workouts/1 \
  -H "Authorization: Bearer ABC123XYZ..." \
  -d '{"title": "Updated Leg Day"}'

# 6. Delete Workout
curl -X DELETE http://localhost:8080/workouts/1 \
  -H "Authorization: Bearer ABC123XYZ..."
```
