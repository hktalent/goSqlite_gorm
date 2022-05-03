package common

import (
	"bytes"
	"fmt"
	mymod "goSqlite_gorm/pkg/models"
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

func GetMacWhereAmI(wami *mymod.WhereAmI) {
	p, err := os.Getwd()
	if nil != err {
		return
	}
	a, err := DoCmd("bash", "-c", "echo $PPSSWWDD|sudo -S "+p+"/tools/whereami")
	if nil != err {
		fmt.Println(err)
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

//
//func main() {
//	//	out, err := DoCmd("ls", "-l", "/var/log/*.log")
//	//	if err != nil {
//	//		fmt.Printf("combined out:\n%s\n", string(out))
//	//		log.Fatalf("cmd.Run() failed with %s\n", err)
//	//	}
//	wami := mymod.WhereAmI{}
//	GetMacWhereAmI(&wami)
//}
