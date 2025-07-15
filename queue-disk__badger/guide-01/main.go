package main

import (
	"fmt"
	"math/rand/v2"
	"thanhldt060802/queuedisk"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var EXAMPLE_NUM int = 1
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

/*
- Example for Enqueue() and Dequeue()
*/
func Example1() {
	queueDisk := queuedisk.NewQueueDisk[string]("disk_storage")

	for i := 1; i <= 30; i++ {
		dataEnq := fmt.Sprintf("message %v", i)
		if err := queueDisk.Enqueue(dataEnq); err != nil {
			log.Errorf("Enqueue failed: %v", err.Error())
			break
		}
	}

	for {
		dataDeq, err := queueDisk.Dequeue()
		if err != nil {
			log.Errorf("Dequeue failed: %v", err.Error())
			break
		}
		fmt.Println(dataDeq)
	}

	queueDisk.Close()
}

/*
- Ref: Example1

- Using data struct
*/
func Example2() {
	type SubDataStruct struct {
		Field1 string `json:"field1"`
		Field2 int32  `json:"field2"`
		Field3 int64  `json:"field3"`
	}
	type DataStruct struct {
		Field1 string        `json:"field1"`
		Field2 int32         `json:"field2"`
		Field3 int64         `json:"field3"`
		Field4 float32       `json:"field4"`
		Field5 float64       `json:"field5"`
		Field6 time.Time     `json:"field6"`
		Field7 SubDataStruct `json:"field7"`
	}

	queueDisk := queuedisk.NewQueueDisk[*DataStruct]("disk_storage")

	for i := 1; i <= 30; i++ {
		dataEnq := DataStruct{
			Field1: uuid.New().String(),
			Field2: rand.Int32(),
			Field3: rand.Int64(),
			Field4: rand.Float32(),
			Field5: rand.Float64(),
			Field6: time.Now(),
			Field7: SubDataStruct{
				Field1: uuid.New().String(),
				Field2: rand.Int32(),
				Field3: rand.Int64(),
			},
		}
		if err := queueDisk.Enqueue(&dataEnq); err != nil {
			log.Errorf("Enqueue failed: %v", err.Error())
			break
		}
	}

	for {
		dataDeq, err := queueDisk.Dequeue()
		if err != nil {
			log.Errorf("Dequeue failed: %v", err.Error())
			break
		}
		fmt.Println(*dataDeq)
	}

	queueDisk.Close()
}

/*
- Example for Enqueue() and Dequeue()

- Calculate time for performance when handle 10000 element
*/
func Example3() {
	queueDisk := queuedisk.NewQueueDisk[string]("disk_storage")

	{
		dataEnqs := make([]string, 10000)
		for i := 0; i < len(dataEnqs); i++ {
			dataEnqs[i] = fmt.Sprintf("message %v", i)
		}

		count := 0
		startTime := time.Now()
		for _, dataEnq := range dataEnqs {
			if err := queueDisk.Enqueue(dataEnq); err != nil {
				log.Errorf("Enqueue failed: %v", err.Error())
				break
			}
			count++
		}
		endTime := time.Now()
		log.Infof("Total time for enqueue %v elements: %v", count, endTime.Sub(startTime))
	}

	{
		count := 0
		startTime := time.Now()
		for {
			_, err := queueDisk.Dequeue()
			if err != nil {
				log.Errorf("Dequeue failed: %v", err.Error())
				break
			}
			count++
		}
		endTime := time.Now()
		log.Infof("Total time for dequeue %v elements: %v", count, endTime.Sub(startTime))
	}

	queueDisk.Close()
}
