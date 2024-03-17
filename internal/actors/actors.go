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

	// sql query where we SELECT FROM ACTORS in field actor_description where LIKE is param.
	// return the result
	// var query string = "SELECT * FROM actors WHERE actor_description LIKE '%" + param + "%'"

	fmt.Fprintf(w, "The value of 'myParam' is: %s", param)
}
