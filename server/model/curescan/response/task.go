package response

import "github.com/EscapeBearSecond/curescan/server/model/curescan"

type TaskResult struct {
	PortScanList    []*curescan.PortScan
	OnlineCheckList []*curescan.OnlineCheck
	JobResultList   []*curescan.JobResultItem
}

type Stage struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
	Total   int     `json:"total"`
	Running int     `json:"running"`
}
