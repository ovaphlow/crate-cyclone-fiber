package schema

import "ovaphlow/cratecyclone/utilities"

func createSchema(schema string) error {
	q := "create schema " + schema
	statement, err := utilities.Postgres.Prepare(q)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

func removeSchema(schema string) error {
	q := "drop schema " + schema + " cascade"
	statement, err := utilities.Postgres.Prepare(q)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

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

func updateSchema(current, new string) error {
	q := "alter schema " + current + " rename to " + new
	statement, err := utilities.Postgres.Prepare(q)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}
