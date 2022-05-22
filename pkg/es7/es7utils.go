package es7

import (
	"bytes"
	"context"
	"encoding/json"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"log"
	"reflect"
	"strings"
)

type Es7Utils struct {
	Client *elasticsearch7.Client
}

var es7 *Es7Utils

func NewEs7() *Es7Utils {
	if nil != es7 {
		return es7
	}
	//time.Now()
	client, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: []string{"http://localhost:9200"},
		//Username:  "username",
		//Password:  "password",
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	es7 = &Es7Utils{Client: client}
	return es7
}

// get strutct name to index name
func (es7 *Es7Utils) GetIndexName(t1 any) string {
	return reflect.TypeOf(t1).Name()
}

// get Doc
func (es7 *Es7Utils) GetDoc(t1 any, id string) *esapi.Response {
	response, err := es7.Client.Get(es7.GetIndexName(t1), id)
	if nil != err {
		log.Println(err)
	}
	return response
}
func (es7 *Es7Utils) Update(t1 any, id string) *esapi.Response {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(&t1)
	if nil != err {
		return nil
	}
	response, err := es7.Client.Update(es7.GetIndexName(t1), id, body)
	if nil != err {
		log.Println(err)
	}
	return response
}

// 创建索引
func (es7 *Es7Utils) Create(t1 any, id string) string {
	body := &bytes.Buffer{}
	//pubDate := time.Now()
	err := json.NewEncoder(body).Encode(&t1)
	if nil != err {
		return ""
	}
	indexName := strings.ToLower(es7.GetIndexName(t1))
	// 覆盖性更新文档，如果给定的文档ID不存在，将创建文档: bytes.NewReader(data),
	response, err := es7.Client.Index(indexName, body, es7.Client.Index.WithDocumentID(id), es7.Client.Index.WithRefresh("true"))
	if nil == err {
		defer response.Body.Close()
		return response.String()
	}
	return ""
}

/*
{
	"_source":{
	  "excludes": ["author"]
	},
	"query": {
	  "match_phrase": {
		"author": "古龙"
	  }
	},
	"sort": [
	  {
		"pages": {
		  "order": "desc"
		}
	  }
	],
	"from": 0,
	"size": 5
}
*/
func (es7 *Es7Utils) Search(t1 any, query string) *esapi.Response {
	body := &bytes.Buffer{}
	body.WriteString(query)
	response, err := es7.Client.Search(es7.Client.Search.WithIndex(es7.GetIndexName(t1)), es7.Client.Search.WithBody(body))
	if nil == err {
		return response
	}
	return nil
}

// "select caseid,title from xc_cases where title like '%中国电信%'",
// 这里使用mysql的方式来请求，非常简单，符合开发习惯，简化es入门门槛，支持order，支持Limit，那么排序和分页就自己写好了
func (es7 *Es7Utils) QueryBySql(t1 any, query1 string) *esapi.Response {
	query := map[string]interface{}{
		"query": query1,
	}
	jsonBody, _ := json.Marshal(query)
	req := esapi.SQLQueryRequest{Body: bytes.NewReader(jsonBody)}
	res, _ := req.Do(context.Background(), es7.Client)
	return res
	// defer res.Body.Close()
}

//func main() {
//
//}
