package main

import (
	"go-crud-article/connection"
	"go-crud-article/handlers"
)

func main() {
	connection.Connect()

	handlers.HandleReq()
}
