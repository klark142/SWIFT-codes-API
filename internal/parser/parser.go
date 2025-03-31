package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"swift-codes/internal/model"
)

func ParseCSV(filePath string) ([]model.SwiftCode, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("nie udało się otworzyć pliku: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("błąd odczytu CSV: %w", err)
	}

	var swiftCodes []model.SwiftCode

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 7 {
			return nil, fmt.Errorf("nieprawidłowy format w wierszu %d", i+1)
		}

		countryISO2 := strings.ToUpper(strings.TrimSpace(record[0]))
		swiftCode := strings.TrimSpace(record[1])
		bankName := strings.ToUpper(strings.TrimSpace(record[3]))
		address := strings.TrimSpace(record[4])
		countryName := strings.ToUpper(strings.TrimSpace(record[6]))

		isHeadquarter := strings.HasSuffix(swiftCode, "XXX")

		sc := model.SwiftCode{
			BankName:      bankName,
			Address:       address,
			CountryISO2:   countryISO2,
			CountryName:   countryName,
			IsHeadquarter: isHeadquarter,
			SwiftCode:     swiftCode,
		}

		swiftCodes = append(swiftCodes, sc)
	}

	headquarterMap := make(map[string]*model.SwiftCode)
	var result []model.SwiftCode

	for i, sc := range swiftCodes {
		if sc.IsHeadquarter {
			prefix := sc.SwiftCode
			if len(sc.SwiftCode) >= 8 {
				prefix = sc.SwiftCode[:8]
			}
			headquarterMap[prefix] = &swiftCodes[i]
			result = append(result, swiftCodes[i])
		}
	}

	for _, sc := range swiftCodes {
		if !sc.IsHeadquarter {
			prefix := sc.SwiftCode
			if len(sc.SwiftCode) >= 8 {
				prefix = sc.SwiftCode[:8]
			}
			if hq, exists := headquarterMap[prefix]; exists {
				hq.Branches = append(hq.Branches, sc)
			} else {
				result = append(result, sc)
			}
		}
	}

	return result, nil
}
