package main

import (
	"chatty/server/services"
)

func main() {
	s := services.NewServer(&services.Config{
		Host: "localhost",
		Port: "8080",
	})
	s.Run()
}
