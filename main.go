package main

import (
	"customer-profile-crud/connection"
	"customer-profile-crud/handlers"
)

func main() {
	connection.Connect()

	handlers.HandleReq()
}
