package routes

import (
	"fandm/internal/register"
	relationshiphandler "fandm/internal/routes/relationship-handlers"
	resourcehandlers "fandm/internal/routes/route-handlers"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

func Routes() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", register.Register)
	mux.HandleFunc("GET /actors/", resourcehandlers.GetActors)
	mux.HandleFunc("GET /conditions/", resourcehandlers.GetConditions)
	mux.HandleFunc("GET /treatments/", resourcehandlers.GetTreatments)
	mux.HandleFunc("POST /relationship", relationshiphandler.CreateNewRelationship)
	mux.HandleFunc("GET /actors/details/", resourcehandlers.GetSpecificActor)
	mux.HandleFunc("GET /conditions/details/", resourcehandlers.GetSpecificCondition)
	mux.HandleFunc("GET /treatments/details/", resourcehandlers.GetSpecificTreatment)
	mux.HandleFunc("OPTIONS /", handleOptions)

	serveStatic()

	fmt.Println("Server started on localhost:8080")
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}

func serveStatic() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}
