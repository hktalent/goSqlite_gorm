package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var dir, esUrl *string

type fnCbk func(s string)

var step int64 = 0

func Log(msg string) {
	fmt.Printf("\r%8d  %s", step, msg)
}

func GetReq(id string, szLstMdf string) bool {
	// Post "77beaaf8081e4e45adb550194cc0f3a62ebb665f": unsupported protocol scheme ""
	req, err := http.NewRequest("GET", *esUrl+id+"/_source", nil)
	if err != nil {
		Log(fmt.Sprintf("%s error %v", id, err))
		return false
	}
	// 取消全局复用连接
	// tr := http.Transport{DisableKeepAlives: true}
	// client := http.Client{Transport: &tr}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	// keep-alive
	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer func() {
			err := resp.Body.Close() // resp 可能为 nil，不能读取 Body
			if nil != err {
				Log(fmt.Sprintf("%s error %v", id, err))
			}
		}()
	}
	if err != nil {
		Log(fmt.Sprintf("%s error %v", id, err))
		return false
	}
	var m map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if nil == err {
		json.Unmarshal(body, &m)
		if s1, ok := m["modified"]; ok {
			if s1 != szLstMdf {
				Log("will add " + id)
				return true
			}
		}
	}
	//log.Println(id, " 没有发生该改变")
	step++
	Log("")
	return false
}

var nThreads = make(chan struct{}, 1)

func sendReq(data []byte, id string, m map[string]interface{}) {
	if s1, ok := m["modified"]; ok {
		if !GetReq(id, s1.(string)) {
			return
		}
	}
	nThreads <- struct{}{}
	defer func() {
		<-nThreads
	}()
	fmt.Println("start send to ", *esUrl, " es "+id)
	// Post "77beaaf8081e4e45adb550194cc0f3a62ebb665f": unsupported protocol scheme ""
	req, err := http.NewRequest("POST", *esUrl+id, bytes.NewReader(data))
	if err != nil {
		Log(fmt.Sprintf("%s error %v", id, err))
		return
	}
	// 取消全局复用连接
	// tr := http.Transport{DisableKeepAlives: true}
	// client := http.Client{Transport: &tr}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.2 Safari/605.1.15")
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	// keep-alive
	req.Header.Add("Connection", "close")
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer func() {
			err := resp.Body.Close() // resp 可能为 nil，不能读取 Body
			if nil != err {
				Log(fmt.Sprintf("%s error %v", id, err))
			}
		}()
	}
	if err != nil {
		Log(fmt.Sprintf("%s error %v", id, err))
		return
	}

	// body, err := ioutil.ReadAll(resp.Body)
	// _, err = io.Copy(ioutil.Discard, resp.Body) // 手动丢弃读取完毕的数据
	// json.NewDecoder(resp.Body).Decode(&data)
	step++
	Log("")
	// req.Body.Close()
	// go http.Post(resUrl, "application/json",, post_body)
}

// dirents 返回 dir 目录中的条目
func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

func fnReadJson(s string) {
	s1, err := ioutil.ReadFile(s)
	if nil == err {
		var m map[string]interface{}
		err = json.Unmarshal(s1, &m)
		if nil == err {
			if id, ok := m["id"]; ok {
				sendReq(s1, id.(string), m)
			}
		}
	}
}

// wakjDir 递归地遍历以 dir 为根目录的整个文件树,并在 filesizes 上发送每个已找到的文件的大小
func walkDir(dir string, cbk fnCbk) {
	var subdir string
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir = filepath.Join(dir, entry.Name())
			walkDir(subdir, cbk)
		} else {
			subdir = filepath.Join(dir, entry.Name())
			if -1 < strings.Index(subdir, ".json") {
				cbk(subdir)
			}
		}
	}
}

func main() {
	dir = flag.String("dir", "", "json file dir")
	esUrl = flag.String("resUrl", "http://127.0.0.1:9200/intelligence_index/_doc/", "Elasticsearch url, eg: http://127.0.0.1:9200/dht_index/_doc/")
	flag.Parse()
	if "" == *esUrl || "" == *dir {
		return
	}
	walkDir(*dir, fnReadJson)
}
