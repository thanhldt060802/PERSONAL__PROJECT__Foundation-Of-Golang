package main

import (
	"math/rand"
	"thanhldt060802/app"
	"time"
)

func main() {

	rand.New(rand.NewSource(time.Now().UnixNano()))

	myApp := app.NewDistributedSystemDemo("node1@localhost", "node2@localhost", 5)
	myApp.Start()

}
