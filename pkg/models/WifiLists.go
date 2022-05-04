package models

import (
	"gorm.io/gorm"
	"time"
)

// 当前位置 Wi-fi 列表
//  SSID BSSID             RSSI CHANNEL HT CC SECURITY (auth/unicast/group)
type WifiInfo struct {
	gorm.Model
	Wid      int64  `json:"wid" gorm:"foreignkey:Wid"`
	SSID     string `json:"ssid" jsonschema:"title=AP唯一的ID码,description=是你给自己的无线网络所取的名字"`
	BSSID    string `json:"bssid" gorm:"column:bssid;unique_index:bssid" jsonschema:"Basic Service Set 基本服务装置,description=6byte长度为48位bit的二进制标识符"`
	RSSI     string `json:"rssi" jsonschema:"Received Signal Strength Indicator是接收信号的强度指示"`
	CHANNEL  string `json:"channel"`
	HT       string `json:"ht"`
	CC       string `json:"cc"`
	SECURITY string `json:"security"` //  (auth/unicast/group)
}

type WifiInfoWifiListas struct {
	gorm.Model
	ID        int64      `json:"id" gorm:"primary_key"`
	Latitude  string     `json:"latitude" gorm:"column:latitude;unique_index:lat_alo"`
	Longitude string     `json:"longitude" gorm:"column:longitude;unique_index:lat_alo"`
	Accuracy  string     `json:"accuracy"`
	Date      time.Time  `json:"date"`
	WifiLists []WifiInfo `json:"wifiLists"`
}
