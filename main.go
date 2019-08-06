package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"gen/dbmeta"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/droundy/goopt"
	"github.com/jimsmart/schema"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jinzhu/inflection"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	sqlType     = goopt.String([]string{"--sqltype"}, "mysql", "sql database type such as mysql, postgres, etc.")
	sqlConnStr  = goopt.String([]string{"-c", "--connstr"}, "nil", "database connection string")
	sqlDatabase = goopt.String([]string{"-d", "--database"}, "nil", "Database to for connection")
	sqlTable    = goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")
	apiPath     = goopt.String([]string{"-p", "--apipath"}, "", "RESTful api api path")

	packageName = goopt.String([]string{"--package"}, "", "name to set for package")

	jsonAnnotation = goopt.Flag([]string{"--json"}, []string{"--no-json"}, "Add json annotations (default)", "Disable json annotations")
	gormAnnotation = goopt.Flag([]string{"--gorm"}, []string{}, "Add gorm annotations (tags)", "")
	gureguTypes    = goopt.Flag([]string{"--guregu"}, []string{}, "Add guregu null types", "")

	rest = goopt.Flag([]string{"--rest"}, []string{}, "Enable generating RESTful api", "")

	verbose = goopt.Flag([]string{"-v", "--verbose"}, []string{}, "Enable verbose output", "")
)

func init() {
	// Setup goopts
	goopt.Description = func() string {
		return "ORM and RESTful API generator for Mysql"
	}
	goopt.Version = "0.2"
	goopt.Summary = `gen [-v] --connstr "user:password@/dbname" --package pkgName --database databaseName --table tableName [--json] [--gorm] [--guregu]`

	//Parse options
	goopt.Parse(nil)

}

const (
	gen_query = "gen_query"
)

func main() {
	// Username is required
	if sqlConnStr == nil || *sqlConnStr == "" {
		fmt.Println("\nsql connection string is required! Add it with --connstr=s")
		return
	}

	if sqlDatabase == nil || *sqlDatabase == "" {
		fmt.Println("\nDatabase can not be null")
		return
	}

	var db, err = sql.Open(*sqlType, *sqlConnStr)
	if err != nil {
		fmt.Println("\nError in open database: " + err.Error())
		return
	}
	defer db.Close()

	// parse or read tables
	var tables []string
	if *sqlTable != "" {
		tables = strings.Split(*sqlTable, ",")
	} else {
		tables, err = schema.TableNames(db)
		if err != nil {
			fmt.Println("\nError in fetching tables information from mysql information schema")
			return
		}
	}
	// if packageName is not set we need to default it
	if packageName == nil || *packageName == "" {
		*packageName = "generated"
	}
	os.Mkdir("model", 0777)
	ctlPackage := "controller"
	//apiRouter := *apiPath //"/man/v1/api/"
	if *rest {
		os.Mkdir(ctlPackage, 0777)
	}
	var structNames []string
	var allStruct = make(map[string]string)
	for _, tableName := range tables {
		fmt.Printf("\n\ntableName %v ", tableName)
		if gen_query == tableName {
			continue
		}
		structName := dbmeta.FmtFieldName(tableName)
		structName = inflection.Singular(structName)
		//	structNames = append(structNames, structName)
		allStruct[structName] = tableName
	}

	// generate go files for each table
	for _, tableName := range tables {
		fmt.Printf("\n\ntableName %v ", tableName)
		if gen_query == tableName {
			continue
		}
		structName := dbmeta.FmtFieldName(tableName)
		structName = inflection.Singular(structName)
		structNames = append(structNames, structName)
		modelInfo := dbmeta.GenerateStruct(db, allStruct, tableName, structName, "model", *jsonAnnotation, *gormAnnotation, *gureguTypes)
		genModel(modelInfo, tableName)
		if *rest {
			genControllers(ctlPackage, *apiPath, tableName, *packageName, structName, modelInfo.RichFields)
		}
	}
	//genDbMigrate
	genDbMigrate(structNames)
	//RouterTmpl
	if *rest {
		genRouters(ctlPackage, structNames)
	}
}
