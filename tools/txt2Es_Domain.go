package main

import (
	"flag"
	sd "goSqlite_gorm/pkg/server"
	"io/ioutil"
)

func main() {
	var szFile *string
	szFile = flag.String("file", "", "json file dir")
	flag.Parse()
	s1, err := ioutil.ReadFile(*szFile)
	if nil == err {
		sd.DoListDomains(string(s1))
		select {}
	}
}
