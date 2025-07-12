# GUIDES FOR TTLCACHE

## Example 01 (See at example-01)

### Guide for commonly used functions.

- <code>func New[K comparable, V any](opts ...Option[K, V]) *Cache[K, V]</code>: This function will initialize a <code>*Cache[K, V]</code> with options.
- <code>func (c *Cache[K, V]) Start()</code>: This function will start an automatic cleanup process that periodically deletes expired items.
- <code>func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) *Item[K, V]</code>: This function will set an element with <code>key: value</code> format and come with TTL.
- <code>func (c *Cache[K, V]) Get(key K, opts ...Option[K, V]) *Item[K, V]</code>: This function will get an element by key, it returns <code>*Item[K, V]</code> which contains key and value, and refresh the TTL of element.
- <code>func (c *Cache[K, V]) Delete(key K)</code>: This function will delete an element by key.

## Example 02 (See att example-02)

### Wrap for Example 01.