package request

import (
	"github.com/EscapeBearSecond/eaglesnest/server/model/common/request"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
)

type CreatePolicy struct {
	PolicyName     string          `json:"policyName"`
	PolicyDesc     string          `json:"policyDesc"`
	PolicyConfig   []*JobConfig    `json:"policyConfig"`
	OnlineConfig   *OnlineConfig   `json:"onlineConfig"`
	PortScanConfig *PortScanConfig `json:"portScanConfig"`
	IgnoredIP      []string        `json:"ignoredIp"`
}

type OnlineConfig struct {
	Use         bool   `json:"use"`
	Timeout     string `json:"timeout"`
	Count       int    `json:"count"`
	Format      string `json:"format"`
	RateLimit   int    `json:"rateLimit"`
	Concurrency int    `json:"concurrency"`
}

type PortScanConfig struct {
	Use         bool   `json:"use"`
	Timeout     string `json:"timeout"`
	Count       int    `json:"count"`
	Format      string `json:"format"`
	Ports       string `json:"ports"`
	RateLimit   int    `json:"rateLimit"`
	Concurrency int    `json:"concurrency"`
}

type JobConfig struct {
	Name        string  `json:"name"`
	Kind        string  `json:"kind"`
	Timeout     string  `json:"timeout"`
	Count       int     `json:"count"`
	Format      string  `json:"format"`
	RateLimit   int     `json:"rateLimit"`
	Concurrency int     `json:"concurrency"`
	IsAll       bool    `json:"isAll"`
	Templates   []int64 `json:"templates"`
}

type UpdatePolicy struct {
	ID uint `json:"ID"`
	CreatePolicy
}

type SearchPolicy struct {
	eaglesnest.Policy
	request.PageInfo
	OrderKey string `json:"orderKey"` // 排序
	Desc     bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
