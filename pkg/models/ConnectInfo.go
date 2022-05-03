package models

import (
	"gorm.io/gorm"
)

// ip info
type IpInfo struct {
	Continent     string `json:"continent,omitempty"`
	ContinentCode string `json:"continentCode,omitempty"`
	Country       string `json:"country,omitempty"`
	CountryCode   string `json:"countryCode,omitempty"`
	Region        string `json:"region,omitempty"`
	RegionName    string `json:"regionName,omitempty"`
	City          string `json:"city,omitempty"`
	District      string `json:"district,omitempty"`
	Zip           string `json:"zip,omitempty"`
	Lat           string `json:"lat,omitempty"`
	Lon           string `json:"lon,omitempty"`
	Timezone      string `json:"timezone,omitempty"`
	Offset        string `json:"offset,omitempty"`
	Currency      string `json:"currency,omitempty"`
	Isp           string `json:"isp,omitempty"`
	Org           string `json:"org,omitempty"`
	As            string `json:"as,omitempty"`
	Asname        string `json:"asname,omitempty"`
	Mobile        string `json:"mobile,omitempty"`
	Proxy         string `json:"proxy,omitempty"`
	Hosting       string `json:"hosting,omitempty"`
	Query         string `json:"query" gorm:"unique_index,foreignKey:Query;references:Ip"` // IP
}

// 连接信息
type ConnectInfo struct {
	gorm.Model
	Pid    string  `json:"pid"`
	Ip     string  `json:"ip"`
	Cmd    string  `json:"cmd"`
	IpInfo *IpInfo `json:"ipInfo" gorm:"foreignKey:Query;references:Ip"`
}
