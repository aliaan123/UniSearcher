package handler

// University struct that defines the fields used from the api university
type University struct {
	Name      string                 `json:"name"`
	Country   string                 `json:"country"`
	Isocode   string                 `json:"isocode"`
	Webpages  []string               `json:"webpages"`
	Languages map[string]interface{} `json:"languages"`
	Maps      []string               `json:"map"`
}

// University struct that defines fields of the respone given by the api.
type UniFromHipo struct {
	AlphaTwoCode  string   `json:"alpha_two_code"`
	Country       string   `json:"country"`
	StateProvince string   `json:"state_province"`
	Domains       []string `json:"domains"`
	Name          string   `json:"name"`
	WebPages      []string `json:"web_pages"`
}

type CombinedStruct struct {
	Name      string            `json:"name"`
	Country   string            `json:"country"`
	TwoCode   string            `json:"alpha_two_code"`
	Domains   []string          `json:"domains"`
	Languages map[string]string `json:"languages"`
	Map       string            `json:"maps"`
}

type Country struct {
	Names     Name              `json:"name"`
	Languages map[string]string `json:"languages"`
	Maps      Map               `json:"maps"`
	Cca2      string            `json:"cca2"`
}

type Map struct {
	GoogleMaps     string `json:"googleMaps"`
	OpenStreetMaps string `json:"openStreetMaps"`
}

type Name struct {
	Common      string            `json:"common"`
	Official    string            `json:"official"`
	NativeNames map[string]string `json:"native_names"`
}

type Borders struct {
	Borders []string `json:"borders"`
}

type DiagStruct struct {
	Universitiesapi string  `json:"universitiesapi"`
	Countriesapi    string  `json:"countriesapi"`
	Version         string  `json:"version"`
	Uptime          float64 `json:"uptime"`
}
