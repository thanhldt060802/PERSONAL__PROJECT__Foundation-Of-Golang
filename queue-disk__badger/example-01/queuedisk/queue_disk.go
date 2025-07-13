package queuedisk

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type QueueDisk struct {
	db      *badger.DB
	counter int64
}

type IQueueDisk interface {
	Enqueue(value string) error
	Dequeue() (string, error)
	Close() error
}

func NewQueueDisk(path string) IQueueDisk {
	opts := badger.DefaultOptions(path)
	// opts.WithSyncWrites(true)  // No effect on Window
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}

	qd := &QueueDisk{
		db:      db,
		counter: 0,
	}
	go qd.garbageCollection()

	return qd
}

func (qd *QueueDisk) garbageCollection() {
	if err := qd.db.RunValueLogGC(0.5); err != nil && err != badger.ErrNoRewrite {
		log.Printf("GC error: %v", err)
	}

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := qd.db.RunValueLogGC(0.5); err != nil && err != badger.ErrNoRewrite {
			log.Printf("GC error: %v", err)
		}
	}
}

func (qd *QueueDisk) Enqueue(value string) error {
	key := []byte(fmt.Sprintf("%020d", qd.counter))
	qd.counter++

	return qd.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, []byte(value))
	})
}

func (qd *QueueDisk) Dequeue() (string, error) {
	var keyToDelete []byte
	var value string

	err := qd.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.KeyCopy(nil)
			v, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			value = string(v)
			keyToDelete = k
			break
		}

		if keyToDelete == nil {
			return errors.New("queue empty")
		}

		return txn.Delete(keyToDelete)
	})

	return value, err
}

func (qd *QueueDisk) Close() error {
	return qd.db.Close()
}
