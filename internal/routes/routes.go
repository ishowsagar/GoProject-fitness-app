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

	// attaching mw for token auth with chi r.Group method
	r.Group(func (r chi.Router) {
		r.Use(app.Middleware.Authenticate)
		// ! whatever type struct they was connected to --> accessed via same struct
		// after executing mw fnc --> we pass group of routes where we want them to execute before these
		r.Get("/workouts/{id}",app.Middleware.RequireUser(app.WorkoutHandler.HandleWorkoutByID)) //* checks user if anon or even has context key before before executing the requested func call
		r.Post("/workouts",app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}",app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/workouts/{id}",app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutByID))
	})

	r.Get("/health",app.HealthCheck)
	// ! whatever type struct they was connected to --> accessed via same struct

	r.Post("/users",app.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authentication",app.TokenHandler.HandleCreateToken)
	return r

}
