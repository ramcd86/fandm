package routes

import (
	"fandm/internal/register"
	relationshiphandler "fandm/internal/routes/relationship-handlers"
	resourcehandlers "fandm/internal/routes/route-handlers"
	"fmt"
	"net/http"
)

func Routes() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", register.Register)
	mux.HandleFunc("GET /actors/", resourcehandlers.GetActors)
	mux.HandleFunc("GET /conditions/", resourcehandlers.GetConditions)
	mux.HandleFunc("GET /treatments/", resourcehandlers.GetTreatments)
	mux.HandleFunc("POST /relationship", relationshiphandler.CreateNewRelationship)

	serveStatic()

	fmt.Println("Server started on localhost:8080")
	http.ListenAndServe("127.0.0.1:8080", mux)
}

func serveStatic() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
}
