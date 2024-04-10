package routehandlers

import (
	"context"
	"database/sql"
	"encoding/json"
	dbutils "fandm/internal/database/_dbutils"
	"fandm/internal/utls"
	"fmt"
	"net/http"
	"time"
)

func handleRouteMethods(routeType string, w http.ResponseWriter, r *http.Request) {
	var param string
	var query string

	switch routeType {
	case "actors":
		param = r.URL.Path[len("/actors/"):]
		query = "SELECT * FROM actors WHERE actor_description LIKE '%" + param + "%'"
	case "treatments":
		param = r.URL.Path[len("/treatments/"):]
		query = "SELECT * FROM treatments WHERE treatment_description LIKE '%" + param + "%'"
	case "conditions":
		param = r.URL.Path[len("/conditions/"):]
		query = "SELECT * FROM conditions WHERE condition_description LIKE '%" + param + "%'"
	}

	responseMap := make(map[string]string)

	if len(param) >= 4 {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		db, err := sql.Open("mysql", dbutils.GetDbConnectionString())
		utls.Catch(err)
		defer db.Close()

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("QueryContext timeout")
				return
			}
			fmt.Println("Error: ", err)
		}
		defer rows.Close()

		for rows.Next() {
			var id string
			var name, description, firstInteractions, secondInteractions string
			err := rows.Scan(&id, &name, &description, &firstInteractions, &secondInteractions)
			responseMap[name] = description
			utls.Catch(err)
		}
	}

	jsonResponse, err := json.Marshal(responseMap)
	utls.Catch(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func GetTreatments(w http.ResponseWriter, r *http.Request) {
	handleRouteMethods("treatments", w, r)
}

func GetConditions(w http.ResponseWriter, r *http.Request) {
	handleRouteMethods("conditions", w, r)
}

func GetActors(w http.ResponseWriter, r *http.Request) {
	handleRouteMethods("actors", w, r)
}
