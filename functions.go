package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
)

func EqualsCurrentDate(date string) bool {
	currentDate := time.Now()
	formatedDate := currentDate.Format("2006-01-02")
	//что-то типа
	if date != formatedDate {
		return false
	} else {
		return true
	}
}

func ErrorsToDocker(err error, description string) {
	dt := time.Now()
	if err != nil {
		fmt.Fprintln(os.Stderr, dt.String(), description, err)
		return
	}
}

func PostgresInsert(databaseConfig DBConfig, temperature string, humidity string) error {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		databaseConfig.Username, databaseConfig.Password, databaseConfig.Database, databaseConfig.Host, databaseConfig.Port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		ErrorsToDocker(err, "Failed to connect to the database:")
		return err
	}
	defer db.Close()
	date := time.Now()
	formatedDate := date.Format("2006-01-02")
	if IfTableExists(formatedDate, db, databaseConfig) {
		tmp := InsertIntoTable(formatedDate, db, temperature, humidity)
		if tmp != nil {
			ErrorsToDocker(tmp, "InsertIntoTable failed:")
			return tmp
		}
	} else {
		tmp := CreateTable(db, formatedDate)
		if tmp != nil {
			ErrorsToDocker(tmp, "Couldn't create table "+formatedDate)
			return tmp
		}
		tmp = InsertIntoTable(formatedDate, db, temperature, humidity)
		if tmp != nil {
			ErrorsToDocker(tmp, "InsertIntoTable failed:")
			return tmp
		}
	}
	return nil
}

func IfTableExists(tableName string, db *sql.DB, dbConfig DBConfig) bool {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM   information_schema.tables
			WHERE  table_schema = $1
			AND    table_name = $2
		);
		`
	err := db.QueryRow(query, "table_schema", tableName).Scan(&exists)
	if err != nil {
		ErrorsToDocker(err, "Error checking if table exists:")
		return false
	}
	return true
}

func CreateTable(db *sql.DB, tableName string) error {
	createTableQuery := fmt.Sprintf(`
		CREATE TABLE "%s" (
			id SERIAL PRIMARY KEY, 
			temperature VARCHAR(255),
			humidity VARCHAR(255),
			date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`, tableName)

	_, err := db.Exec(createTableQuery)
	if err != nil {
		ErrorsToDocker(err, fmt.Sprintf("Couldn't create table %s", tableName))
		return err
	}
	return nil
}

func InsertIntoTable(tableName string, db *sql.DB, values ...interface{}) error {
	valuePlaceHolders := make([]string, len(values))
	for i := range values {
		valuePlaceHolders[i] = fmt.Sprintf("$%d", i+1)
	}
	res := strings.Join(valuePlaceHolders, ", ")
	query := fmt.Sprintf("INSERT INTO %s (temperature, humidity) VALUES (%s)", tableName, res)
	_, err := db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}
	return nil
}

// временная функция, чтобы не делать запросы в Postman'e
func SendTestData(dbConfig DBConfig) error {
	temperature := "23"
	humidity := "59"
	err := PostgresInsert(dbConfig, temperature, humidity)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	return nil
}

func JsonAnswer() {

}
