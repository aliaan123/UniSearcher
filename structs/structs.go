package structs

// University struct that defines the fields used from the api university
type University struct {
	Name      string            `json:"name"`
	Country   string            `json:"country"`
	Isocode   string            `json:"isocode"`
	Webpages  []string          `json:"webpages"`
	Languages map[string]string `json:"languages"`
	Maps      Map               `json:"map"`
}

// UniFromHipo struct that defines fields of the respone given by the api.
type UniFromHipo struct {
	AlphaTwoCode  string   `json:"alpha_two_code"`
	Country       string   `json:"country"`
	StateProvince string   `json:"state_province"`
	Domains       []string `json:"domains"`
	Name          string   `json:"name"`
	WebPages      []string `json:"web_pages"`
}

// CombinedStruct used for combining the struct of the University and Country
type CombinedStruct struct {
	Name      string            `json:"name"`
	Country   string            `json:"country"`
	TwoCode   string            `json:"isocode"`
	WebPages  []string          `json:"webpages"`
	Languages map[string]string `json:"languages"`
	Map       string            `json:"maps"`
}

// Country struct that defines the fields for a country
type Country struct {
	Names     Name              `json:"name"`
	Languages map[string]string `json:"languages"`
	Maps      Map               `json:"maps"`
	Cca2      string            `json:"cca2"`
}

// Map struct that defines the fields for a map, is used in Country
type Map struct {
	GoogleMaps     string `json:"googleMaps"`
	OpenStreetMaps string `json:"openStreetMaps"`
}

// Name struct that defines the fields for a name, is used in Country
type Name struct {
	Common      string            `json:"common"`
	Official    string            `json:"official"`
	NativeNames map[string]string `json:"native_names"`
}

// Borders struct
type Borders struct {
	Borders []string `json:"borders"`
}

// DiagStruct that defines the fields for the diag endpoint
type DiagStruct struct {
	UniversitiesApi string  `json:"universitiesapi"`
	CountriesApi    string  `json:"countriesapi"`
	Version         string  `json:"version"`
	Uptime          float64 `json:"uptime"`
}
