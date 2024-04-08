package routes

import (
	"fandm/internal/register"
	"fandm/internal/routes/actors"
	"fandm/internal/routes/conditions"
	"fandm/internal/routes/treatments"
	"fmt"
	"net/http"
)

func Routes() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", register.Register)
	mux.HandleFunc("GET /actors/", actors.GetActors)
	mux.HandleFunc("GET /conditions/", conditions.GetConditions)
	mux.HandleFunc("GET /treatments/", treatments.GetTreatments)

	serveStatic()

	fmt.Println("Server started on localhost:8080")
	http.ListenAndServe("127.0.0.1:8080", mux)
}

func serveStatic() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
}
