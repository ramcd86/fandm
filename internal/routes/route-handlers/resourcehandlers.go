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
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	FirstInteractions  map[string]interface{} `json:"firstInteractions"`
	SecondInteractions map[string]interface{} `json:"secondInteractions"`
}

func handeBasicSearch(routeType string, w http.ResponseWriter, r *http.Request) {

fmt.Println("routeType: ", routeType)

	var param string
	var query string

	switch routeType {
	case "actors":
		param = r.URL.Path[len("/actors/"):]
		query = "SELECT * FROM actors WHERE actor_description LIKE '%" + param + "%'"
		fmt.Println("query: ", query)
	case "treatments":
		param = r.URL.Path[len("/treatments/"):]
		query = "SELECT * FROM treatments WHERE treatment_description LIKE '%" + param + "%'"
	case "conditions":
		param = r.URL.Path[len("/conditions/"):]
		query = "SELECT * FROM conditions WHERE condition_description LIKE '%" + param + "%'"
	}

	var results []SearchResult

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
			results = append(results, SearchResult{Name: name, Description: description})
			utls.Catch(err)
		}
	}

	fmt.Println("results: ", results)

	jsonResponse, err := json.Marshal(results)
	utls.Catch(err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
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

	go HandleDetailedResources(query, specificResultsSuccessful, &result)

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

func HandleDetailedResources(query string, done chan bool, result *DetailedResult) {

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
		*result = DetailedResult{
			Name:               name,
			Description:        description,
			FirstInteractions:  utils.StringToMap(firstInteractions),
			SecondInteractions: utils.StringToMap(secondInteractions),
		}
		utls.Catch(err)
	}
	done <- true
}

func GetTreatments(w http.ResponseWriter, r *http.Request) {
	handeBasicSearch("treatments", w, r)
}

func GetConditions(w http.ResponseWriter, r *http.Request) {
	handeBasicSearch("conditions", w, r)
}

func GetActors(w http.ResponseWriter, r *http.Request) {
	handeBasicSearch("actors", w, r)
}

func GetSpecificTreatment(w http.ResponseWriter, r *http.Request) {
	handleSpecificSearch("treatments", w, r)
}

func GetSpecificCondition(w http.ResponseWriter, r *http.Request) {
	handleSpecificSearch("conditions", w, r)
}

func GetSpecificActor(w http.ResponseWriter, r *http.Request) {
	handleSpecificSearch("actors", w, r)
}
