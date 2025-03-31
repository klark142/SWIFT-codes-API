package parser

import (
	"io/ioutil"
	"os"
	"testing"
)

func createTempCSV(content string) (string, error) {
	tmpFile, err := ioutil.TempFile("", "test_*.csv")
	if err != nil {
		return "", err
	}
	_, err = tmpFile.Write([]byte(content))
	if err != nil {
		tmpFile.Close()
		return "", err
	}
	tmpFile.Close()
	return tmpFile.Name(), nil
}

func TestParseCSV_Valid(t *testing.T) {
	csvContent := `COUNTRY ISO2 CODE,SWIFT CODE,CODE TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
BG,ABIEBGS1XXX,BIC11,ABV INVESTMENTS LTD,"TSAR ASEN 20  VARNA, VARNA, 9002",VARNA,BULGARIA,Europe/Sofia
BG,ADCRBGS1XXX,BIC11,ADAMANT CAPITAL PARTNERS AD,"JAMES BOURCHIER BLVD 76A HILL TOWER SOFIA, SOFIA, 1421",SOFIA,BULGARIA,Europe/Sofia
UY,AFAAUYM1XXX,BIC11,AFINIDAD A.F.A.P.S.A.,"PLAZA INDEPENDENCIA 743  MONTEVIDEO, MONTEVIDEO, 11000",MONTEVIDEO,URUGUAY,America/Montevideo
`
	tmpFile, err := createTempCSV(csvContent)
	if err != nil {
		t.Fatalf("Nie udało się utworzyć tymczasowego pliku CSV: %v", err)
	}
	defer os.Remove(tmpFile)

	records, err := ParseCSV(tmpFile)
	if err != nil {
		t.Fatalf("Błąd parsowania CSV: %v", err)
	}

	if len(records) != 3 {
		t.Errorf("Oczekiwano 3 rekordów, otrzymano %d", len(records))
	}

	first := records[0]
	if first.CountryISO2 != "BG" {
		t.Errorf("Rekord 1 - Oczekiwano CountryISO2 'BG', otrzymano '%s'", first.CountryISO2)
	}
	if first.SwiftCode != "ABIEBGS1XXX" {
		t.Errorf("Rekord 1 - Oczekiwano SwiftCode 'ABIEBGS1XXX', otrzymano '%s'", first.SwiftCode)
	}
	if !first.IsHeadquarter {
		t.Errorf("Rekord 1 - Oczekiwano, że rekord jest główną siedzibą (isHeadquarter = true)")
	}

	second := records[1]
	if second.CountryISO2 != "BG" {
		t.Errorf("Rekord 2 - Oczekiwano CountryISO2 'BG', otrzymano '%s'", second.CountryISO2)
	}
	if second.SwiftCode != "ADCRBGS1XXX" {
		t.Errorf("Rekord 2 - Oczekiwano SwiftCode 'ADCRBGS1XXX', otrzymano '%s'", second.SwiftCode)
	}
	if !second.IsHeadquarter {
		t.Errorf("Rekord 2 - Oczekiwano, że rekord jest główną siedzibą (isHeadquarter = true)")
	}

	third := records[2]
	if third.CountryISO2 != "UY" {
		t.Errorf("Rekord 3 - Oczekiwano CountryISO2 'UY', otrzymano '%s'", third.CountryISO2)
	}
	if third.SwiftCode != "AFAAUYM1XXX" {
		t.Errorf("Rekord 3 - Oczekiwano SwiftCode 'AFAAUYM1XXX', otrzymano '%s'", third.SwiftCode)
	}
	if !third.IsHeadquarter {
		t.Errorf("Rekord 3 - Oczekiwano, że rekord jest główną siedzibą (isHeadquarter = true)")
	}
}


func TestParseCSV_InvalidFormat(t *testing.T) {

	csvContent := `COUNTRY ISO2 CODE,SWIFT CODE,CODE TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
BG,ABIEBGS1XXX,BIC11,ABV INVESTMENTS LTD,"TSAR ASEN 20  VARNA, VARNA",BULGARIA
`
	tmpFile, err := createTempCSV(csvContent)
	if err != nil {
		t.Fatalf("Nie udało się utworzyć tymczasowego pliku CSV: %v", err)
	}
	defer os.Remove(tmpFile)

	_, err = ParseCSV(tmpFile)
	if err == nil {
		t.Error("Oczekiwano błędu parsowania dla nieprawidłowego formatu, ale błąd nie został zgłoszony")
	}
}
