package config

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

var (
	PostgisDB   *sqlx.DB
	RedisClient *redis.Client
)
