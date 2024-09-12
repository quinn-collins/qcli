package app

type App struct {
	Config *Config
}

func New() *App {
	cfg := NewConfig()

	return &App{
		Config: cfg,
	}
}
