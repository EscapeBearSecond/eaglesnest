package common

// 执行方式
var (
	// 立即执行
	ExecuteImmediately = 1
	// 稍后执行
	ExecuteLater = 2
	// 定时执行
	ExecuteTiming = 3
)

// 模板类型
var (
	// 资产发现
	AssetDiscovery = "1"
	// 漏洞扫描
	VulnerabilityScan = "2"
	// 弱口令
	WeakPassword = "3"
)

var JobTypeName = map[string]string{
	AssetDiscovery:    "资产发现",
	VulnerabilityScan: "漏洞扫描",
	WeakPassword:      "弱口令",
}

// 普通任务状态
var (
	// 创建
	Created = 7
	// 执行中
	Running = 1
	// 成功
	Success = 2
	// 失败
	Failed = 3
	// 终止
	Stopped = 4
)

// 定时任务状态
var (
	// 运行
	TimeRunning = 5
	// 停止
	TimeStopped = 6
)
