package main

import (
	"fmt"
	"log"
	"thanhldt060802/queue"
	"time"
)

var QueueDiskInstance *queue.QueueDisk

func main() {

	QueueDiskInstance = queue.NewQueueDisk("disk_storage")

	data := make([]string, 10000)
	for i := 0; i < len(data); i++ {
		data[i] = fmt.Sprintf("message %v", i)
	}
	EnqueueDemo(data)

	DequeueDemo(false)

	QueueDiskInstance.Close()

}

func EnqueueDemo(data []string) {
	count := 0
	startTime := time.Now()
	for _, element := range data {
		if err := QueueDiskInstance.Enqueue(element); err != nil {
			log.Fatal(err.Error())
		}
		count++
	}
	endTime := time.Now()
	log.Printf("Total time for enqueue %v elements: %v\n", count, endTime.Sub(startTime))
}

func DequeueDemo(printable bool) {
	count := 0
	startTime := time.Now()
	if printable {
		for {
			value, err := QueueDiskInstance.Dequeue()
			if err != nil {
				break
			}
			fmt.Println(value)
			count++
		}
	} else {
		for {
			_, err := QueueDiskInstance.Dequeue()
			if err != nil {
				break
			}
			count++
		}
	}
	endTime := time.Now()
	log.Printf("Total time for dequeue %v elements: %v\n", count, endTime.Sub(startTime))
}
