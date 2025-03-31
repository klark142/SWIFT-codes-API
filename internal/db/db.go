package db

import (
	"database/sql"
	"fmt"
	"strings"
	
	"swift-codes/internal/model"

	_ "github.com/lib/pq" 
)

func InitDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("błąd przy łączeniu z bazą: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("błąd pingowania bazy: %w", err)
	}

	if err = createSchema(db); err != nil {
		return nil, fmt.Errorf("błąd przy tworzeniu schematu: %w", err)
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS swift_codes (
		swift_code VARCHAR(20) PRIMARY KEY,
		bank_name TEXT NOT NULL,
		address TEXT NOT NULL,
		country_iso2 VARCHAR(2) NOT NULL,
		country_name TEXT NOT NULL,
		is_headquarter BOOLEAN NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_country_iso2 ON swift_codes(country_iso2);
	`
	_, err := db.Exec(schema)
	return err
}

func GetSwiftCode(db *sql.DB, code string) (model.SwiftCode, error) {
	var sc model.SwiftCode
	query := `
		SELECT swift_code, bank_name, address, country_iso2, country_name, is_headquarter
		FROM swift_codes
		WHERE swift_code = $1
	`
	row := db.QueryRow(query, code)
	err := row.Scan(&sc.SwiftCode, &sc.BankName, &sc.Address, &sc.CountryISO2, &sc.CountryName, &sc.IsHeadquarter)
	if err != nil {
		return sc, err
	}
	return sc, nil
}


func GetBranchesByHeadquarter(db *sql.DB, headquarterCode string) ([]model.SwiftCode, error) {
	var branches []model.SwiftCode
	prefix := headquarterCode
	if len(headquarterCode) >= 8 {
		prefix = headquarterCode[:8]
	}

	query := `
		SELECT swift_code, bank_name, address, country_iso2, country_name, is_headquarter
		FROM swift_codes
		WHERE swift_code LIKE $1 AND is_headquarter = FALSE
	`
	rows, err := db.Query(query, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sc model.SwiftCode
		if err := rows.Scan(&sc.SwiftCode, &sc.BankName, &sc.Address, &sc.CountryISO2, &sc.CountryName, &sc.IsHeadquarter); err != nil {
			return nil, err
		}
		branches = append(branches, sc)
	}
	return branches, nil
}


func GetSwiftCodesByCountry(db *sql.DB, iso2 string) ([]model.SwiftCode, error) {
	var codes []model.SwiftCode
	query := `
		SELECT swift_code, bank_name, address, country_iso2, country_name, is_headquarter
		FROM swift_codes
		WHERE country_iso2 = $1
	`
	rows, err := db.Query(query, strings.ToUpper(iso2))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sc model.SwiftCode
		if err := rows.Scan(&sc.SwiftCode, &sc.BankName, &sc.Address, &sc.CountryISO2, &sc.CountryName, &sc.IsHeadquarter); err != nil {
			return nil, err
		}
		codes = append(codes, sc)
	}
	return codes, nil
}


func InsertSwiftCode(db *sql.DB, sc model.SwiftCode) error {
	query := `
		INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (swift_code) DO UPDATE 
		SET bank_name = EXCLUDED.bank_name,
		    address = EXCLUDED.address,
		    country_iso2 = EXCLUDED.country_iso2,
		    country_name = EXCLUDED.country_name,
		    is_headquarter = EXCLUDED.is_headquarter
	`
	_, err := db.Exec(query, sc.SwiftCode, sc.BankName, sc.Address, sc.CountryISO2, sc.CountryName, sc.IsHeadquarter)
	return err
}

func DeleteSwiftCode(db *sql.DB, code string) error {
	query := `DELETE FROM swift_codes WHERE swift_code = $1`
	_, err := db.Exec(query, code)
	return err
}


