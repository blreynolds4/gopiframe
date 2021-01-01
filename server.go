package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// StartServer Wraps the mux Router and uses the Negroni Middleware
func StartServer(ctx AppContext) {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = makeHandler(ctx, route.HandlerFunc)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// put in the routing of the UI
	router.PathPrefix("/ui/").Handler(
		http.StripPrefix("/ui/", http.FileServer(http.Dir(ctx.UIPath))))

	// security
	// var isDevelopment = false
	// if ctx.Env == local {
	// 	isDevelopment = true
	// }
	// secureMiddleware := secure.New(secure.Options{
	// 	IsDevelopment:      isDevelopment,                                // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
	// 	AllowedHosts:       []string{"localhost:8080", "localhost:8081"}, // AllowedHosts is a list of fully qualified domain names that are allowed (CORS)
	// 	ContentTypeNosniff: false,                                        // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
	// 	BrowserXssFilter:   false,                                        // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
	// })

	// start now
	n := negroni.New()
	n.Use(negroni.NewLogger())
	// n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.UseHandler(router)
	log.Println("===> Starting app (v" + ctx.Version + ") on port " + ctx.Port + " in " + ctx.Env + " mode.")
	if ctx.Env == local {
		n.Run("localhost:" + ctx.Port)
	} else {
		n.Run(":" + ctx.Port)
	}
}
