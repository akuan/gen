package main

import (
	"fmt"
	"gen/dbmeta"

	//dec "github.com/shopspring/decimal"
	"github.com/jinzhu/gorm"
)

//GenQuery gorm map to gen_query table,for config query tables and columns
type GenQuery struct {
	ID              int    `gorm:"column:id;primary_key" json:"ID"`
	Table           string `gorm:"column:table_name" json:"table"`
	Column          string `gorm:"column:column_name" json:"column"`
	Key             string `gorm:"-" json:"Key"`
	HasEqualQuery   bool   `gorm:"column:has_equal_query" json:"hasEqualQuery"`
	HasBetweenQuery bool   `gorm:"column:has_between_query" json:"hasBetweenQuery"`
	HasLikeQuery    bool   `gorm:"column:has_like_query" json:"hasLikeQuery"`
}

// TableName sets the insert table name for this struct type
func (l *GenQuery) TableName() string {
	return "gen_query"
}

//GetKey use Table::Column as  key
func (l *GenQuery) GetKey() string {
	return fmt.Sprintf("%s::%s", l.Table, l.Column)
}

//QueryCol
type QueryCol struct {
	Stype string //Field name in struct
	Col   string //column  name in database
}

//FindQueryColums use gorm to get all gen_query table data
func FindQueryColums(sqltype, constr string) (map[string]GenQuery, error) {
	var q = make(map[string]GenQuery)
	Db, err := gorm.Open(sqltype, constr)
	if err != nil {
		fmt.Printf("\nOpen DataBase  Error %v ,dbtype=%v,%v", err, sqltype, constr)
		if Db != nil {
			Db.Close()
		}
		panic(err)
	}
	defer Db.Close()
	var gq []*GenQuery
	if err := Db.Find(&gq).Error; err != nil {
		return q, err
	}
	for _, v := range gq {
		q[v.GetKey()] = *v
	}
	return q, nil
}

//EqualQueryColums  find equal Query columns for table "table"
func EqualQueryColums(sqltype, constr, table string) ([]QueryCol, error) {
	Db, err := gorm.Open(sqltype, constr)
	if err != nil {
		fmt.Printf("\nOpen DataBase  Error %v ,dbtype=%v,%v", err, sqltype, constr)
		if Db != nil {
			Db.Close()
		}
		panic(err)
	}
	defer Db.Close()
	var gq []*GenQuery
	err = Db.Where("table_name = ?", table).Where("has_equal_query=?", true).Find(&gq).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nEqualQueryColums I got for table %s len=%d is %v", table, len(gq), gq)
	return buildQueryCol(gq), nil
}

//BetweenQueryColums find between Query columns for table "table"
func BetweenQueryColums(sqltype, constr, table string) ([]QueryCol, error) {
	Db, err := gorm.Open(sqltype, constr)
	if err != nil {
		fmt.Printf("\nOpen DataBase  Error %v ,dbtype=%v,%v", err, sqltype, constr)
		if Db != nil {
			Db.Close()
		}
		panic(err)
	}
	defer Db.Close()
	var gq []*GenQuery
	err = Db.Where("table_name = ?", table).
		Where("has_between_query=?", true).Find(&gq).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nBetweenQueryColums I got for table %s len=%d is %v", table, len(gq), gq)
	return buildQueryCol(gq), nil
}

//EqualQueryColums  find like Query columns for table "table"
func LikeQueryColums(sqltype, constr, table string) ([]QueryCol, error) {
	Db, err := gorm.Open(sqltype, constr)
	if err != nil {
		fmt.Printf("Open DataBase  Error %v ,dbtype=%v,%v", err, sqltype, constr)
		if Db != nil {
			Db.Close()
		}
		panic(err)
	}
	defer Db.Close()
	var gq []*GenQuery
	err = Db.Where("table_name = ?", table).Where("has_like_query=?", true).Find(&gq).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("\nLikeQueryColums I got for table %s len=%d is %v", table, len(gq), gq)
	return buildQueryCol(gq), nil
}

func buildQueryCol(gq []*GenQuery) []QueryCol {
	var res []QueryCol
	for _, que := range gq {
		sc := que.Column
		if sc != "" {
			fn := dbmeta.FmtFieldName(dbmeta.StringifyFirstChar(sc))
			fn = dbmeta.StrFirstToLower(fn)
			col := QueryCol{
				Stype: fn,
				Col:   sc,
			}
			res = append(res, col)
		}
	}
	return res
}
