package models

type Manufacturer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Country      string `json:"country"`
	FoundingYear int    `json:"foundingYear"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Specifications struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}

type Car struct {
	ID                    int            `json:"id"`
	Name                  string         `json:"name"`
	ManufacturerID        int            `json:"manufacturerId"`
	CategoryID            int            `json:"categoryId"`
	Category              Category       `json:"category"`
	Manufacturer          Manufacturer   `json:"manufacturer"`
	Year                  int            `json:"year"`
	Specifications        Specifications `json:"specifications"`
	Image                 string         `json:"image"`
	PreferredManufacturer string
}

type CarData struct {
	Manufacturers []Manufacturer `json:"manufacturers"`
	Categories    []Category     `json:"categories"`
	CarModels     []Car          `json:"carModels"`
}

type TemplateData struct {
	CarData               CarData
	ModelYears            []int
	PreferredManufacturer string
}