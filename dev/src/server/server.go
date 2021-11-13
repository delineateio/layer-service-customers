package server

import (
	"fmt"

	"delineate.io/customers/src/config"
	_ "delineate.io/customers/src/docs"
	"delineate.io/customers/src/logging"
	"delineate.io/customers/src/routes"
)

const defaultHost = "localhost"
const defaultPort = "1102"

func Start() {
	router := routes.NewRouter()
	host := config.GetStringOrDefault("server.host", defaultHost)
	port := config.GetStringOrDefault("server.port", defaultPort)
	address := host + ":" + port
	logging.Info(fmt.Sprintf("attempting to start server on '%s'", address))
	if err := router.Run(address); err != nil {
		logging.Err(err)
	}
}
