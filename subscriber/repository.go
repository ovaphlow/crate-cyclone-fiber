package subscriber

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"ovaphlow/cratecyclone/utility"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
)

/*
CREATE TABLE crate.subscriber (
	id int8 NOT NULL,
	email varchar NOT NULL,
	"name" varchar NOT NULL,
	phone varchar NOT NULL,
	tags jsonb NOT NULL,
	detail jsonb NOT NULL,
	"time" timestamptz NOT NULL,
	CONSTRAINT subscriber_pk PRIMARY KEY (id)
);
CREATE INDEX subscriber_email_idx ON crate.subscriber USING btree (email);
CREATE INDEX subscriber_name_idx ON crate.subscriber USING btree (name);
CREATE INDEX subscriber_phone_idx ON crate.subscriber USING btree (phone);
*/

var subscriberColumns = []string{"id", "email", "name", "phone", "tags", "detail", "time", "state"}

func repoCreateSubscriber(subscriber *Subscriber) error {
	q := fmt.Sprintf(
		"insert into crate.subscriber (%s) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		strings.Join(subscriberColumns, ", "),
	)
	node, err := snowflake.NewNode(1)
	if err != nil {
		utility.Slogger.Error(err.Error())
		return err
	}
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	state := map[string]interface{}{
		"uuid":       randomUUID,
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
	}
	stateJson, err := json.Marshal(state)
	if err != nil {
		return err
	}
	_, err = utility.Postgres.Exec(
		q,
		node.Generate(),
		subscriber.Email,
		subscriber.Name,
		subscriber.Phone,
		subscriber.Tags,
		subscriber.Detail,
		time.Now().Format("2006-01-02 15:04:05"),
		stateJson,
	)
	if err != nil {
		return err
	}
	return nil
}

func repoRetrieveSubscriberById(id int64, uuid string) (*Subscriber, error) {
	q := fmt.Sprintf(
		`
		select %s from crate.subscriber
		where id = $1 and state @> jsonb_build_object('uuid', '%s')
		limit 1
		`,
		strings.Join(subscriberColumns, ", "),
		uuid,
	)
	statement, err := utility.Postgres.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	subscriber := &Subscriber{}
	err = statement.QueryRow(id).Scan(
		&subscriber.Id,
		&subscriber.Email,
		&subscriber.Name,
		&subscriber.Phone,
		&subscriber.Tags,
		&subscriber.Detail,
		&subscriber.Time,
		&subscriber.State,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return subscriber, nil
}

func repoRetrieveSubscriberByUsername(username string) (*Subscriber, error) {
	q := fmt.Sprintf(
		"select %s from crate.subscriber where email = $1 or name = $2 or phone = $3",
		strings.Join(subscriberColumns, ", "),
	)
	result, err := utility.Postgres.Query(q, username, username, username)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	var rows []Subscriber
	for result.Next() {
		var subscriber Subscriber
		err = result.Scan(
			&subscriber.Id,
			&subscriber.Email,
			&subscriber.Name,
			&subscriber.Phone,
			&subscriber.Tags,
			&subscriber.Detail,
			&subscriber.Time,
			&subscriber.State,
		)
		if err != nil {
			return nil, err
		}
		rows = append(rows, subscriber)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	if len(rows) > 1 {
		return nil, fmt.Errorf("duplicate subscriber")
	}
	return &rows[0], nil
}
