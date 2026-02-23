package config

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() Config {
	return Config{
		DBHost:     "localhost",
		DBPort:     "5433",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "ai_interview",
	}
}
