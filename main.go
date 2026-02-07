package main

// imports
import (
	"fem/internal/app"
	"fmt"
	"net/http"
	"time"
	"flag"
)

// Main function where go application spins up
func main() {

	// fallback port if not specified
	var port int
	flag.IntVar(&port,"port",8080,"GO BACKEND SERVER!")
	flag.Parse() 

	app,err := app.NewApplication() //! returns Logger's output

	//  if caught any error intiting app
	if err !=nil {
		panic(err) // exits the app 
	}

	// ? - otherwise successfully imported function and executed
	fmt.Println("app is running!")

	//! server management
	
	// creating instance of a server
	server := &http.Server{
		Addr: fmt.Sprintf(":%d",port),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("App is running on port : %d\n",port)

	// handles request on this path
	http.HandleFunc("/app",HealthCheck)

	// * server listens for any incoming request
	err = server.ListenAndServe() // returns error if failed to listen for a sever
   // if caught error listening for a server
	if err !=nil {
		app.Logger.Fatal(err)
	} 

}

// cb functions  --> r pointer to keep the data persisted 
func HealthCheck(w http.ResponseWriter,req *http.Request) {
	fmt.Fprintf(w,"Status is available\n")
	 w.Write([]byte("hello from Go developer!"))
}