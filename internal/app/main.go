package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.Use(mux.CORSMethodMiddleware(router))

	// //Initializing app routes
	// router = routes.StartRouter(router)

	//Starting API and Service Servers
	http.ListenAndServe(":8080", router)
}
