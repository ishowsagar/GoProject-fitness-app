package app

import (
	"database/sql"
	"fem/internal/api"
	"fem/internal/middleware"
	"fem/internal/store"
	"fem/migrations"
	"fmt"
	"log"
	"net/http"
	"os"
)

//! types declarement
//! Application struct --> holds all dependencies needed across the app
type Application struct {
	Logger *log.Logger //* centralized logger for error tracking
	WorkoutHandler *api.WorkoutHandler //* handles workout CRUD operations
	UserHandler *api.UserHandler //* handles user registration
	TokenHandler *api.TokenHandler //* handles authentication token creation
	Middleware middleware.UserMiddleware //* authentication middleware for protected routes
	DB *sql.DB //* database connection pool
}

//! NewApplication --> constructor that initializes entire app with all dependencies
func NewApplication() (*Application,error) {

	//* establishing database connection
	pgDb,err := store.Open()
	if err != nil {
		return nil,err
	}

	//* running database migrations --> ensures tables are up to date
	err = store.Migratefs(pgDb,migrations.FS,".")
	if err != nil {
		panic(err)
	}


	//* creating logger instance with date and time stamps
	logger := log.New(os.Stdout,"",log.Ldate | log.Ltime) 

	//! Initializing all store instances --> database layer that talks to postgres
	workoutStore := store.NewPostgresWorkoutStore(pgDb) //* workout operations
	userStore := store.NewPostUserStore(pgDb) //* user operations
	tokenStore := store.NewPostgresTokenStore(pgDb) //* token operations

	//! Initializing all handler instances --> HTTP request handlers
	workoutHandler := api.NewWorkoutHandler(workoutStore,logger) //* workout endpoints
	userHandler := api.NewUserHandler(userStore,logger) //* user registration endpoint
	tokenHandler := api.NewTokenHandler(tokenStore,userStore,logger) //* authentication endpoint
	mwHandler := middleware.UserMiddleware{UserStore: userStore} //* middleware for auth checks

	//* creating Application instance with all dependencies wired up
	app := &Application{
		Logger : logger,
		WorkoutHandler: workoutHandler,
		UserHandler: userHandler,
		TokenHandler: tokenHandler,
		Middleware : mwHandler,
		DB: pgDb,
	}
	
	return app,nil //* return initialized app ready to handle requests

}

//! HealthCheck --> simple endpoint to verify server is running
//! GET /health --> returns status message
func (a *Application) HealthCheck(w http.ResponseWriter,req *http.Request) {
	fmt.Fprintf(w," 🦖FitTrack API is healthy 🪐 and running with Docker + Air! 🔥\n") //* simple text response
}