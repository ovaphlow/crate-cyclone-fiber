package genericimplementation

import (
	"encoding/json"
	"fmt"
	"ovaphlow/cratecyclone/utilities"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

type Column struct {
	OrdinalPosition int    `json:"ordinalPosition"`
	ColumnName      string `json:"columnName"`
	DataType        string `json:"dataType"`
}

func create(schema *string, table *string, data map[string]interface{}) error {
	columns, err := retrieveColumns(schema, table)
	if err != nil {
		return err
	}
	var c []string
	var v []string
	for _, column := range columns {
		if column.ColumnName == "id" && column.DataType == "bigint" {
			node, err := snowflake.NewNode(1)
			if err != nil {
				return err
			}
			data["id"] = node.Generate().Int64()
		}
		if column.ColumnName == "time" {
			data["time"] = time.Now().Format("2006-01-02 15:04:05")
		}
		if column.ColumnName == "state" && column.DataType == "jsonb" {
			randomUUID, err := uuid.NewRandom()
			if err != nil {
				return err
			}
			state := map[string]interface{}{
				"created_at": time.Now().Format("2006-01-02 15:04:05"),
				"uuid":       randomUUID,
			}
			stateJson, err := json.Marshal(state)
			if err != nil {
				return err
			}
			data["state"] = string(stateJson)
		}
		c = append(c, column.ColumnName)
		v = append(v, fmt.Sprintf("'%v'", data[column.ColumnName]))
	}
	q := fmt.Sprintf(
		`insert into %s.%s (%s) values (%s)`,
		*schema,
		*table,
		strings.Join(c, ", "),
		strings.Join(v, ", "),
	)
	_, err = utilities.Postgres.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func remove(schema *string, table *string, id *string, uuid *string) error {
	columns, err := retrieveColumns(schema, table)
	if err != nil {
		return err
	}
	var hasState bool
	var q string
	for _, column := range columns {
		if column.ColumnName == "state" && column.DataType == "jsonb" {
			hasState = true
			q = fmt.Sprintf(
				`
				update %s.%s
				set state = state || jsonb_build_object('deleted_at', '%s')
				where id = $1 and state->>'uuid' = '%s'
				`,
				*schema,
				*table,
				time.Now().Format("2006-01-02 15:04:05"),
				*uuid,
			)
		} else {
			q = fmt.Sprintf(
				`
				update %s.%s
				set state = state || jsonb_build_object('deleted_at', '%s')
				where id = $1
				`,
				*schema,
				*table,
				time.Now().Format("2006-01-02 15:04:05"),
			)
		}
	}
	if !hasState {
		return fmt.Errorf("table %s.%s does not have state(jsonb) column", *schema, *table)
	}
	_, err = utilities.Postgres.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}

func retrieveColumns(schema *string, table *string) ([]Column, error) {
	q := `
	select ordinal_position, column_name, data_type
	from information_schema.columns
	where table_schema = $1
		and table_name = $2
	`
	result, err := utilities.Postgres.Query(q, schema, table)
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

func retrieve(schema *string, table *string, o *Option, f *Filter) ([]map[string]interface{}, error) {
	columns, err := retrieveColumns(schema, table)
	if err != nil {
		return nil, err
	}
	var c []string
	for _, column := range columns {
		c = append(c, column.ColumnName)
	}
	q := fmt.Sprintf(
		`select %s from %s.%s where not (state ? 'deleted_at')`,
		strings.Join(c, ", "),
		*schema,
		*table,
	)
	var conditions []string
	var params []string
	if len(f.ArrayContain) > 0 {
		c, p := ArrayContain(f.ArrayContain)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(f.Equal) > 0 {
		c, p := Equal(f.Equal)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(f.Greater) > 0 {
		c, p := Greater(f.Greater)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(f.In) > 0 {
		c, p := In(f.In)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(f.Lesser) > 0 {
		c, p := Lesser(f.Lesser)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(f.Like) > 0 {
		c, p := Like(f.Like)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(f.ObjectContain) > 0 {
		c, p := ObjectContain(f.ObjectContain)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(f.ObjectLike) > 0 {
		c, p := ObjectLike(f.ObjectLike)
		conditions = append(conditions, c...)
		params = append(params, p...)
	}
	if len(conditions) > 0 {
		q += " and "
		q_ := strings.Join(conditions, " and ")
		for i := 0; i < len(params); i++ {
			q_ = strings.Replace(q_, "?", fmt.Sprintf("$%d", i+1), 1)
		}
		q += q_
	}
	q = fmt.Sprintf("%s order by id desc limit %d offset %d", q, o.Take, o.Skip)
	var params_ []interface{}
	for _, it := range params {
		params_ = append(params_, it)
	}
	utilities.Slogger.Info(q)
	utilities.Slogger.Info(fmt.Sprint(params_))
	rows, err := utilities.Postgres.Query(q, params_...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []map[string]interface{}
	columnNames, _ := rows.Columns()
	values := make([]interface{}, len(columnNames))
	valuePtrs := make([]interface{}, len(columnNames))
	for i := range columnNames {
		valuePtrs[i] = &values[i]
	}
	for rows.Next() {
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

func retrieveByID(schema *string, table *string, id *string, uuid *string) (map[string]interface{}, error) {
	columns, err := retrieveColumns(schema, table)
	if err != nil {
		return nil, err
	}
	var c []string
	for _, column := range columns {
		c = append(c, column.ColumnName)
	}
	var q string
	for _, column := range columns {
		if column.ColumnName == "state" && column.DataType == "jsonb" {
			q = fmt.Sprintf(
				`select %s from %s.%s where id = $1 and state->>'uuid' = '%s' limit 1`,
				strings.Join(c, ", "),
				*schema,
				*table,
				*uuid,
			)
		} else {
			q = fmt.Sprintf(
				`select %s from %s.%s where id = $1 limit 1`,
				strings.Join(c, ", "),
				*schema,
				*table,
			)
		}
	}
	rows, err := utilities.Postgres.Query(q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columnNames, _ := rows.Columns()
	values := make([]interface{}, len(columnNames))
	valuePtrs := make([]interface{}, len(columnNames))
	for i := range columnNames {
		valuePtrs[i] = &values[i]
	}
	var result map[string]interface{}
	for rows.Next() {
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
		result = rowData
	}
	return result, nil
}

func update(schema *string, table *string, id *string, uuid *string, data map[string]interface{}) error {
	columns, err := retrieveColumns(schema, table)
	if err != nil {
		return err
	}
	var s []string
	for _, column := range columns {
		if column.ColumnName == "state" && column.DataType == "jsonb" {
			s = append(
				s,
				fmt.Sprintf(
					`state = state || jsonb_build_object('updated_at', '%s')`,
					time.Now().Format("2006-01-02 15:04:05"),
				),
			)
			continue
		}
		if _, ok := data[column.ColumnName]; ok {
			s = append(s, fmt.Sprintf("%s = '%v'", column.ColumnName, data[column.ColumnName]))
		}
	}
	var q string
	for _, column := range columns {
		if column.ColumnName == "state" && column.DataType == "jsonb" {
			q = fmt.Sprintf(
				`update %s.%s set %s where id = $1 and state->>'uuid' = '%s'`,
				*schema,
				*table,
				strings.Join(s, ", "),
				*uuid,
			)
		} else {
			q = fmt.Sprintf(
				`update %s.%s set %s where id = $1`,
				*schema,
				*table,
				strings.Join(s, ", "),
			)
		}
	}
	_, err = utilities.Postgres.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}
