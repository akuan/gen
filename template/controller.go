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


// @Summary 获取所有的{{pluralize .StructName}}
// @Tags {{.StructName}}
// @Accept  json
// @Produce  json
// @Param page query string false "第几页，>=1"
// @Param pagesize  query string false  "分页大小,默认10"
// @Param order query string false "排序列和排序方式，空格分隔,列: id desc"
// @Success 200 {object} model.JsonResult "{"code":0,"data":[model.{{.StructName}}],"msg":"ok","success":true}"
// @Success 500 {object} model.JsonResult "{"code":500,"data":{},"msg":"服务器错误","success":false}"
// @Router /api/{{.StructName | toLower}}  [GET]
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


// @Summary 根据ID获取单个{{.StructName}}
// @Tags {{.StructName}}
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} model.JsonResult "{"code":0,"data":model.{{.StructName}},"msg":"ok","success":true}"
// @Success 404 {object} model.JsonResult "{"code":404,"data":{},"msg":"{{.StructName}} with id 1 Not found","success":false}"
// @Router /api/{{.StructName | toLower}}/{id}  [GET]
func Get{{.StructName}}(c *gin.Context) {
	id := ParamInt(c, "id")
	{{.StructName | toLower}} := &model.{{.StructName}}{}
	if model.Db.First({{.StructName | toLower}}, id).Error != nil {
		NotFound(c, fmt.Sprintf("{{.StructName}} with id %v Not found ", id))
		return
	}
	JsonData(c, {{.StructName | toLower}}) 
}

// @Summary 新增{{.StructName}}
// @Tags {{.StructName}}
// @Accept  json
// @Produce  json
// @Param {{.StructName}} body model.{{.StructName}} true "新增{{.StructName}}"
// @Success 200 {object} model.JsonResult "{"code":0,"data":model.{{.StructName}},"msg":"ok","success":true}"
// @Success 500 {object} model.JsonResult "{"code":500,"data":{},"msg":"服务器错误","success":false}"
// @Router /api/{{.StructName | toLower}}   [POST]
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


// @Summary 更新{{.StructName}}
// @Tags {{.StructName}}
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Param {{.StructName}} body model.{{.StructName}} true "待更新的{{.StructName}}"
// @Success 200 {object} model.JsonResult "{"code":0,"data":model.{{.StructName}},"msg":"ok","success":true}"
// @Success 404 {object} model.JsonResult "{"code":404,"data":{},"msg":"{{.StructName}} with id 1 Not found","success":false}"
// @Success 500 {object} model.JsonResult "{"code":500,"data":{},"msg":"服务器错误","success":false}"
// @Router /api/{{.StructName | toLower}}/{id}  [PUT]
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


// @Summary 删除{{.StructName}}
// @Tags {{.StructName}}
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} model.JsonResult "{"code":0,"data":{},"msg":"ok","success":true}"
// @Success 404 {object} model.JsonResult "{"code":404,"data":{},"msg":"{{.StructName}} with id 1 Not found","success":false}"
// @Success 500 {object} model.JsonResult "{"code":500,"data":{},"msg":"服务器错误","success":false}"
// @Router /api/{{.StructName | toLower}}/{id}  [DELETE]
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
