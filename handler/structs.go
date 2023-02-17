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

/*
type UniversitiesResponse struct {
	//Universities []University `json:"data"`
	Name       string       `json:"name"`
	University []University `json:"data"`
}
*/
