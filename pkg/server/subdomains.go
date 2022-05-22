package server

import (
	"github.com/gin-gonic/gin"
	"github.com/hktalent/go4Hacker/lib/hacker"
	kv "goSqlite_gorm/pkg/common"
	"goSqlite_gorm/pkg/es7"
	"log"
	"net/http"
	"strings"
)

type SubDomain struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
}

type Domain struct {
	Domain string   `json:"domain"`
	Ips    []string `json:"ips"`
}

func SaveDomain(domain string, ips []string) {
	var d = Domain{Domain: domain, Ips: ips}
	s := es7.NewEs7().Create(d, domain)
	log.Println(s)
}

var cache = kv.NewKvDbOp()

func GetIps(domain string) []string {
	a := kv.GetAny[[]string](cache, domain)
	if nil != a {
		return a
	}
	a1 := hacker.GetDomian2IpsAll(domain)
	if nil != a && 0 < len(a) {
		go kv.PutAny[[]string](cache, domain, a1)
	}
	return a1
}

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
	if nil != m.Subdomains {
		for _, x := range m.Subdomains {
			if -1 < strings.Index(x, ":") {
				x = strings.Split(x, ":")[0]
			}
			a = hacker.GetDomian2IpsAll(x)
			if 0 < len(a) {
				go SaveDomain(x, a)
			}
		}
	}
	g.JSON(http.StatusOK, gin.H{"msg": "ok", "code": 200})
}
func InitSubDomainRoute(router *gin.RouterGroup) {
	router.POST("/subdomian", SaveSubDomain)
}
