package relationshiphandler

import (
	"context"
	"database/sql"
	"encoding/json"
	dbutils "fandm/internal/database/_dbutils"
	utils "fandm/internal/utls"
	"fmt"
	"net/http"
	"time"
)

type Relationship struct {
	ReporterID        int32  `json:"reporter_id"`
	ReporterActor     string `json:"reporter_actor"`
	ReporterCondition string `json:"reporter_condition"`
	ReporterTreatment string `json:"reporter_treatment"`
}

func CreateNewRelationship(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	var incomingRelationship Relationship
	err := json.NewDecoder(r.Body).Decode(&incomingRelationship)
	if err != nil {
		http.Error(w, "No data sent.", http.StatusBadRequest)
		return
	}

	if incomingRelationship.ReporterID == 0 || incomingRelationship.ReporterActor == "" || incomingRelationship.ReporterCondition == "" || incomingRelationship.ReporterTreatment == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reportInsertDone := make(chan bool)
	newRelationshipDone := make(chan bool)

	go insertReport(&incomingRelationship, reportInsertDone)

	reportInserted := <-reportInsertDone
	if !reportInserted {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to insert data into database: User has already submitted this report."))
		return
	}

	go insertRelationship(&incomingRelationship, newRelationshipDone)

	relationshipInserted := <-newRelationshipDone

	if !relationshipInserted {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to insert data into database: Relationship insertion error."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully inserted data into database."))
}

func insertReport(incomingRelationship *Relationship, done chan bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := sql.Open("mysql", dbutils.GetDbConnectionString())
	if err != nil {
		done <- false
		return
	}
	defer db.Close()

	query, err := db.QueryContext(ctx, "SELECT 1 FROM reports WHERE reporter_id = ? LIMIT 1", incomingRelationship.ReporterID)
	if err != nil {
		done <- false
		return
	}
	defer query.Close()

	if query.Next() {
		fmt.Println("Relationship already exists in database.")
		done <- false
		return
	}

	statement, err := db.PrepareContext(ctx, "INSERT INTO reports (reporter_id, reporter_actor, reporter_condition, reporter_treatment) VALUES (?, ?, ?, ?)")
	if err != nil {
		done <- false
		return
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, incomingRelationship.ReporterID, incomingRelationship.ReporterActor, incomingRelationship.ReporterCondition, incomingRelationship.ReporterTreatment)
	if err != nil {
		done <- false
		return
	}

	fmt.Println("Successfully inserted relationship into database.")
	done <- true
}

func insertRelationship(incomingRelationship *Relationship, done chan bool) {

	var checkString string
	var queryString string
	var newInsertQueryString string
	var itemToUpdate string
	var firstRelationship string
	var secondRelationship string

	actorInserted := make(chan bool)
	conditionInserted := make(chan bool)
	treatmentInserted := make(chan bool)

	for _, entry := range []string{incomingRelationship.ReporterActor, incomingRelationship.ReporterCondition, incomingRelationship.ReporterTreatment} {

		// if actor is of type 'ReporterActor' then continue;
		if entry == incomingRelationship.ReporterActor {
			itemToUpdate = entry
			checkString = "SELECT treatment_interactions, condition_interactions FROM actors WHERE actor_name = ? LIMIT 1"
			queryString = "UPDATE actors SET treatment_interactions = ?, condition_interactions = ? WHERE actor_name = ?"
			firstRelationship = incomingRelationship.ReporterTreatment
			secondRelationship = incomingRelationship.ReporterCondition
			newInsertQueryString = "INSERT INTO actors (actor_name, treatment_interactions, condition_interactions) VALUES (?, ?, ?)"
			go performInsert(itemToUpdate, checkString, queryString, newInsertQueryString,
				firstRelationship, secondRelationship, actorInserted)
		}
		// if condition is of type 'ReporterCondition' then continue;
		if entry == incomingRelationship.ReporterCondition {
			itemToUpdate = entry
			checkString = "SELECT actor_interactions, treatment_interactions FROM conditions WHERE condition_name = ? LIMIT 1"
			queryString = "UPDATE conditions SET actor_interactions = ?, treatment_interactions = ? WHERE condition_name = ?"
			firstRelationship = incomingRelationship.ReporterActor
			secondRelationship = incomingRelationship.ReporterTreatment
			newInsertQueryString = "INSERT INTO conditions (condition_name, actor_interactions, treatment_interactions) VALUES (?, ?, ?)"
			go performInsert(itemToUpdate, checkString, queryString, newInsertQueryString,
				firstRelationship, secondRelationship, conditionInserted)
		}

		// if treatment is of type 'ReporterTreatment' then continue;
		if entry == incomingRelationship.ReporterTreatment {
			itemToUpdate = entry
			checkString = "SELECT actor_interactions, condition_interactions FROM treatments WHERE treatment_name = ? LIMIT 1"
			queryString = "UPDATE treatments SET actor_interactions = ?, condition_interactions = ? WHERE treatment_name = ?"
			firstRelationship = incomingRelationship.ReporterActor
			secondRelationship = incomingRelationship.ReporterCondition
			newInsertQueryString = "INSERT INTO treatments (treatment_name, actor_interactions, condition_interactions) VALUES (?, ?, ?)"
			go performInsert(itemToUpdate, checkString, queryString, newInsertQueryString,
				firstRelationship, secondRelationship, treatmentInserted)
		}
	}

	actorInsertDone := <-actorInserted
	conditionInsertDone := <-conditionInserted
	treatmentInsertDone := <-treatmentInserted

	if actorInsertDone && conditionInsertDone && treatmentInsertDone {
		done <- true
	} else {
		done <- false
	}
}

func performInsert(itemToUpdate string, checkString string, queryString string, newInsertQueryString string, firstRelationship string, secondRelationship string, done chan bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := sql.Open("mysql", dbutils.GetDbConnectionString())
	if err != nil {
		done <- false
		return
	}
	defer db.Close()

	query, err := db.QueryContext(ctx, checkString, itemToUpdate)
	if err != nil {
		done <- false
		return
	}
	defer query.Close()

	var firstInteractions, secondInteractions string
	if query.Next() {
		err = query.Scan(&firstInteractions, &secondInteractions)
		if err != nil {
			done <- false
			return
		}

		firstMap := utils.StringToMap(firstInteractions)
		secondMap := utils.StringToMap(secondInteractions)

		secondMap[secondRelationship] = secondMap[secondRelationship].(int) + 1
		firstMap[firstRelationship] = firstMap[firstRelationship].(int) + 1

		firstInteractions = utils.MapToString(firstMap)
		secondInteractions = utils.MapToString(secondMap)

		statement, err := db.PrepareContext(ctx, queryString)
		if err != nil {
			done <- false
			return
		}
		defer statement.Close()

		_, err = statement.ExecContext(ctx, firstInteractions, secondInteractions, itemToUpdate)
		if err != nil {
			done <- false
			return
		}
		done <- true
	} else {
		secondMap := map[string]interface{}{secondRelationship: 1}
		firstMap := map[string]interface{}{firstRelationship: 1}

		secondInteractions = utils.MapToString(secondMap)
		firstInteractions = utils.MapToString(firstMap)

		statement, err := db.PrepareContext(ctx, newInsertQueryString)
		if err != nil {
			done <- false
			return
		}
		defer statement.Close()

		_, err = statement.ExecContext(ctx, itemToUpdate, firstInteractions, secondInteractions)
		if err != nil {
			done <- false
			return
		}
		done <- true
	}
}
