package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"goSqlite_gorm/pkg/db"
	mymod "goSqlite_gorm/pkg/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// 最佳的方法是将命令写到临时文件，并通过bash进行执行
func DoCmd(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	// out, err := cmd.CombinedOutput()
	if nil != err {
		return "", err
	}
	return string(outStr + "\n" + errStr), err
}

var dbCC *gorm.DB = db.GetDb("mydbfile", &mymod.ConnectInfo{})

func GetCurConnInfo() []mymod.ConnectInfo {
	p, err := os.Getwd()
	if nil != err {
		log.Println(err)
		return nil
	}
	a, err := DoCmd("/bin/bash", p+"/tools/getCurNetConn.sh", "f")
	if nil != err {
		log.Println(err)
		return nil
	}
	x := strings.Split(a, "\n")
	clst := []mymod.ConnectInfo{}
	dbCC.AutoMigrate(&mymod.IpInfo{})
	for _, y := range x {
		k := mymod.ConnectInfo{}
		w := strings.Index(y, " ")
		if 1 < w {
			k.Pid = y[0:w]
			y = y[w+1:]
			w = strings.Index(y, " ")
			if 1 < w {
				k.Ip = y[0:w]
				y = y[w+1:]
				k.Cmd = y
				if -1 == strings.Index(k.Ip, "192.168.") && -1 == strings.Index(k.Ip, "172.16.") && -1 == strings.Index(k.Ip, "127.0.0") {
					var xx0 *mymod.IpInfo
					rst := dbCC.Model(&mymod.IpInfo{}).First(&xx0, "ip=?", k.Ip)
					if 0 < rst.RowsAffected && nil != xx0 {
						k.IpInfo = *xx0
					} else {
						var xx6 *mymod.IpInfo
						xx6 = GetIpInfo(k.Ip)
						if nil != xx6 {
							k.IpInfo = *xx6
						}
						//log.Println(k.IpInfo)
					}
				}
				clst = append(clst, k)
			}
		}
	}
	return clst
}

// 获取ip信息
func GetIpInfo(ip string) *mymod.IpInfo {
	req, err := http.NewRequest("GET", "http://ip-api.com/json/"+ip, nil)
	if err == nil {
		req.Header.Set("User-Agent", "curl/1.0")
		req.Header.Add("Cache-Control", "no-cache")
		// keep-alive
		req.Header.Add("Connection", "close")
		req.Close = true

		resp, err := http.DefaultClient.Do(req)
		if resp != nil {
			defer resp.Body.Close() // resp 可能为 nil，不能读取 Body
		}
		if err != nil {
			log.Println(err)
			return nil
		}
		var ipInfo *mymod.IpInfo
		err = json.NewDecoder(resp.Body).Decode(&ipInfo)
		if nil == err {
			return ipInfo
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	return nil
}

func DoWifiListsInfo() {
	var k *mymod.WifiLists = GetAirPortBSSID()
	dbCC.AutoMigrate(&mymod.WifiLists{})
	dbCC.AutoMigrate(&mymod.WifiInfo{})
	var x2 mymod.WifiLists
	xx1 := dbCC.Model(&mymod.WifiLists{}).Where("latitude=? and longitude=?", k.Latitude, k.Longitude)
	rst := xx1.Find(&x2)
	if 0 < rst.RowsAffected {
		rst = xx1.Updates(k)
		if nil != rst.Error {
			log.Println(rst.RowsAffected, rst.Error)
		}
	} else {
		rst = dbCC.Create(k)
	}
}

/*
SSID BSSID             RSSI CHANNEL HT CC SECURITY (auth/unicast/group)
*/
// /System/Library/PrivateFrameworks/Apple80211.framework/Versions/A/Resources/airport -s
func GetAirPortBSSID() *mymod.WifiLists {
	sCmd := "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/A/Resources/airport"
	if _, err := os.Stat(sCmd); errors.Is(err, os.ErrNotExist) {
		sCmd = "/usr/local/bin/airport"
	}
	a, err := DoCmd("/bin/bash", "-c", "echo $PPSSWWDD|sudo -S "+sCmd+" -s")
	if nil != err {
		log.Println(err)
		return nil
	}
	x := strings.Split(a, "\n")

	rstObj := mymod.WifiLists{}
	var wflst []mymod.WifiInfo = []mymod.WifiInfo{}
	for _, j := range x {
		r1, err := regexp.Compile(` ([0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}) `)
		if nil != err {
			continue
		}
		a1 := r1.FindAllString(j, -1)
		if 0 < len(a1) {
			wf1 := mymod.WifiInfo{}
			wf1.BSSID = strings.TrimSpace(a1[0])
			a1 = strings.Split(j, a1[0])
			wf1.SSID = strings.TrimSpace(a1[0])
			j = strings.TrimSpace(a1[1])
			a1 = strings.Split(j, "  ")
			wf1.RSSI = strings.TrimSpace(a1[0])
			j = j[len(a1[0])+2:]
			r1, err = regexp.Compile(` {2,}`)
			if nil == err {
				a1 = r1.Split(j, -1)
				wf1.CHANNEL = strings.TrimSpace(a1[0])
				wf1.HT = strings.TrimSpace(a1[1])
				j = strings.TrimSpace(a1[2])
				wf1.CC = j[0:2]
				wf1.SECURITY = strings.TrimSpace(j[2:])
			}
			wflst = append(wflst, wf1)
		}
	}
	var wam *mymod.WhereAmI = &mymod.WhereAmI{}
	GetMacWhereAmI(wam)
	s, err1 := json.Marshal(wam)
	if nil == err1 {
		json.Unmarshal(s, &rstObj)
	}
	rstObj.WifiInfos = wflst
	return &rstObj
}

// 获取当前坐标位置信息
func GetMacWhereAmI(wami *mymod.WhereAmI) {
	p, err := os.Getwd()
	if nil != err {
		return
	}
	a, err := DoCmd("/bin/bash", "-c", "echo $PPSSWWDD|sudo -S "+p+"/tools/whereami")
	if nil != err {
		log.Println(err)
		return
	}
	x := strings.Split(a, "\n")
	for _, j := range x {
		j = strings.TrimSpace(j)
		v := strings.Split(j, ": ")
		if 2 == len(v) {
			//fmt.Println(v[0], v[1])
			v[1] = strings.TrimSpace(v[1])
			switch v[0] {
			case "Latitude":
				wami.Latitude = v[1]
				break
			case "Longitude":
				wami.Longitude = v[1]
				break
			case "Accuracy (m)":
				wami.Accuracy = v[1]
				break
			case "Timestamp":
				wami.Date = time.Now()
				break
			}
		}
	}
	//fmt.Println(wami)
}

////
//func main() {
//	DoWifiListsInfo()
//}
