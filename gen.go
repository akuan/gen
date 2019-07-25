package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"gen/dbmeta"
	gtmpl "gen/template"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/inflection"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/serenize/snaker"
)

//genDbMigrate generate gorm database auto Migrate
func genDbMigrate(structNames []string) {
	//generate db
	tplDb, err := getTemplate(gtmpl.DbTmpl)
	if err != nil {
		fmt.Println("Error in lading db  template")
		return
	}
	var dbBuf bytes.Buffer
	err = tplDb.Execute(&dbBuf, structNames)
	if err != nil {
		fmt.Println("Error in rendering router: " + err.Error())
		return
	}
	data, err := format.Source(dbBuf.Bytes())
	if err != nil {
		fmt.Println("Error in formating source: " + err.Error())
		return
	}
	ioutil.WriteFile(filepath.Join("model", "db.go"), data, 0777)
	//end generate db
}

//generate routers
func genRouters(apiName string, structNames []string) {
	tplRt, err := getTemplate(gtmpl.RouterTmpl)
	if err != nil {
		fmt.Println("Error in lading router template")
		return
	}
	var buf bytes.Buffer
	err = tplRt.Execute(&buf, structNames)
	if err != nil {
		fmt.Println("Error in rendering router: " + err.Error())
		return
	}
	data, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Error in formating source: " + err.Error())
		return
	}
	ioutil.WriteFile(filepath.Join(apiName, "router.go"), data, 0777)
}

//genControllers
func genControllers(apiName, apiRouter, tableName, packageName, structName string, richFields []string) {
	tplCtl, err := getTemplate(gtmpl.ControllerTmpl)
	if err != nil {
		fmt.Println("Error in loading controller template: " + err.Error())
		return
	}
	//add query fields
	v, _ := EqualQueryColums(*sqlType, *sqlConnStr, tableName)
	q, _ := BetweenQueryColums(*sqlType, *sqlConnStr, tableName)
	l, e := LikeQueryColums(*sqlType, *sqlConnStr, tableName)
	fmt.Printf("\n len(EqualQueryCols)=%d,len(BetweenQueryCols)=%d,len(LikeQueryCols)=%d",
		len(v), len(q), len(l))
	if e != nil {
		fmt.Println("\nError in EqualQueryColums : " + e.Error())
	}
	fmt.Printf("\nRichFields is %v : ", richFields)
	//write api
	var buf bytes.Buffer
	err = tplCtl.Execute(&buf, map[string]interface{}{
		"PackageName":      packageName + "/model",
		"StructName":       structName,
		"EqualQueryCols":   v,
		"BetweenQueryCols": q,
		"LikeQueryCols":    l,
		"ApiRouter":        apiRouter,
		"RichFields":       richFields,
	})
	if err != nil {
		fmt.Println("\nError in rendering controller: " + err.Error())
		return
	}
	//
	data, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("\n Error in formating source: " + err.Error())
		fmt.Println(string(buf.Bytes()))
		//	return
	}
	ioutil.WriteFile(filepath.Join(apiName, inflection.Singular(tableName)+".go"), data, 0777)
}

//genModel generate model files
func genModel(modelInfo *dbmeta.ModelInfo, tableName string) {
	var buf bytes.Buffer
	tplModel, err := getTemplate(gtmpl.ModelTmpl)
	if err != nil {
		fmt.Println("Error in loading model template: " + err.Error())
		return
	}
	err = tplModel.Execute(&buf, modelInfo)
	if err != nil {
		fmt.Println("Error in rendering model: " + err.Error())
		return
	}
	data, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("Error in formating source: " + err.Error())
		return
	}
	ioutil.WriteFile(filepath.Join("model", inflection.Singular(tableName)+".go"), data, 0777)
}

//
func getTemplate(t string) (*template.Template, error) {
	var funcMap = template.FuncMap{
		"pluralize":        inflection.Plural,
		"title":            strings.Title,
		"toLower":          strings.ToLower,
		"toLowerCamelCase": camelToLowerCamel,
		"toSnakeCase":      snaker.CamelToSnake,
	}

	tmpl, err := template.New("model").Funcs(funcMap).Parse(t)

	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func camelToLowerCamel(s string) string {
	ss := strings.Split(s, "")
	ss[0] = strings.ToLower(ss[0])

	return strings.Join(ss, "")
}
