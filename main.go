package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

// https://github.com/swaggo/swag
// https://github.com/emicklei/go-restful
// github.com/alecthomas/jsonschema
// github.com/projectdiscovery/yamldoc-go
// gopkg.in/yaml.v3
// github.com/alecthomas/jsonschema

// 远程链接信息
type RemouteServerce struct {
	gorm.Model
	Title string `json:"title"`
	Ip    string `gorm:"primary_key" yaml:"ip,omitempty" json:"ip,omitempty"  jsonschema:"title=ip or domain Required parameters for connection,description=ip or domain Required parameters for connection"`
	Port  int    `gorm:"primary_key" yaml:"port,omitempty" json:"port,omitempty" jsonschema:"title=remote port,description=ssh default 22"`
	User  string `gorm:"index"  yaml:"user,omitempty" json:"user,omitempty" jsonschema:"title=user name,description=user name"`
	P5wd  string `yaml:"p5wd,omitempty"  json:"p5wd,omitempty" jsonschema:"title=password,description=password"`
	Key   string `yaml:"key,omitempty" json:"key,omitempty" jsonschema:"title=ssh -i identity_file,description=Selects a file from which the identity (private key) for public key authentication is read.  You can also specify a public key file to use the corresponding
             private key that is loaded in ssh-agent(1) when the private key file is not present locally.  The default is ~/.ssh/id_rsa, ~/.ssh/id_ecdsa,
             ~/.ssh/id_ecdsa_sk, ~/.ssh/id_ed25519, ~/.ssh/id_ed25519_sk and ~/.ssh/id_dsa.  Identity files may also be specified on a per-host basis in the configuration
             file.  It is possible to have multiple -i options (and multiple identities specified in configuration files).  If no certificates have been explicitly
             specified by the CertificateFile directive, ssh will also try to load certificate information from the filename obtained by appending -cert.pub to identity
             filenames"`
	KeyP5wd   string `yaml:"keyP5wd,omitempty"  json:"keyP5wd,omitempty" jsonschema:"title=key paswd,description=key paswd"`
	Type      string `yaml:"type,omitempty" json:"type,omitempty" jsonschema:"title=type:vnc ssh rdp,description=type:vnc ssh rdp"`
	Tags      string `gorm:"index" yaml:"tags,omitempty" json:"tags,omitempty" jsonschema:"title=tags hackerone butian,description=tags hackerone butian"` // 比较时hackerone，还是其他
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
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
		log.Println("query end", rst.RowsAffected)
		if 0 < rst.RowsAffected {
			g.JSON(http.StatusOK, rsv)
			return
		}
	}
	g.JSON(http.StatusOK, gin.H{"msg": "not found", "code": -1})
}

// @Summary      保存ssh、vnc、rdp等远程连接配置信息
// @Description  保存连接信息,保存api,GORM V2 将使用 upsert 来保存关联记录
// @Tags         remoute,server,config
// @Accept       json
// @Produce      json
// @Router       /api/v1/rsc  [post]
func SaveRsCc(g *gin.Context) {
	var rsv RemouteServerce
	if err := g.BindJSON(&rsv); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"msg": err, "code": -1})
		return
	}
	log.Println(rsv)
	rst := dbCC.Create(&rsv)
	msg := "error"
	if 1 == rst.RowsAffected {
		msg = "ok"
	}

	g.JSON(http.StatusOK, gin.H{"msg": msg, "code": rst.RowsAffected})
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
		router.Use(gin.Recovery())
		//docs.SwaggerInfo.BasePath = "/api/v1"
		// 同时运行多个gin服务并使用不同的swagger文档
		// https://xiaoliu.org/posts/2021/1230-gin-multi-swag/
		// https://golang.hotexamples.com/examples/github.com.gin-gonic.gin/RouterGroup/Group/golang-routergroup-group-method-examples.html
		v1 := router.Group("/api/v1")
		{
			rscc := v1.Group("/rsc")
			rscc.POST("", SaveRsCc)
			rscc.GET("/:ip/:port", GetIPort)

			//route.DELETE("/:id/likings/:userId", userPermission.AuthRequired(deleteLikingOnUser))
			//route.GET("/:id/liked", retrieveLikedOnUsers)
			// SiteInfo
			//eg := v1.Group("/si")
			//{
			//	eg.GET("/helloworld", Helloworld)
			//}
		}
		// swagger 似乎成了所有例子的路径
		//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		router.Run(":8081")
		// Create
		//db.Create(&Product{Code: "D42", Price: 100})

		// Read
		//var product RemouteServerce
		////db.First(&product, 1)                 // find product with integer primary key
		//db.First(&product, "code = ?", "D42") // find product with code D42
		//
		//// Update - update product's price to 200
		//db.Model(&product).Update("Price", 200)
		//// Update - update multiple fields
		////db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
		//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

		// Delete - delete product
		//db.Delete(&product, 1)
	}
}
