package relationshiphandler

import (
	"context"
	"database/sql"
	"encoding/json"
	dbutils "fandm/internal/database/_dbutils"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Relationship struct {
	ReporterID        int32  `json:"reporter_id"`
	ReporterActor     string `json:"reporter_actor"`
	ReporterCondition string `json:"reporter_condition"`
	ReporterTreatment string `json:"reporter_treatment"`
}

func stringToMap(str string) map[string]int { // HL
	items := strings.Split(str, ",")
	m := make(map[string]int)
	for _, item := range items {
		parts := strings.Split(item, ":")
		if len(parts) == 2 {
			key := parts[0]
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				value = 0
			}
			m[key] = value
		}
	}
	return m
}

func mapToString(m map[string]int) string {
	str := ""
	for k, v := range m {
		str += k + ":" + strconv.Itoa(v) + ","
	}
	return str
}

func CreateNewRelationship(w http.ResponseWriter, r *http.Request) {
	var incomingRelationship Relationship
	err := json.NewDecoder(r.Body).Decode(&incomingRelationship)
	if err != nil {
		http.Error(w, string("No data sent."), http.StatusBadRequest)
		return
	}

	if incomingRelationship.ReporterID == 0 || incomingRelationship.ReporterActor == "" || incomingRelationship.ReporterCondition == "" || incomingRelationship.ReporterTreatment == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var dataInserted bool
	var actorInserted bool
	// var conditionInserted bool
	// var treatmentInserted bool

	actorInsertdone := make(chan bool)
	dataInsertdone := make(chan bool)
	// conditionInsertdone := make(chan bool)
	// treatmentInsertdone := make(chan bool)

	go insertReport(&incomingRelationship, dataInsertdone)
	go insertActor(&incomingRelationship, actorInsertdone)

	dataInserted = <-dataInsertdone

	if !dataInserted || !actorInserted {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func insertReport(incomingRelationship *Relationship, done chan bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := sql.Open("mysql", dbutils.GetDbConnectionString())
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query, err := db.QueryContext(ctx, "SELECT 1 FROM reports WHERE reporter_id = "+strconv.Itoa(int(incomingRelationship.ReporterID))+" LIMIT 1")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("QueryContext timeout")
			done <- false
			return
		}
		panic(err.Error())
	}

	if query.Next() {
		fmt.Println("Relationship already exists in database.")
		done <- false
		return
	}

	statement, err := db.PrepareContext(ctx, "INSERT INTO reports (reporter_id, reporter_actor, reporter_condition, reporter_treatment) VALUES (?, ?, ?, ?)")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("QueryContext timeout")
			done <- false
			return
		}
		panic(err.Error())
	}

	_, err = statement.ExecContext(ctx, incomingRelationship.ReporterID, incomingRelationship.ReporterActor, incomingRelationship.ReporterCondition, incomingRelationship.ReporterTreatment)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("QueryContext timeout")
			done <- false
			return
		}
		panic(err.Error())
	}

	fmt.Println("Successfully inserted relationship into database.")

	done <- true
}

func insertActor(incomingRelationship *Relationship, done chan bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := sql.Open("mysql", dbutils.GetDbConnectionString())
	if err != nil {
		done <- false
	}
	defer db.Close()

	query, err := db.QueryContext(ctx, "SELECT * FROM actors WHERE actor_name = '"+incomingRelationship.ReporterActor+"' LIMIT 1")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("QueryContext timeout")
		}
		done <- false
	}
	defer query.Close()

	for query.Next() {
		var id string
		var name, description, firstInteractions, secondInteractions string
		err := query.Scan(&id, &name, &description, &firstInteractions, &secondInteractions)
		if err != nil {
			fmt.Println(err)
			done <- false
		}
		firstInteractionsMap := stringToMap(firstInteractions)
		secondInteractionsMap := stringToMap(secondInteractions)
		firstInteractionsMap[incomingRelationship.ReporterCondition]++
		secondInteractionsMap[incomingRelationship.ReporterTreatment]++
		firstInteractions = mapToString(firstInteractionsMap)
		secondInteractions = mapToString(secondInteractionsMap)
		statement, err := db.PrepareContext(ctx, "UPDATE actors SET treatment_interactions = ?, condition_interactions = ? WHERE actor_name = ?")
		if err != nil {
			fmt.Println(err)
			done <- false
		}
		_, err = statement.ExecContext(ctx, secondInteractions, firstInteractions, incomingRelationship.ReporterActor)
		if err != nil {
			fmt.Println(err)
			done <- false
		}
		done <- true
	}

	done <- true
}
