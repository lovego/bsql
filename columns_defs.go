package bsql

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/lovego/struct_tag"
)

var typesMap = make(map[string]bool)

func init() {
	var types = []string{
		"text", "character", "char", "varchar",
		"date", "time", "timetz", "timestamp", "timestamptz", "interval",
		"boolean", "bool",
		"json", "jsonb", "xml",

		"bigint", "int8", "integer", "int", "int4", "smallint", "int2",
		"bigserial", "serial8", "serial", "serial4", "smallserial", "serial2",
		"real", "float4", "double precision", "float8",
		"decimal", "numeric", "money",

		"bit", "varbit", "bytea",
		"inet", "cidr", "macaddr", "macaddr8",
		"box", "circle", "line", "lseg", "path", "point", "polygon",
		"tsquery", "tsvector",
		"uuid", "txid_snapshot", "pg_lsn",
	}
	for _, typ := range types {
		typesMap[typ] = true
	}
}

func ColumnsDefs(strct interface{}) string {
	return strings.Join(columnsFromStruct(strct), ",\n")
}

func columnsFromStruct(model interface{}) []string {
	columns := make([]string, 0)
	traverseStructFields(reflect.TypeOf(model), func(field reflect.StructField) {
		columns = append(columns, Field2Column(field.Name)+" "+getColumnDefinition(field))
	})
	return columns
}

func getColumnDefinition(field reflect.StructField) string {
	var def []string
	tag, ok := struct_tag.Lookup(string(field.Tag), `sql`)
	if ok {
		tag = strings.TrimSpace(tag)
	}
	if hasColumnType(tag) {
		def = append(def, tag)
	} else {
		def = append(def, getColumnType(field))
	}
	if !hasNullConstraint(tag) {
		def = append(def, "NOT NULL")
	}
	if field.Name == "Id" && !hasPrimaryKeyConstraint(tag) {
		def = append(def, "PRIMARY KEY")
	}
	if tag != "" && tag != "-" && !hasColumnType(tag) {
		def = append(def, tag)
	}
	return strings.Join(def, " ")
}

var firstWordRegexp = regexp.MustCompile("^\\w+")
var nullConstraintRegexp = regexp.MustCompile("(?i)\\bnull\\b")
var primaryKeyConstraintRegexp = regexp.MustCompile("(?i)\\bprimary\\s+key\\b")

func hasColumnType(s string) bool {
	word := firstWordRegexp.FindString(s)
	return typesMap[strings.ToLower(word)]
}

func hasNullConstraint(s string) bool {
	return nullConstraintRegexp.MatchString(s)
}

func hasPrimaryKeyConstraint(s string) bool {
	return primaryKeyConstraintRegexp.MatchString(s)
}

func getColumnType(field reflect.StructField) string {
	typ := field.Type
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.String:
		return "text"
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
		if field.Name == "Id" {
			return "serial8"
		} else {
			return "int8"
		}
	case reflect.Int32, reflect.Uint32:
		if field.Name == "Id" {
			return "serial4"
		} else {
			return "int4"
		}
	case reflect.Int8, reflect.Int16, reflect.Uint8, reflect.Uint16:
		if field.Name == "Id" {
			return "serial2"
		} else {
			return "int2"
		}
	case reflect.Bool:
		return "bool"
	case reflect.Float32:
		return "float4"
	case reflect.Float64:
		return "float8"
	default:
		switch typ.Name() {
		case "Time":
			return "timestamptz"
		case "Date":
			return "date"
		case "Decimal":
			return "decimal"
		case "BoolArray":
			return "bool[]"
		case "Int64Array":
			return "bigint[]"
		case "StringArray":
			return "text[]"
		default:
			return "jsonb"
		}
	}
}
