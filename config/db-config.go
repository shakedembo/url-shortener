package config

type DbConfig struct {
	Host string
}

const HOST = "localhost" //os.Getenv("DB_HOST")

func NewDbConfig() *DbConfig {
	return &DbConfig{
		Host: HOST,
	}
}
