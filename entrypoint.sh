#!/bin/sh
set -e

echo "Czekam na uruchomienie bazy danych..."
sleep 10

echo "Sprawdzanie, czy tabela swift_codes istnieje..."
table_exists=$(psql "$DB_CONN" -t -c "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema='public' AND table_name='swift_codes');" | xargs)

if [ "$table_exists" = "f" ]; then
    echo "Tabela swift_codes nie istnieje. Tworzę schemat..."
    psql "$DB_CONN" -c "CREATE TABLE IF NOT EXISTS swift_codes (
        swift_code VARCHAR(20) PRIMARY KEY,
        bank_name TEXT NOT NULL,
        address TEXT NOT NULL,
        country_iso2 VARCHAR(2) NOT NULL,
        country_name TEXT NOT NULL,
        is_headquarter BOOLEAN NOT NULL
    ); CREATE INDEX IF NOT EXISTS idx_country_iso2 ON swift_codes(country_iso2);"
    echo "Schemat utworzony."
else
    echo "Tabela swift_codes już istnieje."
fi

echo "Sprawdzanie, czy tabela swift_codes jest pusta..."
row_count=$(psql "$DB_CONN" -t -c "SELECT COUNT(*) FROM swift_codes;" | xargs)

if [ "$row_count" -eq 0 ]; then
    echo "Baza jest pusta. Importowanie danych z CSV..."
    ./swift-codes-import --import
else
    echo "Baza już zawiera dane. Pominę import."
fi

echo "Uruchamianie serwera API..."
exec ./swift-codes
