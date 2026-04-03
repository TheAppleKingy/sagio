package configs

type appConfig struct {
	Debug bool `env:"DEBUG" envDefault:"false"`
}

var AppConfig = appConfig{}
