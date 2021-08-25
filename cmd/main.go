package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/lotostudio/financial-api/docs"
	"github.com/lotostudio/financial-api/internal/app"
)

const configPath = "configs/main.yml"

func main() {
	app.Run(configPath)
}
