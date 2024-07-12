package request

import (
	"47.103.136.241/goprojects/curesan/server/model/common/request"
	"47.103.136.241/goprojects/curesan/server/model/curescan"
)

type CreatePolicy struct {
	PolicyName     string   `json:"policyName"`
	PolicyDesc     string   `json:"policyDesc"`
	ScanType       []string `json:"scanType"`
	PolicyConfig   string   `json:"policyConfig"`
	OnlineCheck    bool     `json:"onlineCheck"`
	OnlineConfig   string   `json:"onlineConfig"`
	PortScan       bool     `json:"portScan"`
	PortScanConfig string   `json:"portScanConfig"`
	Templates      []int64  `json:"templates"`
}

type UpdatePolicy struct {
	ID uint `json:"id"`
	CreatePolicy
}

type SearchPolicy struct {
	curescan.Policy
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
