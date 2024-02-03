package table

import (
	"ovaphlow/cratecyclone/utilities"
)

func retrieveTables(schema string) ([]string, error) {
	q := `
	select table_name
	from information_schema.tables
	where table_schema = $1;
	`
	statement, err := utilities.Postgres.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	result, err := statement.Query(schema)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	var tables []string
	for result.Next() {
		var table string
		err = result.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}
