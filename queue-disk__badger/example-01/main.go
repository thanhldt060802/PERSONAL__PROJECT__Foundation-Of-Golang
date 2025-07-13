package main

import (
	"fmt"
	"log"
	"thanhldt060802/queuedisk"
	"time"
)

var EXAMPLE_NUM int = 1
var EXAMPLES map[int]func()

var QueueDisk queuedisk.IQueueDisk

func init() {
	EXAMPLES = map[int]func(){
		2: Example2,
	}

	QueueDisk = queuedisk.NewQueueDisk("disk_storage")
}

func main() {

	EXAMPLES[EXAMPLE_NUM]()

}

/*
- Example for Enqueue() and Dequeue()
- Calculate time for performance when handle 10000 element
*/
func Example2() {
	{
		data := make([]string, 10000)
		for i := 0; i < len(data); i++ {
			data[i] = fmt.Sprintf("message %v", i)
		}

		count := 0
		startTime := time.Now()
		for _, element := range data {
			if err := QueueDisk.Enqueue(element); err != nil {
				log.Fatal(err.Error())
			}
			count++
		}
		endTime := time.Now()
		log.Printf("Total time for enqueue %v elements: %v\n", count, endTime.Sub(startTime))
	}

	{
		count := 0
		startTime := time.Now()
		for {
			_, err := QueueDisk.Dequeue()
			if err != nil {
				break
			}
			count++
		}
		endTime := time.Now()
		log.Printf("Total time for dequeue %v elements: %v\n", count, endTime.Sub(startTime))
	}

	QueueDisk.Close()
}
