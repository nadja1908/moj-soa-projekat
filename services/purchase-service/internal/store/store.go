package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	Client *mongo.Client
	DB     *mongo.Database

	CartsCollection  *mongo.Collection
	TokensCollection *mongo.Collection
}

func NewStore(uri, dbName string) (*Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(dbName)

	store := &Store{
		Client:           client,
		DB:               db,
		CartsCollection:  db.Collection("shoppingCarts"),
		TokensCollection: db.Collection("tourPurchaseTokens"),
	}

	fmt.Println("MongoDB connection established")
	return store, nil
}

func (s *Store) Close(ctx context.Context) error {
	return s.Client.Disconnect(ctx)
}
