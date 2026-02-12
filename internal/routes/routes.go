package routes

import (
	"fem/internal/app"

	"github.com/go-chi/chi/v5"
)

//! SetupRoutes --> configures all HTTP routes for the application
//! Request flow: Server → Router → Middleware → Handler → Response
func SetupRoutes(app *app.Application) *chi.Mux {

	//* create new chi router instance
	r := chi.NewRouter()

	//! Protected routes group --> requires valid authentication token
	//! Middleware chain: Authenticate → RequireUser → Handler
	r.Group(func (r chi.Router) {
		r.Use(app.Middleware.Authenticate) //* extracts token from Authorization header and validates it
		//* all routes in this group are protected by authentication
		r.Get("/workouts/{id}",app.Middleware.RequireUser(app.WorkoutHandler.HandleWorkoutByID)) //* GET single workout
		r.Post("/workouts",app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout)) //* CREATE new workout
		r.Put("/workouts/{id}",app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID)) //* UPDATE existing workout
		r.Delete("/workouts/{id}",app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutByID)) //* DELETE workout
	})

	//! Public routes --> no authentication required
	r.Get("/health",app.HealthCheck) //* health check endpoint

	r.Post("/users",app.UserHandler.HandleRegisterUser) //* user registration
	r.Post("/tokens/authentication",app.TokenHandler.HandleCreateToken) //* login / get auth token
	return r //* return configured router

}
