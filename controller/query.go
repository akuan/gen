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

//BuildBetweenQuery
func BuildBetweenQuery(c *gin.Context, tx *gorm.DB, cols map[string]string) *gorm.DB {
	for pname, col := range cols {
		sname := fmt.Sprintf("%sStart", pname)
		ename := fmt.Sprintf("%sEnd", pname)
		sq := c.Query(sname)
		eq := c.Query(ename)
		if sq != "" && eq != "" {
			log.Debugf("The %s param is start value %s,end value %s ", pname, sq, eq)
			sql := fmt.Sprintf("%s between ? and ?", col)
			tx = tx.Where(sql, sq, eq)
			continue
		}
		if sq != "" {
			log.Debugf("The %s param is start value %s", pname, sq)
			sql := fmt.Sprintf("%s >= ? ", col)
			tx = tx.Where(sql, sq)
			continue
		}
		if eq != "" {
			log.Debugf("The %s param is  end value %s ", pname, eq)
			sql := fmt.Sprintf("%s <= ?  ", col)
			tx = tx.Where(sql, eq)
			continue
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
