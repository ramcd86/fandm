package relationshiphandler

import (
	"database/sql"
	"encoding/json"
	dbutils "fandm/internal/database/_dbutils"
	"fmt"
	"net/http"
	"strconv"
)

type Relationship struct {
	ReporterID        int32  `json:"reporter_id"` // This is the primary key
	ReporterActor     string `json:"reporter_actor"`
	ReporterCondition string `json:"reporter_condition"`
	ReporterTreatment string `json:"reporter_treatment"`
}

func CreateNewRelationship(w http.ResponseWriter, r *http.Request) {
	var incomingRelationship Relationship
	err := json.NewDecoder(r.Body).Decode(&incomingRelationship)
	if err != nil {
		http.Error(w, string("No data sent."), http.StatusBadRequest)
		return
	}

	var dataInserted bool
	dataInsertdone := make(chan bool)

	go insertReport(&incomingRelationship, dataInsertdone)

	dataInserted = <-dataInsertdone

	if !dataInserted {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func insertReport(incomingRelationship *Relationship, done chan bool) {
	if incomingRelationship.ReporterID == 0 || incomingRelationship.ReporterActor == "" || incomingRelationship.ReporterCondition == "" || incomingRelationship.ReporterTreatment == "" {
		done <- false
		return
	}

	db, err := sql.Open("mysql", dbutils.GetDbConnectionString())
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	query, err := db.Query("SELECT 1 FROM reports WHERE reporter_id = " + strconv.Itoa(int(incomingRelationship.ReporterID)) + " LIMIT 1")
	if err != nil {
		panic(err.Error())
	}

	if query.Next() {
		fmt.Println("Relationship already exists in database.")
		done <- false
		return
	}

	statement, err := db.Prepare("INSERT INTO reports (reporter_id, reporter_actor, reporter_condition, reporter_treatment) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = statement.Exec(incomingRelationship.ReporterID, incomingRelationship.ReporterActor, incomingRelationship.ReporterCondition, incomingRelationship.ReporterTreatment)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Successfully inserted relationship into database.")

	done <- true
}
