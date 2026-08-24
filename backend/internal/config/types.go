package config

type Config struct {
	Port       string
	LogLevel   string
	LogFormat  string
	CorsOrigin string
	DB         DBConfig
	S3         S3Config
}

type S3Config struct {
	Bucket   string
	Endpoint string
	Region   string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}
