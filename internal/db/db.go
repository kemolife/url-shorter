package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func Connect(ctx context.Context, dsn string) (*DB, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	return &DB{client: client}, nil
}

func (d *DB) GetClient() *mongo.Client {
	return d.client
}

func (d *DB) CloseConnection(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}
