package resourcehandlers

import (
	"context"
	"database/sql"
	"encoding/json"
	dbutils "fandm/internal/database/_dbutils"
	"fandm/internal/utls"
	utils "fandm/internal/utls"
	"fmt"
	"net/http"
	"time"
)

type SearchResult struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DetailedResult struct {
	Name                  string                 `json:"name"`
	Description           string                 `json:"description"`
	TreatmentInteractions map[string]interface{} `json:"treatmentInteractions"`
	ActorInteractions     map[string]interface{} `json:"actorInteractions"`
	ConditionInteractions map[string]interface{} `json:"conditionInteractions"`
}

func handeBasicSearch(routeType string, w http.ResponseWriter, r *http.Request) {

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

	var results []SearchResult

	if len(param) >= 4 {

		basicResultsSuccessful := make(chan bool)
		go handleBasicResources(query, basicResultsSuccessful, &results)

		resultSuccessful := <-basicResultsSuccessful

		if !resultSuccessful {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Query failed."))
			return
		}

		jsonResponse, err := json.Marshal(results)
		utls.Catch(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Query too short."))
		return
	}

}

func handleBasicResources(query string, done chan bool, results *[]SearchResult) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("mysql", dbutils.GetDbConnectionString())
	utls.Catch(err)
	defer db.Close()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("QueryContext timeout")
			done <- false
		}
		done <- false
		fmt.Println("Error: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name, description, firstInteractions, secondInteractions string
		err := rows.Scan(&id, &name, &description, &firstInteractions, &secondInteractions)
		*results = append(*results, SearchResult{Name: name, Description: description})
		utls.Catch(err)
	}
	done <- true

}

func handleSpecificSearch(routeType string, w http.ResponseWriter, r *http.Request) {
	var param string
	var query string
	switch routeType {
	case "actors":
		param = r.URL.Path[len("/actors/details/"):]
		query = "SELECT * FROM actors WHERE actor_name = '" + param + "'"
	case "treatments":
		param = r.URL.Path[len("/treatments/details/"):]
		query = "SELECT * FROM treatments WHERE treatment_name = '" + param + "'"
	case "conditions":
		param = r.URL.Path[len("/conditions/details/"):]
		query = "SELECT * FROM conditions WHERE condition_name = '" + param + "'"
	}

	var result DetailedResult
	specificResultsSuccessful := make(chan bool)

	go HandleDetailedResources(routeType, query, specificResultsSuccessful, &result)

	resultSuccessful := <-specificResultsSuccessful

	if !resultSuccessful {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Query failed."))
		return
	}

	jsonResponse, err := json.Marshal(result)
	utls.Catch(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func HandleDetailedResources(routeType string, query string, done chan bool, result *DetailedResult) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("mysql", dbutils.GetDbConnectionString())
	utls.Catch(err)
	defer db.Close()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("QueryContext timeout")
			done <- false
		}
		done <- false
		fmt.Println("Error: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name, description, firstInteractions, secondInteractions string
		err := rows.Scan(&id, &name, &description, &firstInteractions, &secondInteractions)
		switch routeType {
		case "actors":
			*result = DetailedResult{
				Name:                  name,
				Description:           description,
				TreatmentInteractions: utils.StringToMap(firstInteractions),
				ConditionInteractions: utils.StringToMap(secondInteractions),
			}
		case "treatments":
			*result = DetailedResult{
				Name:                  name,
				Description:           description,
				ActorInteractions:     utils.StringToMap(firstInteractions),
				ConditionInteractions: utils.StringToMap(secondInteractions),
			}
		case "conditions":
			*result = DetailedResult{
				Name:                  name,
				Description:           description,
				ActorInteractions:     utils.StringToMap(firstInteractions),
				TreatmentInteractions: utils.StringToMap(secondInteractions),
			}

		}

		utls.Catch(err)
	}
	done <- true
}

func GetTreatments(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	handeBasicSearch("treatments", w, r)
}

func GetConditions(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	handeBasicSearch("conditions", w, r)
}

func GetActors(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	handeBasicSearch("actors", w, r)
}

func GetSpecificTreatment(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	handleSpecificSearch("treatments", w, r)
}

func GetSpecificCondition(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	handleSpecificSearch("conditions", w, r)
}

func GetSpecificActor(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	handleSpecificSearch("actors", w, r)
}
