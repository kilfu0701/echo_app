package core

import (
	"context"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kilfu0701/echo_app/config"
)

type AppContext struct {
	echo.Context

	AppDB       *mongo.Database
	AppDBCtx    context.Context
	AppDBClient *mongo.Client

	AppMemcached *memcache.Client

	Env *AppEnv
}

func NewAppContext(c echo.Context, cfg *config.Config) (*AppContext, error) {
	app_env, err := NewAppEnv(cfg)
	if err != nil {
		return nil, err
	}

	db_config := app_env.GetDatabase()
	adb, err := NewAppDB(db_config)
	if err != nil {
		return nil, err
	}

	//memcache_config := app_env.GetMemcache()
	mc := memcache.New("localhost:11211")

	ac := &AppContext{
		Context: c,

		AppDB:       adb.DB,
		AppDBCtx:    adb.Ctx,
		AppDBClient: adb.Client,

		AppMemcached: mc,

		Env: app_env,
	}

	return ac, nil
}
