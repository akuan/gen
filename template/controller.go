package template

var ControllerTmpl = `package controller

import (
	"fmt"
 
	"{{.PackageName}}"
	"github.com/gin-gonic/gin"
)

func config{{pluralize .StructName}}Router(router *gin.RouterGroup) {
	router.GET("/{{.StructName | toLower}}", GetAll{{pluralize .StructName}})
	router.POST("/{{.StructName | toLower}}", Add{{.StructName}})
	router.GET("/{{.StructName | toLower}}/:id", Get{{.StructName}})
	router.PUT("/{{.StructName | toLower}}/:id", Update{{.StructName}})
	router.DELETE("/{{.StructName | toLower}}/:id", Delete{{.StructName}})
}

func GetAll{{pluralize .StructName}}(c *gin.Context) {
	page := QueryInt(c, "page")
	if page < 1 {
		page = 1
	}
	pagesize := QueryInt(c, "pagesize")
	if pagesize <= 0 {
		pagesize = 10
	}
	offset := (page - 1) * pagesize
	order := c.Query("order") 
	{{pluralize .StructName | toLower}} := []*model.{{.StructName}}{}	
	var err error
	if order != "" {
		err = model.Db.Model(&model.{{.StructName}}{}).Order(order).Offset(offset).Limit(pagesize).Find(&{{pluralize .StructName | toLower}}).Error
	} else {
		err = model.Db.Model(&model.{{.StructName}}{}).Offset(offset).Limit(pagesize).Find(&{{pluralize .StructName | toLower}}).Error
	}

	if err != nil {
		ServerError(c, err.Error())
		return
	}
	JsonData(c, {{pluralize .StructName | toLower}})
}

func Get{{.StructName}}(c *gin.Context) {
	id := ParamInt(c, "id")
	{{.StructName | toLower}} := &model.{{.StructName}}{}
	if model.Db.First({{.StructName | toLower}}, id).Error != nil {
		NotFound(c, fmt.Sprintf("{{.StructName}} with id %v Not found ", id))
		return
	}
	JsonData(c, {{.StructName | toLower}}) 
}

func Add{{.StructName}}(c *gin.Context) {
	{{.StructName | toLower}} := &model.{{.StructName}}{}
   if err := c.ShouldBindJSON({{.StructName | toLower}}); err != nil {
		ServerError(c, err.Error())
		return
	}
	if err := model.Db.Save({{.StructName | toLower}}).Error; err != nil {
		ServerError(c, err.Error())
		return
	}
	JsonData(c, {{.StructName | toLower}}) 
}

func Update{{.StructName}}(c *gin.Context) {	 
    id := ParamInt(c, "id")
		
	{{.StructName | toLower}} := &model.{{.StructName}}{}
	if model.Db.First({{.StructName | toLower}}, id).Error != nil {
		NotFound(c, fmt.Sprintf(" update Error {{.StructName | toLower}} with id %v not Found", id))
		return
	}

	updated := &model.{{.StructName}}{}
	if err := c.ShouldBindJSON(updated); err != nil {
		ServerError(c, err.Error())
		return
	}

	if err := model.Copy({{.StructName | toLower}}, updated); err != nil{
		ServerError(c, err.Error())
		return
	}

	if err := model.Db.Save({{.StructName | toLower}}).Error; err != nil {
		ServerError(c, err.Error())
		return
	}
	JsonData(c, {{.StructName | toLower}})
}

func Delete{{.StructName}}(c *gin.Context) {
	id := ParamInt(c, "id")
	{{.StructName | toLower}} := &model.{{.StructName}}{}

	if model.Db.First({{.StructName | toLower}}, id).Error != nil {
		NotFound(c, fmt.Sprintf(" delete Error {{.StructName | toLower}} with id %v not Found", id))
		return
	}
	if err := model.Db.Delete({{.StructName | toLower}}).Error; err != nil {
		ServerError(c, err.Error())
		return
	}
	JsonData(c, "")
}
`
