# Go Method Receivers Explained

## Quick Summary

- **Receivers provide direct access to struct fields without passing them as parameters every time**
- **Pointer receivers directly influence the actual struct, not a copy**
- **Methods with receivers can access struct data internally without explicit parameter passing**
- **All methods on the same receiver are grouped together and share access to the same struct instance**
- **Pointer ensures changes affect the real struct, value receiver only affects a copy**
- **Receivers eliminate repetitive parameter passing for things like database connections**
- **The receiver acts like a built-in first parameter that's always available inside the method**

## How to Call Methods & Interface Connection

- **You DON'T call methods using the receiver name (pg) - that's only inside the function**
- **You call methods on the INSTANCE of the struct that you created**
- **Example: `store.GetWorkoutByID(5)` where store is your instance, not pg**
- **Inside the method, Go automatically passes the instance as `pg` (the receiver)**
- **Interfaces list method SIGNATURES (names, parameters, return types)**
- **If a struct has methods matching ALL signatures in an interface, it automatically implements that interface**
- **No explicit declaration needed - if methods match, the struct satisfies the interface**
- **This allows you to use the struct wherever the interface type is expected**
- **Different structs can implement the same interface with different internal logic**
- **Interface defines WHAT methods must exist, receiver implementation defines HOW they work**

## How Interfaces Actually Work

- **Interface is a CONTRACT/CHECKLIST - it doesn't trigger or call methods itself**
- **Interface says "whoever uses me MUST have these specific methods"**
- **PostgresWorkoutStore is the TYPE (like a blueprint), not a parent**
- **NewPostgresWorkoutStore is a FACTORY FUNCTION that creates an instance of that type**
- **All methods with the same receiver are automatically grouped and belong to that type**
- **When you call a method on an instance, Go looks up which method matches that type and calls it**
- **Interface enables POLYMORPHISM - different types can be used interchangeably if they have the same methods**
- **The interface doesn't control the methods - it just defines requirements**
- **Your code decides when to call methods - interface just ensures the methods exist**
- **Think: Interface = job description, Struct with methods = qualified candidate who meets requirements**

# Structure clarification:

PostgresWorkoutStore is the TYPE (blueprint)
NewPostgresWorkoutStore creates an instance of that type
Methods with (pg \*PostgresWorkoutStore) receiver belong to that type
Not parent/child - more like: Type defines methods, instance uses them
How it works:

You create instance: store := NewPostgresWorkoutStore(db)
You call method: store.GetWorkoutByID(5)
Go finds the matching method for that type and executes it
Inside that method, pg refers to your store instance

YES! You got it! 🎯 Let me confirm your understanding:

✅ Define a struct → type PostgresWorkoutStore struct { db *sql.DB }

✅ 3 methods point to/use that struct → They all have (pg *PostgresWorkoutStore) receiver

✅ Function creates and returns instance → NewPostgresWorkoutStore(pgDb) returns &PostgresWorkoutStore{db: pgDb}

✅ Store that instance → workoutStore := store.NewPostgresWorkoutStore(pgDb)

✅ All methods are accessible through this instance → workoutStore.CreateWorkout(), workoutStore.GetWorkoutByID(), workoutStore.UpdateWorkout()

Small clarification on #5:

Methods aren't "stored under" the instance
Methods are attached to the TYPE (PostgresWorkoutStore)
The instance has access to all methods defined for its type
Every instance of PostgresWorkoutStore can use the same 3 methods
In your code:

Inside each method, pg refers to whatever instance you called it on (in this case, workoutStore).

You nailed the concept! 🚀
## What is a Method Receiver?

The syntax `func (pg *PostgresWorkoutStore) GetWorkoutByID(...)` makes this a **method** instead of a regular function.

The `(pg *PostgresWorkoutStore)` part is called a **receiver** - it attaches the function to a type.

## Regular Function vs Method

```go
// ❌ Regular function - NOT attached to a type
func GetWorkoutByID(id int64) (*Workout, error) {
    // How would this access the database? No way to get db connection!
    // Would need to pass db as parameter every time
}

// ✅ Method - attached to PostgresWorkoutStore type
func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
    // Can access pg.db because pg is the receiver
    pg.db.QueryRow(query, id) // ← Works! Has access to database
}
```

## Why Use a Receiver?

### 1. Access to Struct Fields

```go
type PostgresWorkoutStore struct {
    db *sql.DB  // ← Need access to this database connection
}

func (pg *PostgresWorkoutStore) GetWorkoutByID(id int64) (*Workout, error) {
    // pg.db gives you access to the database connection
    // 'pg' is like 'this' or 'self' in other languages
    result := pg.db.QueryRow(query, id)
    return result, nil
}
```

### 2. Clean Method Call Syntax

```go
// Create store instance
store := NewPostgresWorkoutStore(db)

// Call method ON the instance (clean syntax)
workout, err := store.GetWorkoutByID(5)

// vs regular function would be (awkward):
workout, err := GetWorkoutByID(store, 5)
```

### 3. Organize Related Functions Together

All methods with the same receiver are grouped together conceptually:

```go
func (pg *PostgresWorkoutStore) CreateWorkout(...) { }
func (pg *PostgresWorkoutStore) GetWorkoutByID(...) { }
func (pg *PostgresWorkoutStore) UpdateWorkout(...) { }
// All these belong to PostgresWorkoutStore
```

## The `*` (Pointer Receiver)

```go
func (pg *PostgresWorkoutStore) GetWorkoutByID(...)
//      ^ This asterisk means "pointer receiver"
```

### Why Use Pointer Receiver?

**1. Can modify the struct if needed**

```go
func (pg *PostgresWorkoutStore) UpdateConnection(newDB *sql.DB) {
    pg.db = newDB  // ← Can modify because it's a pointer
}
```

**2. Avoids copying the entire struct (more efficient)**

```go
// Value receiver - copies entire struct every call (slow if struct is large)
func (pg PostgresWorkoutStore) GetWorkoutByID(...) { }

// Pointer receiver - only copies the pointer (fast, efficient)
func (pg *PostgresWorkoutStore) GetWorkoutByID(...) { }
```

**3. Consistency**

- If ANY method on a type uses pointer receiver, ALL should use pointer receiver
- Maintains consistency across your codebase

## How Receivers Enable Interfaces

Your interface definition:

```go
type WorkoutStore interface {
    GetWorkoutByID(id int64) (*Workout, error)
    CreateWorkout(*Workout) (*Workout, error)
    UpdateWorkout(*Workout) error
}
```

By defining methods with `(pg *PostgresWorkoutStore)` receiver, you're implementing the interface:

```go
// PostgresWorkoutStore automatically implements WorkoutStore
// because it has all the required methods

var store WorkoutStore = NewPostgresWorkoutStore(db)
store.GetWorkoutByID(5)  // ← Works!

// This allows polymorphism - could swap with different implementation:
// var store WorkoutStore = NewMongoWorkoutStore(mongoClient)
// Same interface, different database!
```

## Real-World Analogy

Think of a receiver like a remote control:

```go
type TV struct {
    channel int
    volume  int
}

// Method with receiver - remote control attached to THIS TV
func (tv *TV) ChangeChannel(newChannel int) {
    tv.channel = newChannel  // Changes THIS TV's channel
}

// Usage:
myTV := &TV{channel: 5, volume: 10}
myTV.ChangeChannel(7)  // Remote control for myTV

yourTV := &TV{channel: 3, volume: 15}
yourTV.ChangeChannel(7)  // Remote control for yourTV
```

Each TV instance has its own remote (receiver) that controls its own state.

## Summary

| Concept         | Purpose                        | Example                                     |
| --------------- | ------------------------------ | ------------------------------------------- |
| **Receiver**    | Attaches function to a type    | `func (pg *PostgresWorkoutStore)`           |
| **Pointer `*`** | Efficient, can modify struct   | `*PostgresWorkoutStore`                     |
| **Method Call** | Clean syntax with dot notation | `store.GetWorkoutByID(5)`                   |
| **Interface**   | Defines behavior contract      | Multiple types can implement same interface |

**Key Point:** Receivers give functions access to struct fields and enable object-oriented programming patterns in Go!
