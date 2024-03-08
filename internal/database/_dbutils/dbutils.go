package dbutils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"fandm/environment"

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
	if err != nil {
		fmt.Println("Failed to connect to Mysql Database with the following error: ", err)
		return false
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping the Mysql Database with the following error: ", err)
		return false
	}

	fmt.Println("Successfully connected to the Mysql Database")

	return true
}

func SetUpDatabase() {
	if TestConnection() {
		db, err := sql.Open("mysql", GetDbConnectionString())
		if err != nil {
			fmt.Println("Failed to connect to Mysql Database with the following error: ", err)
			return
		}
		// Create primary interactive tables.

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS treatments ( id INT AUTO_INCREMENT PRIMARY KEY, treatment_name TEXT NOT NULL, treatment_description TEXT NOT NULL, actort_interactions TEXT NOT NULL, condition_interactions TEXT NOT NULL);")
		if err != nil {
			fmt.Println("Failed to create Schema 'treatments'", err)
			return
		}
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS actors ( id INT AUTO_INCREMENT PRIMARY KEY, actor_name TEXT NOT NULL, actor_description TEXT NOT NULL, treatment_interactions TEXT NOT NULL, condition_interactions TEXT NOT NULL);")
		if err != nil {
			fmt.Println("Failed to create Schema 'actors'", err)
			return
		}
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS conditions ( id INT AUTO_INCREMENT PRIMARY KEY, condition_name TEXT NOT NULL, condition_description TEXT NOT NULL, actort_interactions TEXT NOT NULL, treatment_interactions TEXT NOT NULL);")
		if err != nil {
			fmt.Println("Failed to create Schema 'conditions'", err)
			return
		}
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS reports ( id INT AUTO_INCREMENT PRIMARY KEY, reporter_condition VARCHAR(255) NOT NULL, reporter_treatment VARCHAR(255) NOT NULL, reporter_actor VARCHAR(255) NOT NULL);")
		if err != nil {
			fmt.Println("Failed to create Schema 'reports'", err)
			return
		}

		createAndPopulateTables("foods")
		createAndPopulateTables("diseases")
		creareAndPopulateTables("treatments")
	}
}

func createAndPopulateTables(insertType string) {
	var filePath string

	switch insertType {
	case "foods":
		filePath = "internal/_json/actors/foods.json"
	case "diseases":
		filePath = "internal/_json/conditions/diseases.json"
	case "treatments":
		filePath = "internal/_json/treatments/treatments.json"
	}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	var result map[string]string

	json.Unmarshal([]byte(byteValue), &result)

	insertActors(result)
}

func insertActors(actorsMap map[string]string) {
	db, err := sql.Open("mysql", GetDbConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO actors(actor_name, actor_description, treatment_interactions, condition_interactions) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range actorsMap {
		_, err = stmt.Exec(key, value, "", "")
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Data inserted successfully.")
}

func createAndPopulateConditionsTable() {
}

func createAndPopulateTreatmentsTable() {
}
