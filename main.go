package main

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	mycmd "goSqlite_gorm/pkg/common"
	"goSqlite_gorm/pkg/db"
	mymod "goSqlite_gorm/pkg/models"
	task "goSqlite_gorm/pkg/task"
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

// 组件信息
type ComponentInfo struct {
	gorm.Model
	Name    string   `yaml:"name,omitempty" json:"name,omitempty" jsonschema:"title=Component name,description=Component name"` // 组件name
	Info    string   `yaml:"info,omitempty" info:"port,omitempty" jsonschema:"title=Component info,description=Component info"`
	VuLists []string `yaml:"vulists,omitempty" json:"vulists,omitempty" jsonschema:"title=vul lists,description=vul lists"`
}

// 服务信息
type ServicesInfo struct {
	gorm.Model
	Ip            string `yaml:"ip,omitempty" json:"ip,omitempty" jsonschema:"title=ip or domain Required parameters for connection,description=ip or domain Required parameters for connection"`
	Port          int    `yaml:"port,omitempty" json:"port,omitempty" jsonschema:"title=connect to port,description=connect to port"`
	Info          string `yaml:"info,omitempty" info:"port,omitempty" jsonschema:"title=Component info,description=Component info"`
	ComponentInfo ComponentInfo
}

// 远程链接信息
type SiteInfo struct {
	gorm.Model
	Url                string         `yaml:"url,omitempty" json:"url,omitempty" jsonschema:"title=attack url,description=attack url"`
	ServsInfo          []ServicesInfo `yaml:"servsInfo,omitempty" json:"servsInfo,omitempty" jsonschema:"title=Services Info lists,description=Services Info lists"`
	Title              string         `yaml:"title,omitempty" json:"title,omitempty" jsonschema:"title=site title,description=site title"`
	ResponseServerName string         `yaml:"respServerName,omitempty" json:"respServerName,omitempty" jsonschema:"title=Response Server Name,description=Response Server Name"`
	ResponsePowerBy    string         `yaml:"respPowerBy,omitempty" json:"respPowerBy,omitempty" jsonschema:"title=Response Power By,description=Response Power By"`
	Tags               string         `yaml:"tags,omitempty" json:"tags,omitempty" jsonschema:"title=tags hackerone butian,description=tags hackerone butian"` // 比较时hackerone，还是其他
}

var dbCC *gorm.DB = db.GetDb("mydbfile", &mymod.RemouteServerce{})

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
	var rsv mymod.RemouteServerce
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

func GetId(g *gin.Context) {
	var rsv mymod.RemouteServerce
	n, e := strconv.Atoi(g.Param("id"))
	if nil == e {
		rst := dbCC.First(&rsv, "id = ?", n)
		//log.Println("query end", rst.RowsAffected)
		if 0 < rst.RowsAffected {
			g.JSON(http.StatusOK, rsv)
			return
		}
	}
	g.JSON(http.StatusBadRequest, gin.H{"msg": "not found", "code": -1})
}

// 通过泛型调用
func GetRmtsvLists(g *gin.Context) {
	db.GetRmtsvLists(g, mymod.RemouteServerce{}, []mymod.RmtSvIpName{})
}

// 当前连接信息
func GetccLists(g *gin.Context) {
	m, _ := time.ParseDuration("-3h")
	s0 := time.Now().Add(m)
	currentPage, err := strconv.Atoi(g.Request.FormValue("currentPage"))
	if nil != err {
		currentPage = 1
	}
	pageSize, err := strconv.Atoi(g.Request.FormValue("pageSize"))
	if nil != err {
		pageSize = 100
	}
	var rst []mymod.ConnectInfo = db.GetRmtsvLists4List(mymod.ConnectInfo{}, "IpInfo",
		[]mymod.ConnectInfo{}, pageSize, currentPage, "updated_at > ?", s0)
	if nil != rst && 0 < len(rst) {
		//	for i, x := range rst {
		//		if "" == x.IpInfo.Ip {
		//			xx1 := db.GetOne[mymod.IpInfo](&x.IpInfo, "ip=?", x.Ip)
		//			if nil != xx1 {
		//				x.IpInfo = *xx1
		//			}
		//			rst[i] = x
		//		}
		//	}
		m1 := make(map[string]interface{})
		m1["count"] = db.GetCount(mymod.ConnectInfo{}, "updated_at > ?", s0)
		m1["list"] = rst
		g.JSON(http.StatusOK, m1)
	}
	//db.GetRmtsvLists(g, mymod.ConnectInfo{}, []mymod.ConnectInfo{})
}

func ConnRmtSvs(g *gin.Context) *mymod.RemouteServerce {
	var rsv mymod.RemouteServerce
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

func SaveRmtsvImg(g *gin.Context) {
	var rsv mymod.RmtSvImg
	if err := g.BindJSON(&rsv); err == nil {
		mycmd.GetMacWhereAmI(&rsv.WhereAmI)
		xxxD := dbCC.Model(&mymod.RemouteServerce{})
		dbCC.Table("remoute_serverces").AutoMigrate(&mymod.RemouteServerce{})
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
	var rsv, rsvOld mymod.RemouteServerce
	if err := g.BindJSON(&rsv); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"msg": err, "code": ErrCode})
		return
	}
	msg := OkMsg
	mycmd.GetMacWhereAmI(&rsv.WhereAmI)
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
	var rsv *mymod.RemouteServerce
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
	go task.DoAllTask()
	if nil != dbCC {
		router := gin.Default()
		router.Use(static.Serve("/", static.LocalFile("dist", false)))
		router.NoRoute(func(c *gin.Context) {
			c.File("dist/index.html")
		})
		//router.Static("/", "./dist")
		//router.StaticFile("/index.html", "./dist/index.html")
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
			rscc.GET("/s/:id", GetId)

			v1.GET("/rmtsvlists", GetRmtsvLists)
			// curl 'http://127.0.0.1:8081/api/v1/cclsts'
			v1.GET("/cclsts", GetccLists)
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
