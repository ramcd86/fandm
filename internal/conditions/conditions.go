package conditions

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

func GetConditions(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/conditions/"):]
	responseMap := make(map[string]string)
	query := "SELECT * FROM actors WHERE condition_description LIKE '%" + param + "%'"

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
			var conditionName, conditionDescription, actorInteractions, treatmentInteractions string
			err := rows.Scan(&id, &conditionName, &conditionDescription, &actorInteractions, &treatmentInteractions)
			responseMap[conditionName] = conditionDescription
			utls.Catch(err)
		}
	}

	jsonResponse, err := json.Marshal(responseMap)
	utls.Catch(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
