package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Column struct {
	Name     string
	Default  string
	Nullable string
	Type     string
	Comment  string
}

func main() {
	// Open a database connection
	db, err := sql.Open("mysql", "root:root@123@tcp(localhost:3306)/whalescheduler")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	var schema string = "whalescheduler"
	var tableName string = "t_st_datasource"

	sql := fmt.Sprintf("select TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, COALESCE(COLUMN_DEFAULT, 'NULL'), IS_NULLABLE, DATA_TYPE, COLUMN_TYPE, COLUMN_COMMENT from information_schema.COLUMNS where TABLE_SCHEMA = '%s' and TABLE_NAME = '%s';", schema, tableName)
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var columns = []Column{}
	// 处理查询结果
	for rows.Next() {
		var tableSchema, tableName, columnName, columnDefault, nullable, dataType, columnType, columnComment string
		if err := rows.Scan(&tableSchema, &tableName, &columnName, &columnDefault, &nullable, &dataType, &columnType, &columnComment); err != nil {
			log.Fatal(err)
		}
		column := Column{
			Nullable: nullable,
			Name:     columnName,
			Default:  columnDefault,
			Type:     columnType,
			Comment:  columnComment,
		}
		columns = append(columns, column)
	}

	PrintTable(tableName, columns)

	// 检查是否有错误发生
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func PrintTable(table string, columns []Column) {
	fmt.Printf("# %s \n", table)
	fmt.Println("| 编号 | 字段名 | 注释 | 数据类型 | 索引信息 | 默认 | 是否为空 |")
	fmt.Println("| ---- | ---- | ---- | ---- | ---- | ---- | ---- |")
	for i := 0; i < len(columns); i++ {
		column := columns[i]
		index := If(column.Name == "id", "主键", "")
		fmt.Printf("| %d | %s | %s | %s | %s | %s | %s | \n", i+1, column.Name, GetComment(column.Name, column.Comment), column.Type, index, column.Default, column.Nullable)
	}
}

func GetComment(columnName string, comment string) string {
	if columnName == "id" {
		return "自增ID"
	} else if columnName == "user_id" {
		return "用户ID"
	} else if columnName == "project_id" {
		return "项目ID"
	} else if columnName == "project_code" {
		return "项目code"
	} else if columnName == "process_definition_code" {
		return "工作流定义code"
	} else if columnName == "process_definition_version" {
		return "工作流定义版本"
	} else if columnName == "task_definition_code" {
		return "任务定义code"
	} else if columnName == "task_definition_version" {
		return "任务定义版本"
	} else if columnName == "process_instance_id" {
		return "工作流实例ID"
	} else if columnName == "task_instance_id" {
		return "任务实例ID"
	} else if columnName == "create_time" {
		return "创建时间"
	} else if columnName == "update_time" {
		return "更新时间"
	} else if columnName == "description" {
		return "描述"
	} else {
		return comment
	}
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
