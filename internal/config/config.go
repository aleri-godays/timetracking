package config

//Config holds all configuration for this service
type Config struct {
	ServiceName string
	Version     string
	LogLevel    string
	DbPath      string
	HTTPPort    int
	JWTSecret   string
}
