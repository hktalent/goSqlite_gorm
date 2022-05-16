package models

import (
	"gorm.io/gorm"
)

/*
局部性更新文档，下面的代码借助go json的omitempty，在将更新数据对象序列化成json，
可以只序列化非零值字段，实现局部更新。 实际项目采用这种方式时，
需要注意某个字段的零值具有业务意义时，可以采用对应的指针类型实现
*/
type SubDomain struct {
	gorm.Model
	Domain    string `json:"domain"`
	SubDomain string `json:"subDomain"`
	ToolName  string `json:"toolName,omitempty"`
}

// domain to ips
// ip to Domain
type Domain2Ips struct {
	gorm.Model
	Domain   string `json:"domain"`
	Ip       string `json:"ip"`
	ToolName string `json:"toolName,omitempty"`
}

// 端口扫描，及端口漏洞扫描
type Ip2Ports struct {
	gorm.Model
	MyId          string `json:"myId"`
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	Des           string `json:"des,omitempty"`
	ToolName      string `json:"toolName,omitempty"`
	VulsCheckFlag uint64 `json:"vulsCheckFlag,omitempty"` // 每一位表示一个工具，所以，可以支持64种工具、插件对该port进行扫描
	VulsCheckRst  string `json:"vulsCheckRst,omitempty"`
}

// 执行任务
type Task struct {
	gorm.Model
	Target   string `json:"target"`
	TaskType string `json:"taskType"`
	PluginId string `json:"pluginId"`
	Status   int    `json:"status"` // 状态:待执行，执行中，已完成
}

const (
	Task_Status_Pending     = 1 << iota // 待执行
	Task_Status_InExecution = 1 << iota // 执行中
	Task_Status_Completed   = 1 << iota // 已完成
)
const (
	Ip2Ports_VulsCheckFlag_Nmap     = 1 << iota // 端口扫描工具：Nmap
	Ip2Ports_VulsCheckFlag_Masscan  = 1 << iota // 端口扫描工具：masscan
	Ip2Ports_VulsCheckFlag_RustScan = 1 << iota // 端口扫描工具：RustScan
)
