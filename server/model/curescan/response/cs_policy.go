package response

import "47.103.136.241/goprojects/curesan/server/global"

type PolicyDetail struct {
	global.GvaModel
	PolicyName     string         `json:"policyName"`
	PolicyDesc     string         `json:"policyDesc"`
	ScanType       []string       `json:"scanType"`
	PolicyConfig   []JobConfig    `json:"policyConfig"`
	OnlineCheck    bool           `json:"onlineCheck"`
	OnlineConfig   OnlineConfig   `json:"onlineConfig"`
	PortScan       bool           `json:"portScan"`
	PortScanConfig PortScanConfig `json:"portScanConfig"`
	Templates      []int64        `json:"templates"`
	IgnoredIP      []string       `json:"ignoredIp"`
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
	Templates   []int64 `json:"templates"`
}
