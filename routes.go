package main

// Route is the model for the router setup
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc HandlerFunc
}

// Routes are the main setup for our Router
type Routes []Route

var routes = Routes{
	// meta services
	Route{"Healthcheck", "GET", "/healthcheck", HealthcheckHandler},

	//=== Add Photos ===
	Route{"AddPhotos", "POST", "/photos", AddPhotosHandler},

	//=== Front End ===
	// is added in server.go to avoid bad interaction with gorilla mux
}
