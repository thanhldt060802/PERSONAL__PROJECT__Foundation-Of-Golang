# GUIDES FOR BADGER (QUEUE DISK)

## Guide 1.1 (See at guide-1.1)

### Guide for queue disk using commonly used functions.

- <code>func Open(opt Options) (*DB, error)</code>: This function will initialize a <code>*DB</code> with options.
- <code>func (db *badger.DB) RunValueLogGC(discardRatio float64) error</code>: This function will start a garbage collector to cleanup file storage if the proportion of data that is no longer useful is greater than <code>discardRatio</code>.
- <code>func (db *badger.DB) Update(fn func(txn *badger.Txn) error) error</code>: This function will start a transaction for a lot of action and commit changes.
- <code>func (txn *badger.Txn) Set(key []byte, val []byte) error</code>: This function will set an element with <code>key: value</code> format.
- <code>func (txn *badger.Txn) Delete(key []byte) error</code>: This function will delete an element by key.
- <code>func (txn *badger.Txn) NewIterator(opt badger.IteratorOptions) *badger.Iterator</code>: This function will create a iterator for retrieving data in (queue) disk.

## Guide 1.2 (See at guide-1.2)

### Guide for advanced queue disk. Base on Guide 1.1, we implement a batch queue disk to minimize disk reads/writes compared to queue disk of Guide 1.1.

## Guide x (Coming soon...)