package response

import "47.103.136.241/goprojects/curescan/server/model/curescan"

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
