package main

import (
	"fmt"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var EXAMPLE_NUM int = 2
var EXAMPLES map[int]func()

var MemCache *ttlcache.Cache[string, string]

func init() {
	EXAMPLES = map[int]func(){
		1: Example1,
		2: Example2,
		3: Example3,
		4: Example4,
		5: Example5,
		6: Example6,
	}

	MemCache = ttlcache.New(
		ttlcache.WithTTL[string, string](3 * time.Second),
	)
	go MemCache.Start()
}

func main() {

	EXAMPLES[EXAMPLE_NUM]()

}

/*
- Example for Set() and Get()

- Element is setted by Set() with default TTL is 3 seconds

- Get() will be called per 1 second (is less than defaultTTL)
*/
func Example1() {
	MemCache.Set("my-key", "my-value", ttlcache.DefaultTTL)

	for {
		item := MemCache.Get("my-key")
		if item == nil {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", item.Value())

		time.Sleep(1 * time.Second)
	}
}

/*
- Example for Set() and Get()

- Element is setted by Set() with default TTL is 3 seconds

- Get() will be called per 5 seconds (is greater defaultTTL)
*/
func Example2() {
	MemCache.Set("my-key", "my-value", ttlcache.DefaultTTL)

	for {
		item := MemCache.Get("my-key")
		if item == nil {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", item.Value())

		time.Sleep(5 * time.Second)
	}
}

/*
- Test for Set() and Get()

- Element is setted by Set() with other TTL is 5 seconds

- Get() will be called per 1 second (is less than setted TTL)
*/
func Example3() {
	MemCache.Set("my-key", "my-value", 5*time.Second)

	for {
		item := MemCache.Get("my-key")
		if item == nil {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", item.Value())

		time.Sleep(1 * time.Second)
	}
}

/*
- Test for Set() and Get()

- Element is setted by Set() with other TTL is 6 seconds

- Get() will be called per 4 second (is less than setted TTL)

- It means the default TLL is no effect to other setted TTL for element
*/
func Example4() {
	MemCache.Set("my-key", "my-value", 6*time.Second)

	for {
		item := MemCache.Get("my-key")
		if item == nil {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", item.Value())

		time.Sleep(4 * time.Second)
	}
}

/*
- Test for Set() and Get()

- Element is setted by Set() with other TTL is 6 seconds

- Get() will be called per 10 second (is greater setted TTL)
*/
func Example5() {
	MemCache.Set("my-key", "my-value", 6*time.Second)

	for {
		item := MemCache.Get("my-key")
		if item == nil {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", item.Value())

		time.Sleep(10 * time.Second)
	}
}

/*
- Test for Set(), Get() and Del()

- Element is setted by Set() with default TTL is 3 seconds

- Get() will be called per 1 second (is less than defaultTTL)

- When Get() is called 5 times, the element will be deleted
*/
func Example6() {
	MemCache.Set("my-key", "my-value", ttlcache.DefaultTTL)
	count := 0

	for {
		item := MemCache.Get("my-key")
		if item == nil {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", item.Value())

		count++
		if count == 5 {
			MemCache.Delete("my-key")
		}
		time.Sleep(1 * time.Second)
	}
}
