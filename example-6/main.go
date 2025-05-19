package main

import (
	"math/rand"
	"thanhldt060802/app"
	"time"
)

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	myApp := app.NewDistributedSystemHandleTask("node1@localhost", "node2@localhost", 1)
	myApp.Start()
}
