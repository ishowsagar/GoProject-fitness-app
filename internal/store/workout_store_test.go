package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ! setupTestDB --> prepares fresh test database for each test
func setupTestDB(t *testing.T) *sql.DB {
	// * connecting to test db on port 5500 (separate from dev db on 5445)
	db,err := sql.Open("pgx","host=localhost user=postgres password=postgres dbname=postgres port=5500 sslmode=disable")
	if err!= nil {
		t.Fatalf("opening test db : %v",err)
	}

	// ? - running migrations to create tables in test db
	err = Migrate(db,"../../migrations/")
	if err != nil {
		t.Fatalf("caught error while migrating db : %v",err)
	}

	// ! wiping all data so each test starts with clean slate
	_,err = db.Exec(`TRUNCATE workouts,workout_entries CASCADE`)
	if err!= nil {
		t.Fatalf("caught error while truncating db : %v",err)
	}

	return db
}

// ! TestCreateWorkout --> tests workout creation with different scenarios
func TestCreateWorkout(t *testing.T) {
	// * setting up test db and store
	db := setupTestDB(t)
	defer db.Close() // ? - cleanup when test finishes

	store := NewPostgresWorkoutStore(db)

	// ? - TABLE-DRIVEN TESTS: array of test cases to run
	test :=[]struct {
		name string    // * test name shown in output
		workout *Workout  // * input workout data
		wantErr bool   // ! expect error or not		
	}{
	 {  // ! TEST CASE 1: valid workout with all fields
		name: "valid workout",
		workout: &Workout{
			UserID: 1,
			Title: "push day",
			Description: "upper body training day",
			DurationMinutes: 60,
			CaloriesBurned: 200,		
			Entries: []WorkoutEntry{
				{
					ExerciseName: "Bench press",
					Sets: 4,
					Reps: intPointer(10),  // * pointer because reps is optional
					Weight: floatPointer(145.35),
					Notes: "warm up properly",
					OrderIndex: 1,
				},
			},
		},
		wantErr: false,  // ? - should succeed
	 },
	 {  // ! TEST CASE 2: workout with multiple entries
		name: "workout with invalid entries",
		workout: &Workout{
			UserID: 1,
			Title: "full body workout day",
			Description: "complete body training day",
			DurationMinutes: 90,
			CaloriesBurned: 500,
			Entries: []WorkoutEntry{
				{  // * plank - reps-based exercise
					ExerciseName: "plank",
					Sets: 2,
					Reps: intPointer(7),
					Notes: "keep yourself steadly warm",
					OrderIndex: 1,
				},
				{  // * squats - has both reps AND duration
					ExerciseName: "squats",
					Sets: 4,
					Reps: intPointer(12),
					DurationSeconds: intPointer(60),
					Weight: floatPointer(135.24),
					Notes: "intense training workout session",
					OrderIndex: 2,
				},
			},		
		},
		wantErr: false,  // * name says "invalid" but should still work - tests flexible entry types
	 },
	}

	// ! looping through each test case
	for _, tt := range test {
		// * t.Run creates subtest with name --> shows clearly which test failed
		t.Run(tt.name,func(t *testing.T) {
			// ? - attempting to create workout in test db
			createWorkout,err := store.CreateWorkout(tt.workout)
			
			// * if we expect an error
			if tt.wantErr {
				assert.Error(t,err)  // ? - make sure error happened
				return  // ! stop here, don't check other stuff
			}
			
			// * if we DON'T expect error, validate everything worked
			require.NoError(t,err)  // ! fail test immediately if error
			assert.Equal(t,tt.workout.Title, createWorkout.Title)
			assert.Equal(t,tt.workout.Description, createWorkout.Description)
			assert.Equal(t,tt.workout.DurationMinutes, createWorkout.DurationMinutes)
			
			// ? - now fetch the workout from db to verify it was really saved
			retrieved,err := store.GetWorkoutByID(int64(createWorkout.ID))
			require.NoError(t,err)

			// * checking main workout fields match
			assert.Equal(t,createWorkout.ID,retrieved.ID)
			assert.Equal(t,createWorkout.Entries,retrieved.Entries)
		
			// ? - checking each exercise entry was saved correctly
			for i:= range retrieved.Entries {
				assert.Equal(t,tt.workout.Entries[i].ExerciseName,retrieved.Entries[i].ExerciseName)
				assert.Equal(t,tt.workout.Entries[i].Sets,retrieved.Entries[i].Sets)
				assert.Equal(t,tt.workout.Entries[i].Reps,retrieved.Entries[i].Reps)
			}
		})
	}
}

// ! HELPER FUNCTIONS for converting values to pointers

// * intPointer --> some fields like reps/duration are optional pointers
func intPointer( i int) *int {
	return &i // ? --> returns address of the passed "i"
}

// * floatPointer --> weight is optional so needs pointer
func floatPointer( i float64) *float64 {
	return &i // ? --> returns address of the passed "i" so db can store NULL if not provided
}