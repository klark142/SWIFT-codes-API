package main

import (
	"log"
	"net/http"
	"os"
	"swift-codes/internal/db"
	"swift-codes/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		log.Fatal("Brak ustawionej zmiennej środowiskowej DB_CONN")
	}

	database, err := db.InitDB(connStr)
	if err != nil {
		log.Fatalf("Błąd inicjalizacji bazy danych: %v", err)
	}
	defer database.Close()

	router := mux.NewRouter()

	router.HandleFunc("/v1/swift-codes/{swiftCode}", handlers.GetSwiftCodeHandler(database)).Methods("GET")
	router.HandleFunc("/v1/swift-codes/country/{countryISO2code}", handlers.GetSwiftCodesByCountryHandler(database)).Methods("GET")
	router.HandleFunc("/v1/swift-codes", handlers.CreateSwiftCodeHandler(database)).Methods("POST")
	router.HandleFunc("/v1/swift-codes/{swift-code}", handlers.DeleteSwiftCodeHandler(database)).Methods("DELETE")

	log.Println("Serwer uruchomiony na porcie 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
