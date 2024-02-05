package genericimplementation

import (
	"fmt"
	"ovaphlow/cratecyclone/utilities"
	"strings"
)

type Column struct {
	OrdinalPosition int    `json:"ordinalPosition"`
	ColumnName      string `json:"columnName"`
	DataType        string `json:"dataType"`
}

func retrieveColumns(schema *string, table *string) ([]Column, error) {
	q := `
	select ordinal_position, column_name, data_type
	from information_schema.columns
	where table_schema = $1
		and table_name = $2
	`
	statement, err := utilities.Postgres.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	result, err := statement.Query(schema, table)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	var column []Column
	for result.Next() {
		var c Column
		err = result.Scan(&c.OrdinalPosition, &c.ColumnName, &c.DataType)
		if err != nil {
			return nil, err
		}
		column = append(column, c)
	}
	return column, nil
}

func retrieve(schema *string, table *string) ([]map[string]interface{}, error) {
	columns, err := retrieveColumns(schema, table)
	if err != nil {
		return nil, err
	}
	var c []string
	for _, column := range columns {
		c = append(c, column.ColumnName)
	}
	q := fmt.Sprintf(
		`select %s from %s.%s`,
		strings.Join(c, ", "),
		*schema,
		*table,
	)
	rows, err := utilities.Postgres.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []map[string]interface{}
	columnNames, _ := rows.Columns()
	values := make([]interface{}, len(columnNames))
	valuePtrs := make([]interface{}, len(columnNames))
	for rows.Next() {
		for i := range columnNames {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}
		rowData := make(map[string]interface{})
		for i, col := range columnNames {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				rowData[col] = string(b)
			} else {
				rowData[col] = val
			}
		}
		results = append(results, rowData)
	}
	return results, nil
}
