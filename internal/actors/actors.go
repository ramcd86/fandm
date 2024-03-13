package actors

import (
	"fmt"
	"net/http"
)

func GetActors(w http.ResponseWriter, r *http.Request) {
	// Get all actors from the database
	param := r.URL.Path[len("/actors/"):]
	fmt.Println(param)

	if len(param) >= 4 {
		// do query
	}

	fmt.Fprintf(w, "The value of 'myParam' is: %s", param)
}
