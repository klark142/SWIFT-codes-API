package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"swift-codes/internal/db"
	"swift-codes/internal/model"

	"github.com/gorilla/mux"
)


func GetSwiftCodeHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		swiftCodeParam := vars["swiftCode"]

		swiftData, err := db.GetSwiftCode(dbConn, swiftCodeParam)
		if err != nil {
			http.Error(w, "Nie znaleziono wpisu", http.StatusNotFound)
			return
		}

		if swiftData.IsHeadquarter {
			branches, err := db.GetBranchesByHeadquarter(dbConn, swiftData.SwiftCode)
			if err != nil {
				http.Error(w, "Błąd podczas pobierania oddziałów", http.StatusInternalServerError)
				return
			}
			swiftData.Branches = branches
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(swiftData); err != nil {
			http.Error(w, "Błąd podczas kodowania odpowiedzi", http.StatusInternalServerError)
			return
		}
	}
}

func GetSwiftCodesByCountryHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		countryISO2 := strings.ToUpper(vars["countryISO2code"])

		swiftCodes, err := db.GetSwiftCodesByCountry(dbConn, countryISO2)
		if err != nil {
			http.Error(w, "Błąd pobierania danych", http.StatusInternalServerError)
			return
		}

		countryName := ""
		if len(swiftCodes) > 0 {
			countryName = swiftCodes[0].CountryName
		}

		response := struct {
			CountryISO2 string             `json:"countryISO2"`
			CountryName string             `json:"countryName"`
			SwiftCodes  []model.SwiftCode  `json:"swiftCodes"`
		}{
			CountryISO2: countryISO2,
			CountryName: countryName,
			SwiftCodes:  swiftCodes,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Błąd podczas kodowania odpowiedzi", http.StatusInternalServerError)
			return
		}
	}
}

func CreateSwiftCodeHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newSwift model.SwiftCode

		if err := json.NewDecoder(r.Body).Decode(&newSwift); err != nil {
			http.Error(w, "Błędny format danych", http.StatusBadRequest)
			return
		}

		newSwift.CountryISO2 = strings.ToUpper(strings.TrimSpace(newSwift.CountryISO2))
		newSwift.CountryName = strings.ToUpper(strings.TrimSpace(newSwift.CountryName))
		newSwift.BankName = strings.ToUpper(strings.TrimSpace(newSwift.BankName))
		newSwift.SwiftCode = strings.TrimSpace(newSwift.SwiftCode)
		newSwift.Address = strings.TrimSpace(newSwift.Address)

		if err := db.InsertSwiftCode(dbConn, newSwift); err != nil {
			http.Error(w, "Nie udało się dodać wpisu", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"message": "Wpis dodany pomyślnie",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Błąd podczas kodowania odpowiedzi", http.StatusInternalServerError)
			return
		}
	}
}

func DeleteSwiftCodeHandler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		swiftCodeParam := vars["swift-code"]

		if err := db.DeleteSwiftCode(dbConn, swiftCodeParam); err != nil {
			http.Error(w, "Nie udało się usunąć wpisu", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"message": "Wpis usunięty pomyślnie",
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Błąd podczas kodowania odpowiedzi", http.StatusInternalServerError)
			return
		}
	}
}
