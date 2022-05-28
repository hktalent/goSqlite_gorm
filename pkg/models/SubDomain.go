package models

import (
	"gorm.io/gorm"
)

// 存储到ES
type SubDomain struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
	Tags       string   `json:"tags,omitempty"` // 标识属于那个tag，例如hackerone
}

//http://127.0.0.1:9200/domain_index/_search?q=domain:%20in%20*qianxin*
type Domain struct {
	Domain string   `json:"domain"`
	Ips    []string `json:"ips"`
}

/*
局部性更新文档，下面的代码借助go json的omitempty，在将更新数据对象序列化成json，
可以只序列化非零值字段，实现局部更新。 实际项目采用这种方式时，
需要注意某个字段的零值具有业务意义时，可以采用对应的指针类型实现
*/
type SubDomainItem struct {
	gorm.Model
	Domain    string `json:"domain"`
	SubDomain string `json:"subDomain"`
	ToolName  uint64 `json:"toolName,omitempty"` // 支持多个工具
	Tags      string `json:"tags,omitempty"`     // 标识属于那个tag，例如hackerone
}

// domain to ips
// ip to Domain
type Domain2Ips struct {
	gorm.Model
	Domain   string `json:"domain"`
	Ip       string `json:"ip"`
	ToolName uint64 `json:"toolName,omitempty"` // 支持多个工具
}

// 端口扫描，及端口漏洞扫描
type Ip2Ports struct {
	gorm.Model
	MyId          string `json:"myId,omitempty"` // 对应ES domain id，可以为空
	Ip            string `json:"ip"`
	Port          int    `json:"port"`
	Des           string `json:"des,omitempty"`
	ToolName      uint64 `json:"toolName,omitempty"`      // 支持多个工具
	VulsCheckFlag uint64 `json:"vulsCheckFlag,omitempty"` // 每一位表示一个工具，所以，可以支持64种工具、插件对该port进行扫描
	VulsCheckRst  string `json:"vulsCheckRst,omitempty"`
}

// ip 经纬度 info
// curl -H 'User-Agent:curl/1.0' http://ip-api.com/json/107.182.191.202|jq
type IpInfo struct {
	gorm.Model
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	Timezone      string  `json:"timezone"`
	Offset        string  `json:"offset"`
	Currency      string  `json:"currency"`
	Isp           string  `json:"isp"`
	Org           string  `json:"org"`
	As            string  `json:"as"`
	Asname        string  `json:"asname"`
	Mobile        string  `json:"mobile"`
	Proxy         string  `json:"proxy"`
	Hosting       string  `json:"hosting"`
	Ip            string  `json:"query" gorm:"primaryKey,unique_index"` // IP
}

// 执行任务
type Task struct {
	gorm.Model
	Target   string `json:"target"`
	TaskType uint64 `json:"taskType"`
	PluginId string `json:"pluginId"`
	Status   int    `json:"status"` // 状态:待执行，执行中，已完成
}

const (
	// TaskType
	TaskType_Subdomain = 1 << iota // 任务类型：子域名
	TaskType_IP2Port   = 1 << iota // 任务类型：端口扫描
	TaskType_UrlScan   = 1 << iota // 任务类型：url扫描
	TaskType_VulsScan  = 1 << iota // 任务类型：漏洞扫描

	// Status
	Task_Status_Pending     = 1 << iota // 任务状态：待执行
	Task_Status_InExecution = 1 << iota // 任务状态：执行中
	Task_Status_Completed   = 1 << iota // 任务状态：已完成

	// PluginId
	SubDomains_Sublist3r            = 1 << iota // 子域名：Sublist3r
	SubDomains_Ksubdomain           = 1 << iota // 子域名：ksubdomain
	Ip2Ports_VulsCheckFlag_Nmap     = 1 << iota // 端口扫描工具：Nmap
	Ip2Ports_VulsCheckFlag_Masscan  = 1 << iota // 端口扫描工具：masscan
	Ip2Ports_VulsCheckFlag_RustScan = 1 << iota // 端口扫描工具：RustScan
)
