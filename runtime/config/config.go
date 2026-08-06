package config

type Config struct {
	Server ServerConfig
	App    AppConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type ServerConfig struct {
	Port int
}

func Default() *Config {

	return &Config{

		App: AppConfig{

			Name: "Goduck",

			Env: "development",
		},

		Server: ServerConfig{

			Port: 8080,
		},
	}
}