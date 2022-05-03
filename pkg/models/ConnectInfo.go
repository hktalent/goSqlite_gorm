package models

import (
	"gorm.io/gorm"
)

// ip 经纬度 info
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

// domain info
type DomainInfo struct {
	gorm.Model
	Name string   `json:"name" gorm:"primaryKey,unique_index"`
	Ips  []IpInfo `json:"ips"`
}

// 连接信息
type ConnectInfo struct {
	gorm.Model
	Pid    string `json:"pid"`
	Ip     string `json:"ip"`
	Cmd    string `json:"cmd"`
	IpInfo IpInfo `json:"ipInfo" gorm:"foreignkey:Ip;references:Ip"`
}
