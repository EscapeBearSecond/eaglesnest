package response

import "47.103.136.241/goprojects/curesan/server/model/curescan"

type TaskResult struct {
	PortScanList    []*curescan.PortScan
	OnlineCheckList []*curescan.OnlineCheck
}
