package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/threadpulse/internal/db"
	redisinternal "github.com/threadpulse/internal/db/redis"
	"github.com/threadpulse/internal/routes"

	_ "github.com/threadpulse/internal/routes"
)

func main() {
	db.DBint()
	redisinternal.RedisInit()
	routes.Routes()

}
