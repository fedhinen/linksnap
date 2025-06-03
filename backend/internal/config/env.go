package config

import "os"

type Env struct {
	DatabaseDriver string
	DatabaseURL    string
	RedisURL       string
	RedisPassword  string
	ClerkSecretKey string
}

var Envs *Env

func LoadEnv() *Env {
	Envs = &Env{
		DatabaseDriver: os.Getenv("DATABASE_DRIVER"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		RedisURL:       os.Getenv("REDIS_ADDRESS"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		ClerkSecretKey: os.Getenv("CLERK_SECRET_KEY"),
	}

	return Envs
}
