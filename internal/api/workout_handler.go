package api

import (
	"encoding/json"
	"fem/internal/store"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// types declaration
type WorkoutHandler struct {
	workstore store.WorkoutStore

}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
return &WorkoutHandler{
	workstore: workoutStore,
}
}

//! methods --> have base method WorkoutHandler ( points to type which persists changes across app) --> other called via base this one	
func (wh *WorkoutHandler) HandleWorkoutByID(w http.ResponseWriter, req *http.Request) {
// extracting id from url via chi
paramsWorkoutID := chi.URLParam(req,"id") // passing req and "slug" being route params

// if id not found on url params
if paramsWorkoutID == "" {
	http.NotFound(w,req)
	return
}

workoutID,err := strconv.ParseInt(paramsWorkoutID,10,64)
if err != nil {
http.NotFound(w,req)
return
}
fmt.Fprintf(w,"this is the workout id : %d\n",workoutID)
}

func (wh *WorkoutHandler) HandleCreateWorkout (w http.ResponseWriter, req *http.Request) {
var workout  store.Workout // * follows type def of this struct
err := json.NewDecoder(req.Body).Decode(&workout)

if err !=nil {
	fmt.Println(err)
	http.Error(w,"failed to create workout",http.StatusInternalServerError)
	return
}

createWorkout,err := wh.workstore.CreateWorkout(&workout)
if err !=nil {
	fmt.Println(err)
	http.Error(w,"failed to create workout",http.StatusInternalServerError)
	return
}

w.Header().Set("Content-type","application/json")
json.NewEncoder(w).Encode(createWorkout)

}