package models

type Config struct {
	DB     DB
	Server Server
}

type DB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Server struct {
	Port string
}
