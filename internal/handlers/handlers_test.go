package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"swift-codes/internal/db"
	"swift-codes/internal/model"
)

func setupTestServer(t *testing.T) (*mux.Router, *sql.DB) {
	connStr := os.Getenv("TEST_DB_CONN")
	if connStr == "" {
		t.Fatal("Brak ustawionej zmiennej środowiskowej TEST_DB_CONN")
	}
	testDB, err := db.InitDB(connStr)
	if err != nil {
		t.Fatalf("InitDB nie powiodło się: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/v1/swift-codes/{swiftCode}", GetSwiftCodeHandler(testDB)).Methods("GET")
	router.HandleFunc("/v1/swift-codes/country/{countryISO2code}", GetSwiftCodesByCountryHandler(testDB)).Methods("GET")
	router.HandleFunc("/v1/swift-codes", CreateSwiftCodeHandler(testDB)).Methods("POST")
	router.HandleFunc("/v1/swift-codes/{swift-code}", DeleteSwiftCodeHandler(testDB)).Methods("DELETE")

	return router, testDB
}

func clearTable(db *sql.DB, t *testing.T) {
	_, err := db.Exec("TRUNCATE TABLE swift_codes")
	if err != nil {
		t.Fatalf("Nie udało się wyczyścić tabeli: %v", err)
	}
}

func TestCreateAndGetSwiftCodeHandler(t *testing.T) {
	router, testDB := setupTestServer(t)
	defer testDB.Close()
	clearTable(testDB, t)

	testRecord := model.SwiftCode{
		BankName:      "API TEST BANK",
		Address:       "API TEST ADDRESS",
		CountryISO2:   "AT",
		CountryName:   "TESTLAND",
		IsHeadquarter: true,
		SwiftCode:     "APITESTXXX",
	}
	payload, err := json.Marshal(testRecord)
	if err != nil {
		t.Fatalf("Błąd kodowania JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/v1/swift-codes", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Błąd tworzenia żądania POST: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("POST - oczekiwano status 200, otrzymano %d", status)
	}

	req, err = http.NewRequest("GET", "/v1/swift-codes/APITESTXXX", nil)
	if err != nil {
		t.Fatalf("Błąd tworzenia żądania GET: %v", err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET - oczekiwano status 200, otrzymano %d", status)
	}

	var got model.SwiftCode
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("Błąd dekodowania odpowiedzi GET: %v", err)
	}
	if got.BankName != testRecord.BankName {
		t.Errorf("Oczekiwano BankName %s, otrzymano %s", testRecord.BankName, got.BankName)
	}
}

func TestGetSwiftCodesByCountryHandler(t *testing.T) {
	router, testDB := setupTestServer(t)
	defer testDB.Close()
	clearTable(testDB, t)

	records := []model.SwiftCode{
		{
			BankName:      "Bank One",
			Address:       "Address One",
			CountryISO2:   "AA",
			CountryName:   "CountryA",
			IsHeadquarter: true,
			SwiftCode:     "AA_SWIFTXXX1",
		},
		{
			BankName:      "Bank Two",
			Address:       "Address Two",
			CountryISO2:   "BB",
			CountryName:   "CountryB",
			IsHeadquarter: true,
			SwiftCode:     "BB_SWIFTXXX1",
		},
		{
			BankName:      "Bank Three",
			Address:       "Address Three",
			CountryISO2:   "AA",
			CountryName:   "CountryA",
			IsHeadquarter: false,
			SwiftCode:     "AA_SWIFTXXX2",
		},
	}

	for _, rec := range records {
		if err := db.InsertSwiftCode(testDB, rec); err != nil {
			t.Fatalf("InsertSwiftCode nie powiodło się: %v", err)
		}
	}

	req, err := http.NewRequest("GET", "/v1/swift-codes/country/aa", nil)
	if err != nil {
		t.Fatalf("Błąd tworzenia żądania GET: %v", err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GET by country - oczekiwano status 200, otrzymano %d", status)
	}

	var response struct {
		CountryISO2 string            `json:"countryISO2"`
		CountryName string            `json:"countryName"`
		SwiftCodes  []model.SwiftCode `json:"swiftCodes"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Błąd dekodowania odpowiedzi GET by country: %v", err)
	}
	if !strings.EqualFold(response.CountryISO2, "AA") {
		t.Errorf("Oczekiwano CountryISO2 'AA', otrzymano '%s'", response.CountryISO2)
	}
	if len(response.SwiftCodes) != 2 {
		t.Errorf("Oczekiwano 2 rekordów dla kraju AA, otrzymano %d", len(response.SwiftCodes))
	}
}

func TestDeleteSwiftCodeHandler(t *testing.T) {
	router, testDB := setupTestServer(t)
	defer testDB.Close()
	clearTable(testDB, t)

	rec := model.SwiftCode{
		BankName:      "Delete Bank",
		Address:       "Delete Address",
		CountryISO2:   "DL",
		CountryName:   "Deleteland",
		IsHeadquarter: true,
		SwiftCode:     "DELETESWIFTXXX",
	}
	if err := db.InsertSwiftCode(testDB, rec); err != nil {
		t.Fatalf("InsertSwiftCode nie powiodło się: %v", err)
	}

	req, err := http.NewRequest("DELETE", "/v1/swift-codes/DELETESWIFTXXX", nil)
	if err != nil {
		t.Fatalf("Błąd tworzenia żądania DELETE: %v", err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DELETE - oczekiwano status 200, otrzymano %d", status)
	}

	req, err = http.NewRequest("GET", "/v1/swift-codes/DELETESWIFTXXX", nil)
	if err != nil {
		t.Fatalf("Błąd tworzenia żądania GET po DELETE: %v", err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("GET po DELETE - oczekiwano status 404, otrzymano %d", status)
	}
}
