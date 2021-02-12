package main

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	a := App{}

	a.Initialize()

	defer a.DB.Close()
	a.Migrations()

	// log.Println("Server running on Port 8080...")
	a.Run(":8000")
}
