package response

type VulnTopN struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type AssetTopN struct {
	Host  string `json:"host"`
	Count int    `json:"count"`
}
