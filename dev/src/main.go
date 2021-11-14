package main

import (
	"delineate.io/customers/src/config"
	"delineate.io/customers/src/database"
	"delineate.io/customers/src/discovery"
	_ "delineate.io/customers/src/docs"
	"delineate.io/customers/src/server"
)

// @title Customer Service
// @version 1.0
// @description Service for managing customers
// @termsOfService http://www.delineate.io/terms/

// @contact.name API Support
// @contact.url http://www.delineate.io/support
// @contact.email support@delineate.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1102
// @BasePath /
// @schemes http
func main() {
	config.Initialize()
	database.Initialize()
	discovery.Initialize()
	server.Start()
}
