package main

import (
	"math/rand"
	"thanhldt060802/app"
	"thanhldt060802/infrastructure"
	"thanhldt060802/repository"
	"time"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	infrastructure.InitPostgesDB()
	defer infrastructure.PostgresDB.Close()
	repository.InitTaskRepository()

	myApp := app.NewDistributedSystemHandleTask("node1@localhost", "node2@localhost", 5)
	myApp.Start()
}
