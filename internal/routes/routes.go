package routes

import (
	"fem/internal/app"

	"github.com/go-chi/chi/v5"
)

// ! Handler in server configuration --> directs all req to SetupRoutes

//  Req comes to server --> passes to "r" handler --> passes to setupRoutes --> goes to "r" handler
//  --> goes to desired route handler fnc --> respond to the request with desired response
func SetupRoutes(app *app.Application) *chi.Mux {

	// * handler is defined and which process the client's req
	r := chi.NewRouter()

	// routes
	r.Get("/health",app.HealthCheck)
	r.Get("/workouts/{id}",app.WorkoutHandler.HandleWorkoutByID)
	r.Post("/workouts",app.WorkoutHandler.HandleCreateWorkout)
	return r

}
