package main

import (
	"log"

	"taskmanager/routers"
)

//Entry point of the program
func main() {

	//common.StartUp() - Replaced with init method
	// Get the mux router object
	router := routers.InitRoutes()
	// Create a negroni instance
	log.Println("Listening...")
	router.Run()
}
