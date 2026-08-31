// Copyright 2012 James Cooper. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gorp

import (
	"fmt"
	"reflect"
)

type SqliteDialect struct {
	suffix string
}

func (d SqliteDialect) QuerySuffix() string { return ";" }

func (d SqliteDialect) ToSqlType(val reflect.Type, maxsize int, isAutoIncr bool) string {
	switch val.Kind() {
	case reflect.Pointer:
		return d.ToSqlType(val.Elem(), maxsize, isAutoIncr)
	case reflect.Bool:
		return "integer"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float64, reflect.Float32:
		return "real"
	case reflect.Slice:
		if val.Elem().Kind() == reflect.Uint8 {
			return "blob"
		}
	}

	switch val.Name() {
	case "NullInt64":
		return "integer"
	case "NullFloat64":
		return "real"
	case "NullBool":
		return "integer"
	case "Time":
		return "datetime"
	}

	if maxsize < 1 {
		maxsize = 255
	}
	return fmt.Sprintf("varchar(%d)", maxsize)
}

// AutoIncrStr Returns autoincrement
func (d SqliteDialect) AutoIncrStr() string {
	return "autoincrement"
}

func (d SqliteDialect) AutoIncrBindValue() string {
	return "null"
}

func (d SqliteDialect) AutoIncrInsertSuffix(_ *ColumnMap) string {
	return ""
}

// CreateTableSuffix Returns suffix
func (d SqliteDialect) CreateTableSuffix() string {
	return d.suffix
}

func (d SqliteDialect) CreateIndexSuffix() string {
	return ""
}

func (d SqliteDialect) DropIndexSuffix() string {
	return ""
}

// TruncateClause With sqlite, there technically isn't a TRUNCATE statement,
// but a DELETE FROM uses a truncate optimization:
// http://www.sqlite.org/lang_delete.html
func (d SqliteDialect) TruncateClause() string {
	return "delete from"
}

// BindVar Returns "?"
func (d SqliteDialect) BindVar(_ int) string {
	return "?"
}

func (d SqliteDialect) InsertAutoIncr(exec SqlExecutor, insertSql string, params ...any) (int64, error) {
	return standardInsertAutoIncr(exec, insertSql, params...)
}

func (d SqliteDialect) QuoteField(f string) string {
	return `"` + f + `"`
}

// QuotedTableForQuery sqlite does not have schemas like PostgreSQL does, so just escape it like normal
func (d SqliteDialect) QuotedTableForQuery(_ string, table string) string {
	return d.QuoteField(table)
}

func (d SqliteDialect) IfSchemaNotExists(command, _ string) string {
	return fmt.Sprintf("%s if not exists", command)
}

func (d SqliteDialect) IfTableExists(command, _, _ string) string {
	return fmt.Sprintf("%s if exists", command)
}

func (d SqliteDialect) IfTableNotExists(command, _, _ string) string {
	return fmt.Sprintf("%s if not exists", command)
}
