package model

type SwiftCode struct {
	BankName      string      `json:"bankName"`      
	Address       string      `json:"address"`       
	CountryISO2   string      `json:"countryISO2"`  
	CountryName   string      `json:"countryName"`   
	IsHeadquarter bool        `json:"isHeadquarter"`
	SwiftCode     string      `json:"swiftCode"`    
	Branches      []SwiftCode `json:"branches,omitempty"`
}
