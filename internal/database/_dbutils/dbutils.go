package dbutils

import (
	"database/sql"
	"encoding/json"
	"fandm/environment"
	"fandm/internal/utls"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

func GetDbConnectionString() string {
	fmt.Println(environment.DbConnection)
	return environment.DbConnection["DB_USER"] +
		":" + environment.DbConnection["DB_PASS"] +
		"@tcp(" + environment.DbConnection["DB_HOST"] +
		":" + environment.DbConnection["DB_PORT"] +
		")/" + environment.DbConnection["DB_NAME"]
}

func TestConnection() bool {
	db, err := sql.Open("mysql", GetDbConnectionString())
	utls.Catch(err)
	err = db.Ping()
	utls.Catch(err)
	fmt.Println("Successfully connected to the Mysql Database")

	return true
}

func SetUpDatabase() {
	if TestConnection() {
		db, err := sql.Open("mysql", GetDbConnectionString())
		utls.Catch(err)
		// Create primary interactive tables.
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS treatments 
			( id INT AUTO_INCREMENT PRIMARY KEY, 
			treatment_name TEXT NOT NULL, 
			treatment_description TEXT NOT NULL, 
			actor_interactions TEXT NOT NULL, 
			condition_interactions TEXT NOT NULL)`)
		utls.Catch(err)

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS actors 
			( id INT AUTO_INCREMENT PRIMARY KEY, 
			actor_name TEXT NOT NULL, 
			actor_description TEXT NOT NULL, 
			treatment_interactions TEXT NOT NULL, 
			condition_interactions TEXT NOT NULL);`)
		utls.Catch(err)

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS conditions 
			( id INT AUTO_INCREMENT PRIMARY KEY, 
			condition_name TEXT NOT NULL, 
			condition_description TEXT NOT NULL, 
			actor_interactions TEXT NOT NULL, 
			treatment_interactions TEXT NOT NULL);`)
		utls.Catch(err)

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS reports 
			( id INT AUTO_INCREMENT PRIMARY KEY, 
			reporter_condition VARCHAR(255) NOT NULL, 
			reporter_treatment VARCHAR(255) NOT NULL, 
			reporter_actor VARCHAR(255) NOT NULL);`)
		utls.Catch(err)

		createAndPopulateTables("foods")
		createAndPopulateTables("diseases")
		createAndPopulateTables("treatments")
	}
}

func createAndPopulateTables(insertType string) {
	var dataExists bool

	checkData := func(checkType string, done chan bool) {
		fmt.Println("Checking data for " + checkType)
		db, err := sql.Open("mysql", GetDbConnectionString())
		utls.Catch(err)
		defer db.Close()

		rows, err := db.Query("SELECT 1 FROM " + checkType + " LIMIT 1")
		utls.Catch(err)
		defer rows.Close()

		dataExists := rows.Next()

		done <- dataExists
	}

	dataCheckDone := make(chan bool)

	switch insertType {
	case "foods":
		go checkData("actors", dataCheckDone)
	case "diseases":
		go checkData("conditions", dataCheckDone)
	case "treatments":
		go checkData("treatments", dataCheckDone)
	}

	dataExists = <-dataCheckDone

	var filePath, insertQuery string

	switch insertType {
	case "foods":
		filePath = "internal/_json/actors/foods.json"
		insertQuery = `
            INSERT INTO actors
            (actor_name, actor_description, treatment_interactions, condition_interactions)
            VALUES(?, ?, ?, ?)
        `
	case "diseases":
		filePath = "internal/_json/conditions/diseases.json"
		insertQuery = `
            INSERT INTO conditions
            (condition_name, condition_description, actor_interactions, treatment_interactions)
            VALUES(?, ?, ?, ?)
        `
	case "treatments":
		filePath = "internal/_json/treatments/treatments.json"
		insertQuery = `
            INSERT INTO 
            treatments(treatment_name, treatment_description, actor_interactions, condition_interactions)
            VALUES(?, ?, ?, ?)
        `
	}

	if !dataExists {
		insertWaitGroup := sync.WaitGroup{}
		insertWaitGroup.Add(1)

		grabData := func(filePath string, insertQuery string) {
			defer insertWaitGroup.Done()

			var result map[string]string

			jsonFile, err := os.Open(filePath)
			utls.Catch(err)
			defer jsonFile.Close()

			byteValue, _ := ioutil.ReadAll(jsonFile)

			json.Unmarshal([]byte(byteValue), &result)

			insertActors(result, insertQuery)
		}

		go grabData(filePath, insertQuery)

		insertWaitGroup.Wait()
	} else {
		fmt.Println(insertType + " data already exists")
	}
}

func insertActors(dataMap map[string]string, insertQuery string) {
	db, err := sql.Open("mysql", GetDbConnectionString())
	utls.Catch(err)

	defer db.Close()

	stmt, err := db.Prepare(insertQuery)
	utls.Catch(err)

	for key, value := range dataMap {
		_, err = stmt.Exec(key, value, "", "")
		utls.Catch(err)
	}

	fmt.Println("Data inserted successfully.")
}
