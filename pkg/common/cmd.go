package common

import (
	"bytes"
	"encoding/json"
	"goSqlite_gorm/pkg/db"
	mymod "goSqlite_gorm/pkg/models"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/exec"
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
					rst := dbCC.Model(&mymod.IpInfo{}).Where("query=?", k.IpInfo).Find(&xx0)
					if 0 < rst.RowsAffected && nil != xx0 {
						k.IpInfo = xx0
					} else {
						k.IpInfo = GetIpInfo(k.Ip)
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
//	x := getCurConnInfo()
//	fmt.Printf("%v", x)
//}
