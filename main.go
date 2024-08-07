package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("file.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		ErrorsToDocker(err, "Error loading .env file:")
		return
	}
	dbConfig := DBConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	}
	err = SendTestData(dbConfig)
	if err != nil {
		fmt.Println(err)
	}

	// Установка обработчика для маршрута "/"
	/*http.HandleFunc("/home", saveDataHandler(dbConfig))

	// Запуск веб-сервера на порту 80
	log.Fatal(http.ListenAndServe(":80", nil))*/
}

func saveDataHandler(dbConfig DBConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получение параметров из тела запроса
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing form: %s", err), http.StatusBadRequest)
			ErrorsToDocker(err, "Error parsing form:")
			return
		}

		temperature := r.Form.Get("temperature")
		humidity := r.Form.Get("humidity")
		err = PostgresInsert(dbConfig, temperature, humidity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			w.Write([]byte("An error occured"))
		}
		w.Write([]byte("Data successfully saved to the database."))
	}
}
