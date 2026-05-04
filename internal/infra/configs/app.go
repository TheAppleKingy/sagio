package configs

type appConfig struct {
	Debug              bool   `env:"DEBUG" envDefault:"false"`
	StepKeyName        string `env:"STEP_KEY_NAME" envDefault:"x-saga-step"`
	AppShutdownTimeout int    `env:"APP_SHUTDOWN_TIMEOUT envDefault:20"`
}

var AppConfig = appConfig{}
