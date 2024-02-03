package schema

import "ovaphlow/cratecyclone/utilities"

func retrieveSchemas() ([]string, error) {
	q := "select schema_name from information_schema.schemata"
	statement, err := utilities.Postgres.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	result, err := statement.Query()
	if err != nil {
		return nil, err
	}
	var schemas []string
	for result.Next() {
		var schema string
		err = result.Scan(&schema)
		if err != nil {
			return nil, err
		}
		schemas = append(schemas, schema)
	}
	return schemas, nil
}
