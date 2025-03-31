package main

import (
	"log"
	"os"
	"swift-codes/internal/db"
	"swift-codes/internal/parser"
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

	filePath := "data/swiftcodes_data.csv"
	swiftRecords, err := parser.ParseCSV(filePath)
	if err != nil {
		log.Fatalf("Błąd parsowania CSV: %v", err)
	}

	for _, record := range swiftRecords {
		if err := db.InsertSwiftCode(database, record); err != nil {
			log.Printf("Błąd wstawiania rekordu %s: %v", record.SwiftCode, err)
		}
	}
	log.Println("Dane z pliku CSV zostały wstawione do bazy danych.")
}
