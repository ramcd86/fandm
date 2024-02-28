package dbutils

import (
	"database/sql"
	"fmt"

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

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS treatments ( id INT AUTO_INCREMENT PRIMARY KEY, treatment_name TEXT NOT NULL, treatment_description TEXT NOT NULL, actort_interactions TEXT NOT NULL, condition_interactions TEXT NOT NULL, uuid VARCHAR(255) NOT NULL );")
		if err != nil {
			fmt.Println("Failed to create Schema 'treatments'", err)
			return
		}
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS actors ( id INT AUTO_INCREMENT PRIMARY KEY, actor_name TEXT NOT NULL, actor_description TEXT NOT NULL, treatment_interactions TEXT NOT NULL, condition_interactions TEXT NOT NULL, uuid VARCHAR(255) NOT NULL );")
		if err != nil {
			fmt.Println("Failed to create Schema 'actors'", err)
			return
		}
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS conditions ( id INT AUTO_INCREMENT PRIMARY KEY, condition_name TEXT NOT NULL, condition_description TEXT NOT NULL, actort_interactions TEXT NOT NULL, treatment_interactions TEXT NOT NULL, uuid VARCHAR(255) NOT NULL );")
		if err != nil {
			fmt.Println("Failed to create Schema 'conditions'", err)
			return
		}
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS reports ( id INT AUTO_INCREMENT PRIMARY KEY, reporter_condition VARCHAR(255) NOT NULL, reporter_treatment VARCHAR(255) NOT NULL, reporter_actor VARCHAR(255) NOT NULL, uuid VARCHAR(255) NOT NULL );")
		if err != nil {
			fmt.Println("Failed to create Schema 'reports'", err)
			return
		}
	}
}
