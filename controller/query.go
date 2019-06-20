package controller

import (
	"fmt"
	"labgo/modules/log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//BuildEqualQuery 生成=查询参数
//param cols key-参数名称，value-对应数据库列名
func BuildEqualQuery(c *gin.Context, tx *gorm.DB, cols map[string]string) *gorm.DB {
	for pname, col := range cols {
		q := c.Query(pname)
		if q != "" {
			log.Debugf("The %s param is %s \n", pname, q)
			sql := fmt.Sprintf("%s =?", col)
			tx = tx.Where(sql, q)
		}
	}
	return tx
}

//BuildEqualQuery 生成Like查询参数
//param cols key-参数名称，value-对应数据库列名
func BuildLikeQuery(c *gin.Context, tx *gorm.DB, cols map[string]string) *gorm.DB {
	for pname, col := range cols {
		q := c.Query(pname)
		if q != "" {
			log.Debugf("The %s param is %s \n", pname, q)
			sql := fmt.Sprintf("%s like ?", col)
			tx = tx.Where(sql, q+"%")
		}
	}
	return tx
}
