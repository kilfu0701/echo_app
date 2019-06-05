package core

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kilfu0701/echo_app/config"
)

type AppDB struct {
	Ctx    context.Context
	DB     *mongo.Database
	Client *mongo.Client
}

// @params
//   mongodb_uri = mongodb://localhost:27017
//   db_name = app
//   timeout = 10
func NewAppDB(db_config config.DBConfig) (*AppDB, error) {
	mongodb_uri := fmt.Sprintf("%s://%s:%d", db_config.Driver, db_config.Host, db_config.Port)

	client, err := mongo.NewClient(options.Client().ApplyURI(mongodb_uri))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(db_config.Timeout)*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	//defer client.Disconnect(ctx)

	adb := &AppDB{
		Ctx:    ctx,
		DB:     client.Database(db_config.DBName),
		Client: client,
	}
	return adb, nil
}
