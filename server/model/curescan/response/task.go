package response

import "codeup.aliyun.com/66d825f8c06a2fdac7bbfe8c/curescan/server/model/curescan"

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
