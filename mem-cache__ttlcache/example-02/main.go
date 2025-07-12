package main

import (
	"fmt"
	"thanhldt060802/memcache"
	"time"
)

var EXAMPLE_NUM int = 5
var EXAMPLES map[int]func()

var MemCache memcache.IMemCache[string, string]

func init() {
	EXAMPLES = map[int]func(){
		1: Example1,
		2: Example2,
		3: Example3,
		4: Example4,
		5: Example5,
		6: Example6,
	}

	MemCache = memcache.NewMemCache[string, string](3 * time.Second)
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
	MemCache.Set("my-key", "my-value")

	for {
		value, ok := MemCache.Get("my-key")
		if !ok {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", value)

		time.Sleep(1 * time.Second)
	}
}

/*
- Example for Set() and Get()

- Element is setted by Set() with default TTL is 3 seconds

- Get() will be called per 5 seconds (is greater defaultTTL)
*/
func Example2() {
	MemCache.Set("my-key", "my-value")

	for {
		value, ok := MemCache.Get("my-key")
		if !ok {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", value)

		time.Sleep(5 * time.Second)
	}
}

/*
- Test for SetTTL() and Get()

- Element is setted by SetTTL() with other TTL is 5 seconds

- Get() will be called per 1 second (is less than setted TTL)
*/
func Example3() {
	MemCache.SetTTL("my-key", "my-value", 5*time.Second)

	for {
		value, ok := MemCache.Get("my-key")
		if !ok {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", value)

		time.Sleep(1 * time.Second)
	}
}

/*
- Test for SetTTL() and Get()

- Element is setted by SetTTL() with other TTL is 6 seconds

- Get() will be called per 4 second (is less than setted TTL)

- It means the default TLL is no effect to other setted TTL for element
*/
func Example4() {
	MemCache.SetTTL("my-key", "my-value", 6*time.Second)

	for {
		value, ok := MemCache.Get("my-key")
		if !ok {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", value)

		time.Sleep(4 * time.Second)
	}
}

/*
- Test for SetTTL() and Get()

- Element is setted by SetTTL() with other TTL is 6 seconds

- Get() will be called per 10 second (is greater setted TTL)
*/
func Example5() {
	MemCache.SetTTL("my-key", "my-value", 6*time.Second)

	for {
		value, ok := MemCache.Get("my-key")
		if !ok {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", value)

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
	MemCache.Set("my-key", "my-value")
	count := 0

	for {
		value, ok := MemCache.Get("my-key")
		if !ok {
			fmt.Println("my-key doesn't exists in cache")
			return
		}
		fmt.Println("my-key :", value)

		count++
		if count == 5 {
			MemCache.Del("my-key")
		}
		time.Sleep(1 * time.Second)
	}
}
