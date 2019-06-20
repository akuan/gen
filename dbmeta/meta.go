package dbmeta

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jimsmart/schema"
)

type ModelInfo struct {
	PackageName     string
	StructName      string
	ShortStructName string
	TableName       string
	Fields          []string
	HasTimeFiled    bool
}

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}

var intToWordMap = []string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

// Constants for return types of golang
const (
	golangByteArray  = "[]byte"
	gureguNullInt    = "null.Int"
	sqlNullInt       = "sql.NullInt64"
	golangBool       = "bool"
	golangInt        = "int"
	golangInt64      = "int64"
	gureguNullFloat  = "null.Float"
	sqlNullFloat     = "sql.NullFloat64"
	golangFloat      = "float"
	golangFloat32    = "float32"
	golangFloat64    = "float64"
	gureguNullString = "null.String"
	sqlNullString    = "sql.NullString"
	gureguNullTime   = "null.Time"
	golangTime       = "time.Time"
	golangDecimal    = "decimal.Decimal"
	custDateTime     = "DateTime"
	custDate         = "Date"
	custTime         = "STime"
)

// GenerateStruct generates a struct for the given table.
func GenerateStruct(db *sql.DB, allStruct map[string]string, tableName string, structName string, pkgName string, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) *ModelInfo {
	cols, _ := schema.Table(db, tableName)
	fields, hasTime := generateFieldsTypes(db, allStruct, cols, 0, jsonAnnotation, gormAnnotation, gureguTypes)

	//fields := generateMysqlTypes(db, columnTypes, 0, jsonAnnotation, gormAnnotation, gureguTypes)

	var modelInfo = &ModelInfo{
		PackageName:     pkgName,
		StructName:      structName,
		TableName:       tableName,
		ShortStructName: strings.ToLower(string(structName[0])),
		Fields:          fields,
		HasTimeFiled:    hasTime,
	}

	return modelInfo
}

// Generate fields string
func generateFieldsTypes(db *sql.DB, allStruct map[string]string, columns []*sql.ColumnType, depth int, jsonAnnotation bool,
	gormAnnotation bool, gureguTypes bool) ([]string, bool) {

	//sort.Strings(keys)

	var fields []string
	var field = ""
	var hasTime = false
	for i, c := range columns {
		nullable, _ := c.Nullable()
		key := c.Name()
		valueType := sqlTypeToGoType(strings.ToLower(c.DatabaseTypeName()), nullable, gureguTypes)
		if valueType == "" { // unknown type
			fmt.Printf(" unknown type %s \n", c.DatabaseTypeName())
			continue
		}
		if valueType == golangTime {
			hasTime = true
		}
		fieldName := FmtFieldName(stringifyFirstChar(key))

		var annotations []string
		if gormAnnotation == true {
			if i == 0 {
				annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s;primary_key\"", key))
			} else {
				annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s\"", key))
			}

		}
		if jsonAnnotation == true {
			jsAnn := strFirstToLower(fieldName)
			annotations = append(annotations, fmt.Sprintf("json:\"%s\"", jsAnn))
			if strings.HasPrefix(valueType, "Date") || strings.HasSuffix(valueType, "Time") {
				annotations = append(annotations, "swaggertype:\"primitive,string\"")
			}
		}
		if len(annotations) > 0 {
			field = fmt.Sprintf("%s %s `%s`",
				fieldName,
				valueType,
				strings.Join(annotations, " "))

		} else {
			field = fmt.Sprintf("%s %s",
				fieldName,
				valueType)
		}
		fields = append(fields, field)
		cn := c.Name()
		cn = strings.ToLower(cn)
		if strings.HasSuffix(cn, "_id") {
			refFieldName := fieldName[:len(fieldName)-2]
			strName := refFieldName
			_, ok := allStruct[refFieldName]
			if !ok {
				strName = "DicValue"
			}
			refField := genReferFild(c.Name(), refFieldName, strName, jsonAnnotation, gormAnnotation)
			fields = append(fields, refField)
		}
	}
	return fields, hasTime
}

func genReferFild(reginalName, fName, sType string, jsonAnnotation bool, gormAnnotation bool) string {
	var field = ""
	var annotations []string
	if gormAnnotation == true {
		annotations = append(annotations, fmt.Sprintf("gorm:\"save_associations:fase;foreignkey:%s\"", reginalName))
	}
	if jsonAnnotation == true {
		jsAnn := strFirstToLower(fName)
		annotations = append(annotations, fmt.Sprintf("json:\"%s\"", jsAnn))
	}
	if len(annotations) > 0 {
		field = fmt.Sprintf("%s %s `%s`",
			fName,
			sType,
			strings.Join(annotations, " "))

	} else {
		field = fmt.Sprintf("%s %s",
			fName,
			sType)
	}
	return field
}

func sqlTypeToGoType(mysqlType string, nullable bool, gureguTypes bool) string {
	switch mysqlType {
	case "tinyint", "int", "integer", "smallint", "int2", "int4", "mediumint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt
	case "bigint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt64
	case "bit", "bool", "boolean":
		return golangBool
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext":
		if nullable {
			if gureguTypes {
				return gureguNullString
			}
			return sqlNullString
		}
		return "string"
	//
	// case "date", "datetime", "time", "timestamp":
	// 	if nullable && gureguTypes {
	// 		return gureguNullTime
	// 	}
	// 	return golangTime
	//
	case "datetime", "timestamp", "timestamp without time zone", "timestamp with time zone":
		if nullable && gureguTypes {
			return gureguNullTime
		}
		return custDateTime
	case "date":
		if nullable && gureguTypes {
			return gureguNullTime
		}
		return custDate
	case "time", "time without time zone", "time with time zone":
		if nullable && gureguTypes {
			return gureguNullTime
		}
		return custTime
	case "decimal", "numeric", "double":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat64
	case "float":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat32
	case "binary", "blob", "longblob", "mediumblob", "varbinary", "bytea":
		return golangByteArray
	}
	return ""
}
