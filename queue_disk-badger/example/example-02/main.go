package main

import (
	"fmt"
	"log"
	"thanhldt060802/queue"
	"time"
)

var BatchQueueDiskInstance *queue.BatchQueueDisk

func main() {

	BatchQueueDiskInstance = queue.NewBatchQueueDisk("disk_storage", 100)

	data := make([]string, 10000)
	for i := 0; i < len(data); i++ {
		data[i] = fmt.Sprintf("message %v", i)
	}
	EnqueueDemo(data)

	DequeueDemo(false)

	BatchQueueDiskInstance.Close()

}

func EnqueueDemo(data []string) {
	count := 0
	startTime := time.Now()
	for _, element := range data {
		if err := BatchQueueDiskInstance.Enqueue(element); err != nil {
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
			values, err := BatchQueueDiskInstance.Dequeue()
			if err != nil {
				break
			}
			for _, value := range values {
				fmt.Println(value)
			}
			count += len(values)
		}
	} else {
		for {
			values, err := BatchQueueDiskInstance.Dequeue()
			if err != nil {
				break
			}
			count += len(values)
		}
	}
	endTime := time.Now()
	log.Printf("Total time for dequeue %v elements: %v\n", count, endTime.Sub(startTime))
}
