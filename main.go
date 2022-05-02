package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	mycmd "goSqlite_gorm/pkg/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	SshServer = "http://127.0.0.1:2222"
	OkMsg     = "ok"
	ErrMsg    = "error"
	//OkCode  = 1
	ErrCode = -1
)

// https://github.com/swaggo/swag
// https://github.com/emicklei/go-restful
// github.com/alecthomas/jsonschema
// github.com/projectdiscovery/yamldoc-go
// gopkg.in/yaml.v3
// github.com/alecthomas/jsonschema

// how test:
// curl 'http://127.0.0.1:8081/api/v1/rsc/192.168.0.111/222'
type CommonDBInfo struct {
	//ID int `gorm:"column:id;primary_key;auto_increment;not null" json:"id"`

}
type WhereAmI struct {
	Latitude  string    `json:"latitude"`
	Longitude string    `json:"longitude"`
	Accuracy  string    `json:"accuracy"`
	Date      time.Time `json:"date"`
}

// 远程链接信息
type RemouteServerce struct {
	gorm.Model
	CommonDBInfo
	WhereAmI
	Title string `json:"title"`
	Ip    string `gorm:"column:ip;unique_index:ip_port" yaml:"ip,omitempty" json:"ip,omitempty"  jsonschema:"title=ip or domain Required parameters for connection,description=ip or domain Required parameters for connection"`
	Port  int    `gorm:"column:port;unique_index:ip_port" yaml:"port,omitempty" json:"port,omitempty" jsonschema:"title=remote port,description=ssh default 22"`
	User  string `gorm:"index"  yaml:"user,omitempty" json:"user,omitempty" jsonschema:"title=user name,description=user name"`
	P5wd  string `yaml:"p5wd,omitempty"  json:"p5wd,omitempty" jsonschema:"title=password,description=password"`
	Key   string `yaml:"key,omitempty" json:"key,omitempty" jsonschema:"title=ssh -i identity_file,description=Selects a file from which the identity (private key) for public key authentication is read.  You can also specify a public key file to use the corresponding
             private key that is loaded in ssh-agent(1) when the private key file is not present locally.  The default is ~/.ssh/id_rsa, ~/.ssh/id_ecdsa,
             ~/.ssh/id_ecdsa_sk, ~/.ssh/id_ed25519, ~/.ssh/id_ed25519_sk and ~/.ssh/id_dsa.  Identity files may also be specified on a per-host basis in the configuration
             file.  It is possible to have multiple -i options (and multiple identities specified in configuration files).  If no certificates have been explicitly
             specified by the CertificateFile directive, ssh will also try to load certificate information from the filename obtained by appending -cert.pub to identity
             filenames"`
	KeyP5wd string `yaml:"keyP5wd,omitempty"  json:"keyP5wd,omitempty" jsonschema:"title=key paswd,description=key paswd"`
	Type    string `yaml:"type,omitempty" json:"type,omitempty" jsonschema:"title=type:vnc ssh rdp,description=type:vnc ssh rdp"`
	Tags    string `gorm:"index" yaml:"tags,omitempty" json:"tags,omitempty" jsonschema:"title=tags hackerone butian,description=tags hackerone butian"` // 比较时hackerone，还是其他
	ImgData string `json:"imgData"`
}

// 组件信息
type ComponentInfo struct {
	gorm.Model
	CommonDBInfo
	Name    string   `yaml:"name,omitempty" json:"name,omitempty" jsonschema:"title=Component name,description=Component name"` // 组件name
	Info    string   `yaml:"info,omitempty" info:"port,omitempty" jsonschema:"title=Component info,description=Component info"`
	VuLists []string `yaml:"vulists,omitempty" json:"vulists,omitempty" jsonschema:"title=vul lists,description=vul lists"`
}

// 服务信息
type ServicesInfo struct {
	gorm.Model
	CommonDBInfo
	Ip            string `yaml:"ip,omitempty" json:"ip,omitempty" jsonschema:"title=ip or domain Required parameters for connection,description=ip or domain Required parameters for connection"`
	Port          int    `yaml:"port,omitempty" json:"port,omitempty" jsonschema:"title=connect to port,description=connect to port"`
	Info          string `yaml:"info,omitempty" info:"port,omitempty" jsonschema:"title=Component info,description=Component info"`
	ComponentInfo ComponentInfo
}

// 远程链接信息
type SiteInfo struct {
	gorm.Model
	CommonDBInfo
	Url                string         `yaml:"url,omitempty" json:"url,omitempty" jsonschema:"title=attack url,description=attack url"`
	ServsInfo          []ServicesInfo `yaml:"servsInfo,omitempty" json:"servsInfo,omitempty" jsonschema:"title=Services Info lists,description=Services Info lists"`
	Title              string         `yaml:"title,omitempty" json:"title,omitempty" jsonschema:"title=site title,description=site title"`
	ResponseServerName string         `yaml:"respServerName,omitempty" json:"respServerName,omitempty" jsonschema:"title=Response Server Name,description=Response Server Name"`
	ResponsePowerBy    string         `yaml:"respPowerBy,omitempty" json:"respPowerBy,omitempty" jsonschema:"title=Response Power By,description=Response Power By"`
	Tags               string         `yaml:"tags,omitempty" json:"tags,omitempty" jsonschema:"title=tags hackerone butian,description=tags hackerone butian"` // 比较时hackerone，还是其他
}

var dbCC *gorm.DB

func GetDb(dbName string, dst ...interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file:"+dbName+".db?cache=shared&mode=rwc&_journal_mode=WAL&Synchronous=Off&temp_store=memory&mmap_size=30000000000"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(dst[0])
	dbCC = db
	return db, nil
}

type ResultObj struct {
	Msg  string `json:msg`
	Code int    `json:code`
}

// @Summary      通过ip、port返回连接信息
// @Description  通过ip、port返回连接信息, curl 'http://127.0.0.1:8081/api/v1/rsc/192.168.0.111/222'
// @Tags         remoute,server,config
// @Accept       json
// @Produce      json
// @Param        ip path string true  "ip address"
// @Param        port path int true  "port"
// @Success 200 {object} RemouteServerce
// @Router       /api/v1/rsc/:ip/:port  [post]
func GetIPort(g *gin.Context) {
	var rsv RemouteServerce
	n, e := strconv.Atoi(g.Param("port"))
	if nil == e {
		rst := dbCC.First(&rsv, "ip = ? and port = ?", g.Param("ip"), n)
		//log.Println("query end", rst.RowsAffected)
		if 0 < rst.RowsAffected {
			g.JSON(http.StatusOK, rsv)
			return
		}
	}
	g.JSON(http.StatusBadRequest, gin.H{"msg": "not found", "code": -1})
}

type RmtSvIpName struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	ImgData string `json:"imgData"`
}

func GetRmtsvLists(g *gin.Context) {
	var aRst []RmtSvIpName
	// 查询时会自动选择 `id`, `title` 字段
	rst := dbCC.Model(&RemouteServerce{}).Limit(1000).Find(&aRst)
	// SELECT `id`, `name` FROM `users` LIMIT 10
	//log.Println("query end", rst.RowsAffected)
	if 0 < rst.RowsAffected && 0 < len(aRst) {
		g.JSON(http.StatusOK, aRst)
		return
	}
	g.JSON(http.StatusBadRequest, gin.H{"msg": "not found", "code": -1})
}

func ConnRmtSvs(g *gin.Context) *RemouteServerce {
	var rsv RemouteServerce
	//n, e := strconv.Atoi(strings.Split(g.Request.RequestURI, "/conn/")[1])
	n, e := strconv.Atoi(g.Param("id"))
	if nil == e {
		rst := dbCC.First(&rsv, "id = ?", n)
		if 0 < rst.RowsAffected {
			return &rsv
		}
	}
	return nil
}

type RmtSvImg struct {
	ID      uint   `json:"id"`
	ImgData string `json:"imgData"`
}

func SaveRmtsvImg(g *gin.Context) {
	var rsv RmtSvImg
	if err := g.BindJSON(&rsv); err == nil {
		xxxD := dbCC.Model(&RemouteServerce{})
		dbCC.Table("remoute_serverces").AutoMigrate(&RemouteServerce{})
		rst := xxxD.Where("id = ?", rsv.ID).Update("img_data", rsv.ImgData)
		//log.Println(rst.RowsAffected, rsv.ID, rst.Error)
		msg := OkMsg
		if nil != rst.Error {
			msg = fmt.Sprintf("%v", rst.Error)
		}
		g.JSON(http.StatusOK, gin.H{"msg": msg, "code": rst.RowsAffected})
		return
	}
}

// @Summary      保存ssh、vnc、rdp等远程连接配置信息
// @Description  保存连接信息,保存api,GORM V2 将使用 upsert 来保存关联记录
// @Tags         remoute,server,config
// @Accept       json
// @Produce      json
// @Router       /api/v1/rsc  [post]
func SaveRsCc(g *gin.Context) {
	var rsv, rsvOld RemouteServerce
	if err := g.BindJSON(&rsv); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"msg": err, "code": ErrCode})
		return
	}
	msg := OkMsg
	rst := dbCC.First(&rsvOld, "ip = ? and port = ?", rsv.Ip, rsv.Port)
	if 1 == rst.RowsAffected {
		rst = dbCC.Model(&rsv).Where("id = ?", rsvOld.ID).Updates(rsv)
	} else {
		rst = dbCC.Create(&rsv)
	}
	if nil != rst.Error {
		msg = fmt.Sprintf("%v", rst.Error)
	}
	g.JSON(http.StatusOK, gin.H{"msg": msg, "code": rst.RowsAffected})
}

// 反向代理封装
func DoReverseProxy(c *gin.Context, target string) {
	remote, err := url.Parse(target)
	if nil == err {
		director := func(req *http.Request) {
			req.Header = c.Request.Header
			req.URL.Scheme = remote.Scheme
			req.Host = remote.Host
			req.URL.Host = remote.Host
			req.RequestURI = c.Request.RequestURI
			req.URL.Path = c.Request.RequestURI
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func fnWriteHtml(c *gin.Context, szHtml string) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(szHtml))
}

// https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin
// curl 'http://127.0.0.1:8081/conn/0'
func ConnRmtSvsH(c *gin.Context) {
	var rsv *RemouteServerce
	rsv = ConnRmtSvs(c)
	if nil != rsv {
		log.Println(rsv.Tags, strings.Index(rsv.Tags, "vnc"))
		if -1 < strings.Index(rsv.Tags, "vnc") {
			mycmd.DoCmd("open", "vnc://"+rsv.User+"@"+rsv.Ip+":"+strconv.Itoa(rsv.Port))
			fnWriteHtml(c, "<i id=terminal-container>VNC is open from cmd shell</i>")
			return
		}
		if -1 < strings.Index(rsv.Tags, "rdp") {
			mycmd.DoCmd("rdesktop", rsv.Ip+":"+strconv.Itoa(rsv.Port))
			fnWriteHtml(c, "<i id=terminal-container>brew install rdesktop, RDP is open from cmd shell</i>")
			return
		}
		sss := make(url.Values)
		sss.Set("host", rsv.Ip)
		sss.Set("port", strconv.Itoa(rsv.Port))
		sss.Set("username", rsv.User)
		sss.Set("userpassword", rsv.P5wd)
		sss.Set("privatekey", rsv.Key)
		c.Request.Form = sss
		c.Request.PostForm = sss
		newBody := "host=" + rsv.Ip + "&port=" + strconv.Itoa(rsv.Port) + "&username=" + rsv.User + "&userpassword=" + url.QueryEscape(rsv.P5wd) + "&privatekey=" + url.QueryEscape(rsv.Key)
		c.Request.RequestURI = "/ssh/host/" // + newBody
		c.Request.URL.RawQuery = newBody
		// 1. set new header
		c.Request.Header.Set("Content-Length", strconv.Itoa(len(newBody)))

		// 2. also update this field
		c.Request.ContentLength = int64(len(newBody))
		c.Request.Method = "POST"
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded,charset=UTF-8")
		//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(newBody)))
		c.Request.Body = io.NopCloser(strings.NewReader(newBody))
		DoReverseProxy(c, SshServer)
		return
	}
}

// 反向代理分装，最后必须是类似：*id
func ReverseProxy(path, target string, router *gin.Engine) {
	xxx := func(c *gin.Context) {
		DoReverseProxy(c, target)
	}
	//router.Group(path, xxx)
	router.GET(path, xxx)
	router.POST(path, xxx)
}

// http://localhost:8080/swagger/index.html
// @title 51pwn app API
// @version 1.0
// @description This is 51pwn app api docs.
// @license.name Apache 2.0
// @contact.name go-swagger
// @contact.url https://github.com/hktalent/
// @host localhost:8080
// @BasePath /api/v1
func main() {
	db, err := GetDb("mydbfile", &RemouteServerce{})
	if err != nil {
		//panic("failed to connect database")
		log.Println(err)
		return
	}
	if nil != db {
		router := gin.Default()
		// 内部异常返回500
		router.Use(gin.Recovery())

		// docs.SwaggerInfo.BasePath = "/api/v1"
		// 同时运行多个gin服务并使用不同的swagger文档
		// https://xiaoliu.org/posts/2021/1230-gin-multi-swag/
		// https://golang.hotexamples.com/examples/github.com.gin-gonic.gin/RouterGroup/Group/golang-routergroup-group-method-examples.html
		v1 := router.Group("/api/v1")
		{
			// ssh、RDP、vnc 远程连接信息
			rscc := v1.Group("/rsc")
			rscc.POST("", SaveRsCc)
			rscc.GET("/:ip/:port", GetIPort)
			v1.GET("/rmtsvlists", GetRmtsvLists)
			v1.POST("/rmtsvImg", SaveRmtsvImg)
		}

		// ssh，必须在 connGrp.Use(ConnRmtSvsMiddleware()) 之后
		ReverseProxy("/ssh/*id", SshServer, router)
		router.GET("/conn/:id", ConnRmtSvsH)
		//router.Use(ConnRmtSvsMiddleware())

		// swagger 似乎成了所有例子的路径
		//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		router.Run(":8081")
	}
}
