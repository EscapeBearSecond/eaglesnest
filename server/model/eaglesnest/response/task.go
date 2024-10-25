package response

import "github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"

type TaskResult struct {
	PortScanList    []*eaglesnest.PortScan
	OnlineCheckList []*eaglesnest.OnlineCheck
	JobResultList   []*eaglesnest.JobResultItem
}

type Stage struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
	Total   int     `json:"total"`
	Running int     `json:"running"`
}
