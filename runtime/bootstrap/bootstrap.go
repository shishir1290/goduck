package bootstrap

import (
	"github.com/shishir1290/goduck/runtime/container"
	"github.com/shishir1290/goduck/runtime/database"
	"github.com/shishir1290/goduck/runtime/server"
)

func New(port int) *server.Server {

	app := server.New(port)

	db, err := database.Connect(database.Config{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "",
		Database: "",
	})

	if err == nil {

		app.Container().Register(
			(*database.Database)(nil),
			func(c *container.Container) any {
				return db
			},
		)
	}

	return app
}