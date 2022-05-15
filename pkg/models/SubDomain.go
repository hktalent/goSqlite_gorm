package models

import (
	"gorm.io/gorm"
	"time"
)

/*
局部性更新文档，下面的代码借助go json的omitempty，在将更新数据对象序列化成json，
可以只序列化非零值字段，实现局部更新。 实际项目采用这种方式时，
需要注意某个字段的零值具有业务意义时，可以采用对应的指针类型实现
*/
type SubDomain struct {
	gorm.Model
	Domain    string    `json:"domain"`
	SubDomain string    `json:"subDomain"`
	ToolName  string    `json:"toolName,omitempty"`
	Create    time.Time `json:"create"`
}
