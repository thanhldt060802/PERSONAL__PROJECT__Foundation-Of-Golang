package queue

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type BatchQueueDisk struct {
	db      *badger.DB
	counter int64

	batchSize int

	batchEnqueue            []string
	currentBatchEnqueueSize int

	batchDequeue []string
}

func NewBatchQueueDisk(path string, batchSize int) *BatchQueueDisk {
	opts := badger.DefaultOptions(path)
	// opts.WithSyncWrites(true)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	bqd := &BatchQueueDisk{
		db:      db,
		counter: 0,

		batchSize: batchSize,

		batchEnqueue:            make([]string, batchSize),
		currentBatchEnqueueSize: 0,

		batchDequeue: make([]string, batchSize),
	}
	go bqd.GarbageCollection()

	return bqd
}

func (bqd *BatchQueueDisk) GarbageCollection() {
	if err := bqd.db.RunValueLogGC(0.5); err != nil && err != badger.ErrNoRewrite {
		log.Printf("GC error: %v", err)
	}

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := bqd.db.RunValueLogGC(0.5); err != nil && err != badger.ErrNoRewrite {
			log.Printf("GC error: %v", err)
		}
	}
}

func (bqd *BatchQueueDisk) Close() error {
	return bqd.db.Close()
}

func (bqd *BatchQueueDisk) Enqueue(value string) error {
	bqd.batchEnqueue[bqd.currentBatchEnqueueSize] = value
	bqd.currentBatchEnqueueSize++

	if bqd.currentBatchEnqueueSize >= bqd.batchSize {
		return bqd.db.Update(func(txn *badger.Txn) error {
			for _, value := range bqd.batchEnqueue {
				key := []byte(fmt.Sprintf("%020d", bqd.counter))
				bqd.counter++

				if err := txn.Set(key, []byte(value)); err != nil {
					return err
				}
			}

			bqd.currentBatchEnqueueSize = 0

			return nil
		})
	}

	return nil
}

func (bqd *BatchQueueDisk) Dequeue() ([]string, error) {
	var keysToDelete [][]byte
	var values []string

	err := bqd.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		currentBatchDequeueSize := 0
		for it.Rewind(); it.Valid() && currentBatchDequeueSize < bqd.batchSize; it.Next() {
			item := it.Item()
			k := item.KeyCopy(nil)
			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			bqd.batchDequeue[currentBatchDequeueSize] = string(v)
			keysToDelete = append(keysToDelete, k)

			currentBatchDequeueSize++
		}
		values = bqd.batchDequeue[:currentBatchDequeueSize]

		if len(values) == 0 {
			return errors.New("queue empty")
		}

		for _, key := range keysToDelete {
			if err := txn.Delete(key); err != nil {
				return err
			}
		}

		return nil
	})

	return values, err
}
