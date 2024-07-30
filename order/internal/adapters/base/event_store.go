package base

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	DefaultEventTableName = "events"
	loadEventsSQL         = "SELECT event_name, event_data FROM %s WHERE entity_name = $1 AND entity_id = $2 AND event_version > $3 ORDER BY event_version ASC"
	writeEventSQL         = "INSERT INTO %s (entity_name, entity_id, event_version, event_name, event_data, created_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)"
	CreateEventsTableSQL  = `CREATE TABLE %s (
    	entity_name    text        NOT NULL,
    	entity_id      text        NOT NULL,
		event_version  int         NOT NULL,
		event_name     text        NOT NULL,
		event_data     bytea       NOT NULL,
		created_at     timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (entity_name, entity_id, event_version)
	)`
)

var _ Store = (*EventStore)(nil)

type EventStore struct {
	tableName string
	client    Client
}

func NewEventStore(client Client, options ...EventStoreOption) *EventStore {
	store := &EventStore{
		tableName: DefaultEventTableName,
		client:    client,
	}

	for _, option := range options {
		option(store)
	}

	err := client.Migrate(store.tableName, CreateEventsTableSQL)
	if err != nil {
		panic(err)
	}

	return store
}

func (e *EventStore) Load(ctx context.Context, root *AggregateRoot) error {
	name := root.AggregateName()
	id := root.AggregateID()
	version := root.PendingVersion()

	row := e.client.QueryRow(ctx, fmt.Sprintf(loadEventsSQL, e.tableName), name, id, version)

	var eventName string
	var data []byte

	err := row.Scan(&eventName, &data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	event := root.GetEvent(eventName)
	err = json.Unmarshal(data, &event)
	if err != nil {
		return err
	}

	err = root.LoadEvent(event)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventStore) Save(ctx context.Context, root *AggregateRoot) (err error) {
	name := root.AggregateName()
	id := root.AggregateID()
	version := root.Version()

	for i, event := range root.Events() {
		var data []byte

		data, err = json.Marshal(event)
		if err != nil {
			return err
		}
		err = e.client.Exec(ctx, fmt.Sprintf(writeEventSQL, e.tableName), name, id, version+i+1, event.EventName(), data)
		if err != nil {
			return err
		}
	}

	return nil
}

type EventStoreOption func(*EventStore)

func WithEventTableName(tableName string) EventStoreOption {
	return func(store *EventStore) {
		store.tableName = tableName
	}
}
