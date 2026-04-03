package configs

import env "github.com/caarlos0/env/v11"

type baseConfig struct{}

func (cfg *baseConfig) Parse() error {
	return env.Parse(cfg)
}

var confList = []any{
	&AppConfig,
}

func ParseConfigs() error {
	for _, conf := range confList {
		if err := env.Parse(conf); err != nil {
			return err
		}
	}
	return nil
}
