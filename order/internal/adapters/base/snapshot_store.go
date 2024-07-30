package base

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	DefaultSnapshotTableName = "snapshots"
	loadSnapshotSQL          = "SELECT snapshot_name, snapshot_data, snapshot_version FROM %s WHERE entity_name = $1 AND entity_id = $2 LIMIT 1"
	saveSnapshotSQL          = `INSERT INTO %s (entity_name, entity_id, snapshot_name, snapshot_data, snapshot_version, modified_at) 
VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP) 
ON CONFLICT (entity_name, entity_id) DO
UPDATE SET snapshot_name = EXCLUDED.snapshot_name, snapshot_data = EXCLUDED.snapshot_data, snapshot_version = EXCLUDED.snapshot_version, modified_at = EXCLUDED.modified_at`
	CreateSnapshotsTableSQL = `CREATE TABLE %s (
		entity_name      text        NOT NULL,
		entity_id        text        NOT NULL,
		snapshot_name    text        NOT NULL,
		snapshot_data    bytea       NOT NULL,
		snapshot_version int         NOT NULL,
		modified_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (entity_name, entity_id)
	)`
)

type SnapshotStore struct {
	tableName string
	client    Client
	strategy  SnapshotStrategy
	next      Store
}

func NewSnapshotStore(client Client, options ...SnapshotStoreOption) StoreMiddleware {
	store := &SnapshotStore{
		tableName: DefaultSnapshotTableName,
		client:    client,
		strategy:  DefaultSnapshotStrategies,
	}

	for _, option := range options {
		option(store)
	}

	err := client.Migrate(store.tableName, CreateSnapshotsTableSQL)
	if err != nil {
		panic(err)
	}

	return func(next Store) Store {
		store.next = next
		return store
	}
}

func (s *SnapshotStore) Load(ctx context.Context, root *AggregateRoot) error {
	name := root.AggregateName()
	id := root.AggregateID()

	row := s.client.QueryRow(ctx, fmt.Sprintf(loadSnapshotSQL, s.tableName), name, id)

	var data []byte
	var version int

	err := row.Scan(&data, &version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s.next.Load(ctx, root)
		}
		return err
	}

	snapshot := root.GetSnapshotType()
	err = json.Unmarshal(data, &snapshot)
	if err != nil {
		return err
	}

	err = root.LoadSnapshot(snapshot, version)
	if err != nil {
		return err
	}

	return nil
}

func (s *SnapshotStore) Save(ctx context.Context, root *AggregateRoot) error {
	err := s.next.Save(ctx, root)
	if err != nil {
		return err
	}

	if !s.strategy.ShouldSnapshot(root) {
		return nil
	}

	snapshot, err := root.Aggregate().ToSnapshot()
	if err != nil {
		return err
	}

	name := root.AggregateName()
	id := root.AggregateID()
	version := root.PendingVersion()
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}

	err = s.client.Exec(ctx, fmt.Sprintf(saveSnapshotSQL, s.tableName), name, id, snapshot.SnapshotName(), data, version)
	if err != nil {
		return err
	}

	return nil
}

type SnapshotStoreOption func(*SnapshotStore)

func WithSnapshotStoreTableName(tableName string) SnapshotStoreOption {
	return func(store *SnapshotStore) {
		store.tableName = tableName
	}
}

func WithSnapshotStoreStrategy(strategy SnapshotStrategy) SnapshotStoreOption {
	return func(store *SnapshotStore) {
		store.strategy = strategy
	}
}
