package db

import (
	"database/sql"
	"os"
	"testing"

	"swift-codes/internal/model"
)

func getTestDB(t *testing.T) *sql.DB {
	connStr := os.Getenv("TEST_DB_CONN")
	if connStr == "" {
		t.Fatal("Brak ustawionej zmiennej środowiskowej TEST_DB_CONN")
	}
	db, err := InitDB(connStr)
	if err != nil {
		t.Fatalf("InitDB nie powiodło się: %v", err)
	}
	return db
}

func clearTable(db *sql.DB, t *testing.T) {
	_, err := db.Exec("TRUNCATE TABLE swift_codes")
	if err != nil {
		t.Fatalf("Nie udało się wyczyścić tabeli: %v", err)
	}
}

func TestInsertAndGetSwiftCode(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	clearTable(db, t)

	testRecord := model.SwiftCode{
		BankName:      "Test Bank",
		Address:       "Test Address",
		CountryISO2:   "TS",
		CountryName:   "Testland",
		IsHeadquarter: true,
		SwiftCode:     "TESTSWIFTXXX",
	}

	if err := InsertSwiftCode(db, testRecord); err != nil {
		t.Fatalf("InsertSwiftCode nie powiodło się: %v", err)
	}

	retrieved, err := GetSwiftCode(db, testRecord.SwiftCode)
	if err != nil {
		t.Fatalf("GetSwiftCode nie powiodło się: %v", err)
	}
	if retrieved.BankName != testRecord.BankName {
		t.Errorf("Oczekiwano BankName %s, otrzymano %s", testRecord.BankName, retrieved.BankName)
	}
}

func TestGetBranchesByHeadquarter(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	clearTable(db, t)

	headquarter := model.SwiftCode{
		BankName:      "HQ Bank",
		Address:       "HQ Address",
		CountryISO2:   "HQ",
		CountryName:   "HQland",
		IsHeadquarter: true,
		SwiftCode:     "HQSWIFTXXXT",
	}
	if err := InsertSwiftCode(db, headquarter); err != nil {
		t.Fatalf("InsertSwiftCode dla głównej siedziby nie powiodło się: %v", err)
	}

	branch := model.SwiftCode{
		BankName:      "Branch Bank",
		Address:       "Branch Address",
		CountryISO2:   "HQ",
		CountryName:   "HQland",
		IsHeadquarter: false,
		SwiftCode:     "HQSWIFTXXXB",
	}
	if err := InsertSwiftCode(db, branch); err != nil {
		t.Fatalf("InsertSwiftCode dla oddziału nie powiodło się: %v", err)
	}

	branches, err := GetBranchesByHeadquarter(db, headquarter.SwiftCode)
	if err != nil {
		t.Fatalf("GetBranchesByHeadquarter nie powiodło się: %v", err)
	}

	if len(branches) != 1 {
		t.Errorf("Oczekiwano 1 oddziału, otrzymano %d", len(branches))
	}

	if branches[0].SwiftCode != branch.SwiftCode {
		t.Errorf("Oczekiwano SwiftCode oddziału %s, otrzymano %s", branch.SwiftCode, branches[0].SwiftCode)
	}
}

func TestGetSwiftCodesByCountry(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	clearTable(db, t)

	record1 := model.SwiftCode{
		BankName:      "Bank One",
		Address:       "Address One",
		CountryISO2:   "AA",
		CountryName:   "CountryA",
		IsHeadquarter: true,
		SwiftCode:     "AA_SWIFTXXX1",
	}
	record2 := model.SwiftCode{
		BankName:      "Bank Two",
		Address:       "Address Two",
		CountryISO2:   "BB",
		CountryName:   "CountryB",
		IsHeadquarter: true,
		SwiftCode:     "BB_SWIFTXXX1",
	}
	record3 := model.SwiftCode{
		BankName:      "Bank Three",
		Address:       "Address Three",
		CountryISO2:   "AA",
		CountryName:   "CountryA",
		IsHeadquarter: false,
		SwiftCode:     "AA_SWIFTXXX2",
	}

	for _, rec := range []model.SwiftCode{record1, record2, record3} {
		if err := InsertSwiftCode(db, rec); err != nil {
			t.Fatalf("InsertSwiftCode nie powiodło się: %v", err)
		}
	}

	records, err := GetSwiftCodesByCountry(db, "aa")
	if err != nil {
		t.Fatalf("GetSwiftCodesByCountry nie powiodło się: %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Oczekiwano 2 rekordów dla kraju AA, otrzymano %d", len(records))
	}
}

func TestDeleteSwiftCode(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()
	clearTable(db, t)

	record := model.SwiftCode{
		BankName:      "Delete Bank",
		Address:       "Delete Address",
		CountryISO2:   "DL",
		CountryName:   "Deleteland",
		IsHeadquarter: true,
		SwiftCode:     "DELETESWIFTXXX",
	}
	if err := InsertSwiftCode(db, record); err != nil {
		t.Fatalf("InsertSwiftCode nie powiodło się: %v", err)
	}

	if err := DeleteSwiftCode(db, record.SwiftCode); err != nil {
		t.Fatalf("DeleteSwiftCode nie powiodło się: %v", err)
	}

	_, err := GetSwiftCode(db, record.SwiftCode)
	if err == nil {
		t.Error("Oczekiwano błędu przy pobieraniu usuniętego rekordu, ale błąd nie wystąpił")
	}
}
