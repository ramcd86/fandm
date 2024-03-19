package actors

import (
	"database/sql"
	"encoding/json"
	dbutils "fandm/internal/database/_dbutils"
	"fandm/internal/utls"
	"net/http"
)

func GetActors(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/actors/"):]
	responseMap := make(map[string]string)
	query := "SELECT * FROM actors WHERE actor_description LIKE '%" + param + "%'"

	if len(param) >= 4 {

		db, err := sql.Open("mysql", dbutils.GetDbConnectionString())
		utls.Catch(err)
		defer db.Close()

		rows, err := db.Query(query)
		utls.Catch(err)
		defer rows.Close()

		for rows.Next() {
			var id string
			var actorName, actorDescription, treatmentInteractions, conditionInteractions string
			err := rows.Scan(&id, &actorName, &actorDescription, &treatmentInteractions, &conditionInteractions)
			responseMap[actorName] = actorDescription
			utls.Catch(err)
		}
	}

	jsonResponse, err := json.Marshal(responseMap)
	utls.Catch(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
	w.WriteHeader(http.StatusOK)
}
