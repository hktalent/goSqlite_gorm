package models

import (
	"gorm.io/gorm"
	"time"
)

// 返回界面列表
type RmtSvIpName struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	ImgData string `json:"imgData"`
	Tags    string `gorm:"index" yaml:"tags,omitempty" json:"tags,omitempty" jsonschema:"title=tags hackerone butian,description=tags hackerone butian"` // 比较时hackerone，还是其他
}

// 位置坐标
type WhereAmI struct {
	Latitude  string    `json:"latitude"`
	Longitude string    `json:"longitude"`
	Accuracy  string    `json:"accuracy"`
	Date      time.Time `json:"date"`
}

// sudo arp-scan --localnet|grep '\d*\.\d*\.\d*\.\d*' | grep -v DUP
type Localnet struct {
	gorm.Model
	WhereAmI
	Ip   string `json:"ip"`
	Mac  string `json:"mac"`
	Name string `json:"name"`
}

// 图片更新
type RmtSvImg struct {
	ID      uint   `json:"id"`
	ImgData string `json:"imgData"`
	WhereAmI
}

// 远程链接信息
type RemouteServerce struct {
	gorm.Model
	WhereAmI
	Title string `json:"title"`
	Ip    string `gorm:"column:ip;unique_index:ip_port" yaml:"ip,omitempty" json:"ip,omitempty"  jsonschema:"title=ip or domain Required parameters for connection,description=ip or domain Required parameters for connection"`
	Port  int    `gorm:"column:port;unique_index:ip_port" yaml:"port,omitempty" json:"port,omitempty" jsonschema:"title=remote port,description=ssh default 22"`
	User  string `gorm:"index"  yaml:"user,omitempty" json:"user,omitempty" jsonschema:"title=user name,description=user name"`
	P5wd  string `yaml:"p5wd,omitempty"  json:"p5wd,omitempty" jsonschema:"title=password,description=password"`
	Key string `yaml:"key,omitempty" json:"key,omitempty" jsonschema:"title=ssh -i identity_file,description=Selects a file from which the identity (private key) for public key authentication is read.  You can also specify a public key file to use the corresponding
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
