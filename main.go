package main

import (
	"ecommerce/database"
	"ecommerce/routes"
)

func main() {
	database.ConnectDatabase()
	r := routes.SetUpRouter()
	r.Run(":8888")
}
