# GUIDES FOR GO-REDIS (REDIS)

## Guide 1.1 (See at guide-1.1)

### Guide for redis pub/sub using commonly used functions.

- <code>func redis.NewClient(opt *redis.Options) *redis.Client</code>: This function will initialize a <code>*redis.Client</code> with options, it is Redis client to Redis server.
- <code>func (c redis.cmdable) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd</code>: This function will publish message to channel.
- <code>func (c *redis.Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub</code>: This function will subscribe to channel.
- <code>func (c *redis.PubSub) Channel(opts ...redis.ChannelOption) <-chan *redis.Message</code>: This function will start a listener to channel for receiving messages.

## Guide xx (Coming soon...)