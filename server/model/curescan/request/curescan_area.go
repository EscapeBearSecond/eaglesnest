package request

type CreateArea struct {
	AreaName string   `json:"areaName"`
	AreaIP   []string `json:"areaIP"`
	AreaDesc string   `json:"areaDesc"`
}
