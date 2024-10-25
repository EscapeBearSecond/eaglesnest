package response

type VulnTopN struct {
	Name     string `json:"name"`
	Count    int    `json:"count"`
	Severity string `json:"severity"`
}

type AssetTopN struct {
	Host     string `json:"host"`
	Count    int    `json:"count"`
	Critical int    `json:"critical"`
	High     int    `json:"high"`
	Medium   int    `json:"medium"`
	Low      int    `json:"low"`
}

type SeverityVuln struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Total    int `json:"total"`
}
