package core

import (
	"github.com/kilfu0701/echo_app/config"
)

type AppEnv struct {
	Config *config.Config
}

func NewAppEnv(cfg *config.Config) (*AppEnv, error) {
	a_env := &AppEnv{
		Config: cfg,
	}

	return a_env, nil
}

func (ae *AppEnv) GetEnvName() string {
	return ae.Config.EnvName
}

func (ae *AppEnv) GetDatabase() config.DBConfig {
	return ae.Config.Database
}
