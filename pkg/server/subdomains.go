package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hktalent/go4Hacker/lib/hacker"
	kv "goSqlite_gorm/pkg/common"
	"goSqlite_gorm/pkg/es7"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type SubDomain struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
}

//http://127.0.0.1:9200/domain_index/_search?q=domain:%20in%20*qianxin*
type Domain struct {
	Domain string   `json:"domain"`
	Ips    []string `json:"ips"`
}

func SaveDomain(domain string, ips []string) {
	var d = Domain{Domain: domain, Ips: ips}
	x1 := es7.NewEs7()
	x2 := x1.GetDoc(d, domain)
	if nil != x2 {
		if -1 < strings.Index(x2.String(), domain) {
			return
		}
	}
	s := x1.Create(d, domain)
	log.Println(s)
}

// var cache = kv.NewKvDbOp()

var nGetIp = make(chan struct{}, 12800)

func GetIps(domain string) []string {
	defer func() {
		<-nGetIp
	}()
	nGetIp <- struct{}{}
	a, err := kv.GetAny[[]string](domain)
	if nil == err {
		return a
	}
	a1 := hacker.GetDomian2IpsAll(domain)
	if nil != a && 0 < len(a) {
		log.Println("ok ", domain)
		go kv.PutAny[[]string](domain, a1)
	}
	return a1
}

func DoDomainLists(a []string) {
	if nil != a {
		for _, x := range a {
			// 跳过ip
			xreg, err := regexp.Compile(`(\d{1,3}\.){3}\d{1,3}`)
			if nil == err {
				x11 := xreg.FindAllString(x, -1)
				if nil != x11 && 0 < len(x11) {
					continue
				}
			}
			if -1 < strings.Index(x, "://") {
				x = strings.Split(x, "://")[1]
			}
			if -1 < strings.Index(x, ":") {
				x = strings.Split(x, ":")[0]
			}
			a = GetIps(x)
			if 0 < len(a) {
				go SaveDomain(x, a)
			}
		}
	}
}

func DoDomain2Ips(s string) {
	if "" != s {
		xreg, err := regexp.Compile(`[\n;]`)
		if nil == err {
			a := xreg.Split(s, -1)
			for i, x := range a {
				if -1 < strings.Index(x, "//") {
					a[i] = strings.Split(x, "//")[1]
				}
				if -1 < strings.Index(a[i], ":") {
					a[i] = strings.Split(a[i], "//")[0]
				}
				a[i] = strings.Replace(a[i], "*", "", -1)
				go (func(s00 string) []string {
					return GetIps(s00)
				})(a[i])
			}

		}
	}
}

func DoListDomains(s string) {
	if "" != s {
		xreg, err := regexp.Compile(`[\n;]`)
		if nil == err {
			a := xreg.Split(s, -1)
			for i, x := range a {
				xreg, err = regexp.Compile(`(\d{1,3}\.){3}\d{1,3}`)
				if nil == err {
					x11 := xreg.FindAllString(s, -1)
					if nil != x11 && 0 < len(x11) {
						continue
					}
				}

				if -1 < strings.Index(x, "//") {
					a[i] = strings.Split(x, "//")[1]
				}
				if -1 < strings.Index(a[i], ":") {
					a[i] = strings.Split(a[i], "//")[0]
				}
				a[i] = strings.Replace(a[i], "*", "", -1)
			}
			go DoDomainLists(a)
		}
	}
}

// dlst=
func SaveDomainLst(g *gin.Context) {
	s := g.Request.FormValue("dlst")
	DoListDomains(s)
	g.JSON(http.StatusOK, gin.H{"msg": "ok", "code": 200})
}

/*
{
"domain": "xx.com",
"subdomains": []
}

*/
func SaveSubDomain(g *gin.Context) {
	var m SubDomain
	if err := g.BindJSON(&m); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"msg": err, "code": ErrCode})
		return
	}

	if -1 < strings.Index(m.Domain, "//") {
		m.Domain = strings.Split(m.Domain, "//")[1]
	}
	a := GetIps(m.Domain)
	if 0 < len(a) {
		go SaveDomain(m.Domain, a)
	}
	go DoDomainLists(m.Subdomains)
	g.JSON(http.StatusOK, gin.H{"msg": "ok", "code": 200})
}
func InitSubDomainRoute(router *gin.RouterGroup) {
	router.POST("/subdomian", SaveSubDomain)
	router.POST("/dlists", SaveDomainLst)
}
