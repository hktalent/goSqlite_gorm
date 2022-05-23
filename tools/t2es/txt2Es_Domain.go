package main

import (
	"flag"
	sd "goSqlite_gorm/pkg/server"
	"io/ioutil"
)

func main() {
	var szFile *string
	var d2i *string
	szFile = flag.String("file", "", "json file dir")
	d2i = flag.String("d2i", "", "domain 2 ip")
	flag.Parse()
	//xx1 := "/Users/51pwn/MyWork/mybugbounty/dlst1.txt"
	//*szFile = xx1
	//*d2i = "xx"
	s1, err := ioutil.ReadFile(*szFile)
	if nil == err {
		if "" != *d2i {
			go sd.DoDomain2Ips(string(s1))
		} else {
			sd.DoListDomains(string(s1))
		}
		select {}
	}
}
