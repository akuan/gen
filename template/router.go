package template

var RouterTmpl = `package controller

import (
	"fmt"
	"labgo/modules/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


func ConfigRouter( router *gin.RouterGroup)   {
    {{range .}}config{{pluralize .}}Router(router)
    {{end}}
}

//从查询字符串中获取Int值。
func QueryInt(c *gin.Context, key string) int {
	si := c.Query(key)
	i, err := strconv.Atoi(si)
	if err != nil {
		log.Error(fmt.Sprintf("Parse Query int error,key=%v ", key))
		log.Error(err)
	}
	return i
}

func ParamInt(c *gin.Context, key string) int {
	si := c.Param(key)
	i, err := strconv.Atoi(si)
	if err != nil {
		log.Error(fmt.Sprintf("Parse Param int error,key=%v ", key))
		log.Error(err)
	}
	return i
}
func BadRequest(c *gin.Context, msg string) {
	JsonError(c, http.StatusBadRequest, msg)
}
func NotFound(c *gin.Context, msg string) {
	JsonError(c, http.StatusNotFound, msg)
}
func ServerError(c *gin.Context, msg string) {
	JsonError(c, http.StatusInternalServerError, msg)
}

func JsonError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"code":    code,
		"msg":     msg,
		"success": false,
	})
}
func JsonData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"msg":     "ok",
		"success": true,
		"data":    data,
	})
}

//totalPage 计算总页数
func totalPage(count, pageSize int) int {
	if pageSize == 0 {
		return 1
	}
	return int(math.Ceil(float64(count) / float64(pageSize)))
}

//parsePageParam 解析page参数
func parsePageParam(c *gin.Context) (page int, pagesize int) {
	page = QueryInt(c, "page")
	if page < 1 {
		page = 1
	}
	pagesize = QueryInt(c, "pagesize")
	if pagesize <= 0 {
		pagesize = 10
	}
	return page, pagesize
}
`
