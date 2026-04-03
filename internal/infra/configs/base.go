package configs

import env "github.com/caarlos0/env/v11"

var confList = []any{
	&AppConfig,
}

func ParseConfig() error {
	for _, conf := range confList {
		if err := env.Parse(conf); err != nil {
			return err
		}
	}
	return nil
}
