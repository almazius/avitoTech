package config

import "time"

type Config struct {
	Postgres struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		DbName   string `json:"dbName"`
	} `json:"postgres"`
	ServerPort   string        `json:"server-port"`
	ReadTimeout  time.Duration `json:"read-timeout"`
	WriteTimeout time.Duration `json:"write-timeout"`
}
