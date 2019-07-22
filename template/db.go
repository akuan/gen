package template

var DbTmpl = `package  model

import (
	"errors"
	pg "labgo/connections/database/postgresql"
	"reflect"
	"github.com/jinzhu/gorm"
	dec "github.com/shopspring/decimal"
)

var Db *gorm.DB

func init() {
	dec.DivisionPrecision = 2
	Db = pg.OpenPgDb()
	Db.LogMode(true)
	autoMigrate()
	initDbVal()
}
func autoMigrate() {
	{{range .}}Db.AutoMigrate(&{{.}}{})
	{{end}}
}

func initDbVal() {
	//some other init operation
}

func Copy(dst interface{}, src interface{}) error {
	dstV := reflect.Indirect(reflect.ValueOf(dst))
	srcV := reflect.Indirect(reflect.ValueOf(src))
	if !dstV.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}
	if srcV.Type() != dstV.Type() {
		return errors.New("different types can be copied")
	}
	for i := 0; i < dstV.NumField(); i++ {
		f := srcV.Field(i)
		if !isZeroOfUnderlyingType(f.Interface()) {
			dstV.Field(i).Set(f)
		}
	}
	return nil
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

//IsNotFound
func IsNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
`
