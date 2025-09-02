package main

import (
	"fmt"
	"math/rand/v2"
	"thanhldt060802/common/queuedisk"
	"thanhldt060802/model"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

var EXAMPLE_NUM int = 3
var EXAMPLES map[int]func()

func init() {
	EXAMPLES = map[int]func(){
		1: Example1,
		2: Example2,
		3: Example3,
	}
}

func main() {

	EXAMPLES[EXAMPLE_NUM]()

}

// Example for Enqueue() and Dequeue().
func Example1() {
	queuedisk.BatchQueueDiskInstance1 = queuedisk.NewBatchQueueDisk[string]("disk_storage", 8)

	for i := 1; i <= 30; i++ {
		dataEnq := fmt.Sprintf("message %v", i)
		if err := queuedisk.BatchQueueDiskInstance1.Enqueue(dataEnq); err != nil {
			log.Errorf("Enqueue failed: %v", err.Error())
			break
		}
	}

	for {
		dataDeqs, err := queuedisk.BatchQueueDiskInstance1.Dequeue()
		if err != nil {
			log.Errorf("Dequeue failed: %v", err.Error())
			break
		}
		for _, dataDeq := range dataDeqs {
			fmt.Println(dataDeq)
		}
	}

	queuedisk.BatchQueueDiskInstance1.Close()
}

// Ref: Example1(), use data struct.
func Example2() {
	queuedisk.BatchQueueDiskInstance2 = queuedisk.NewBatchQueueDisk[*model.DataStruct]("disk_storage", 8)

	for i := 1; i <= 30; i++ {
		dataEnq := model.DataStruct{
			Field1: uuid.New().String(),
			Field2: rand.Int32(),
			Field3: rand.Int64(),
			Field4: rand.Float32(),
			Field5: rand.Float64(),
			Field6: time.Now(),
			Field7: model.SubDataStruct{
				Field1: uuid.New().String(),
				Field2: rand.Int32(),
				Field3: rand.Int64(),
			},
		}
		if err := queuedisk.BatchQueueDiskInstance2.Enqueue(&dataEnq); err != nil {
			log.Errorf("Enqueue failed: %v", err.Error())
			break
		}
	}

	for {
		dataDeqs, err := queuedisk.BatchQueueDiskInstance2.Dequeue()
		if err != nil {
			log.Errorf("Dequeue failed: %v", err.Error())
			break
		}
		for _, dataDeq := range dataDeqs {
			fmt.Println(*dataDeq)
		}
	}

	queuedisk.BatchQueueDiskInstance2.Close()
}

// Example for Enqueue() and Dequeue().
// Calculate time for performance when handle 10000 element.
func Example3() {
	queuedisk.BatchQueueDiskInstance1 = queuedisk.NewBatchQueueDisk[string]("disk_storage", 33)

	{
		dataEnqs := make([]string, 10000)
		for i := 0; i < len(dataEnqs); i++ {
			dataEnqs[i] = fmt.Sprintf("message %v", i)
		}

		count := 0
		startTime := time.Now()
		for _, dataEnq := range dataEnqs {
			if err := queuedisk.BatchQueueDiskInstance1.Enqueue(dataEnq); err != nil {
				log.Errorf("Enqueue failed: %v", err.Error())
				break
			}
			count++
		}
		endTime := time.Now()
		log.Printf("Total time for enqueue %v elements: %v", count, endTime.Sub(startTime))
	}

	{
		count := 0
		startTime := time.Now()
		for {
			dataDeqs, err := queuedisk.BatchQueueDiskInstance1.Dequeue()
			if err != nil {
				log.Errorf("Dequeue failed: %v", err.Error())
				break
			}
			count += len(dataDeqs)
		}
		endTime := time.Now()
		log.Printf("Total time for dequeue %v elements: %v", count, endTime.Sub(startTime))
	}

	queuedisk.BatchQueueDiskInstance1.Close()
}
