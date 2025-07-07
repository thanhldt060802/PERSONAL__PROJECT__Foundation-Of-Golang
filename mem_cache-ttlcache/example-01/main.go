package main

import (
	"fmt"
	"thanhldt060802/cache"
	"time"
)

var MemCache cache.IMemCache[string, string] = cache.NewMemCache[string, string](3 * time.Second)

func main() {

	test1()
	// test2()
	// test3()
	// test4()
	// test5()
	// test6()

}

func test1() {
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

func test2() {
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

func test3() {
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

func test4() {
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

func test5() {
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

func test6() {
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
