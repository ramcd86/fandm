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
	wg := sync.WaitGroup{}

	grabData := func(filePath string, insertQuery string) {
		defer wg.Done()

		var result map[string]string

		jsonFile, err := os.Open(filePath)
		utls.Catch(err)

		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)

		json.Unmarshal([]byte(byteValue), &result)

		insertActors(result, insertQuery)
	}

	switch insertType {
	case "foods":
		wg.Add(1)
		go grabData("internal/_json/actors/foods.json",
			`
            INSERT INTO actors
                (actor_name, actor_description, treatment_interactions, condition_interactions)
            VALUES(?, ?, ?, ?)
        `)
	case "diseases":
		wg.Add(1)
		go grabData("internal/_json/conditions/diseases.json",
			`
            INSERT INTO conditions
                (condition_name, condition_description, actor_interactions, treatment_interactions)
            VALUES(?, ?, ?, ?)
        `)
	case "treatments":
		wg.Add(1)
		go grabData("internal/_json/treatments/treatments.json",
			`
            INSERT INTO 
                treatments(treatment_name, treatment_description, actor_interactions, condition_interactions)
            VALUES(?, ?, ?, ?)
        `)
	}

	wg.Wait()
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
