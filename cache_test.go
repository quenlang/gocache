package gocache

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

type customer struct {
	Name string
	Age  int
}

func TestCacheDefault(t *testing.T) {
	cache, err := New()
	if err != nil {
		fmt.Println(err)
	}
	err = cache.Set("c1", customer{Name: "c1", Age: 1})
	if err != nil {
		fmt.Println(err)
	}
	err = cache.Set("c2", customer{Name: "c2", Age: 2})
	if err != nil {
		fmt.Println(err)
	}
	c1, err := cache.Get("c1")
	if err != nil {
		fmt.Println(err)
	}
	if c1cast, ok := c1.(customer); ok {
		fmt.Println(c1cast)
	}
	if cache.Delete("c2") != nil {
		fmt.Println(err)
	}
	c2, err := cache.Get("c2")
	if err != nil {
		fmt.Println(err)
	}
	if c2cast, ok := c2.(customer); ok {
		fmt.Println(c2cast)
	}
	if cache.Set("n1", 1) != nil {
		fmt.Println(err)
	}
	if cache.Set("s1", "s1") != nil {
		fmt.Println(err)
	}
	n1, err := cache.Get("n1")
	if err != nil {
		fmt.Println(err)
	}
	if n1cast, ok := n1.(int); ok {
		fmt.Println(n1cast)
	}
	s1, err := cache.Get("s1")
	if err != nil {
		fmt.Println(err)
	}
	if s1cast, ok := s1.(string); ok {
		fmt.Println(s1cast)
	}
}

func TestCacheExist(t *testing.T) {
	sizeOps := WithSizeInMB(2048)
	ttlOps := WithTTL(48 * time.Hour)
	cleanFreqOps := WithCleanFrequency(1 * time.Hour)

	cache, err := New(sizeOps, ttlOps, cleanFreqOps)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	cache.Set("zhoujk", 1)
	fmt.Println(cache.Exist("zhoujk"))
	fmt.Println(cache.Exist("zhoujk1"))
}

func TestCacheWithSetup(t *testing.T) {
	sizeOps := WithSizeInMB(2048)
	ttlOps := WithTTL(48 * time.Hour)
	cleanFreqOps := WithCleanFrequency(1 * time.Hour)

	cache, err := New(sizeOps, ttlOps, cleanFreqOps)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	err = cache.Set("c1", customer{Name: "c1", Age: 1})
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	err = cache.Set("c1", customer{Name: "c1", Age: 11})
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	err = cache.Set("c2", customer{Name: "c2", Age: 2})
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	c1, err := cache.Get("c1")
	if err != nil {
		fmt.Printf("get c1 error: %v\n", err)
	}
	if c1cast, ok := c1.(customer); ok {
		fmt.Println(c1cast)
	}
	if cache.Delete("c2") != nil {
		fmt.Printf("error: %v\n", err)
	}
	c2, err := cache.Get("c2")
	if err != nil {
		fmt.Printf("get c2 error: %v\n", err)
	}
	if c2cast, ok := c2.(customer); ok {
		fmt.Println(c2cast)
	}
	if cache.Set("n1", 1) != nil {
		fmt.Printf("error: %v\n", err)
	}
	if cache.Set("s1", "s1") != nil {
		fmt.Printf("error: %v\n", err)
	}
	n1, err := cache.Get("n1")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	if n1cast, ok := n1.(int); ok {
		fmt.Println(n1cast)
	}
	s1, err := cache.Get("s1")
	if err != nil {
		fmt.Println(err)
	}
	if s1cast, ok := s1.(string); ok {
		fmt.Println(s1cast)
	}
}

func BenchmarkName(b *testing.B) {
	cache, err := New()
	if err != nil {
		fmt.Println(err)
	}
	start := time.Now()
	for i := 0; i < b.N; i++ {
		if err := cache.Set("key"+strconv.Itoa(i), i); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("set %v items cost %v\n", b.N, time.Since(start))
}

func Test(t *testing.T) {
	cache, err := New()
	if err != nil {
		fmt.Println(err)
	}

	start := time.Now()
	for i := 0; i < 1000000; i++ {
		if err := cache.Set("key"+strconv.Itoa(i), i); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("set 1000000 items cost %v\n", time.Since(start))
	var producerCh = make(chan int, 1)
	var consumerCh = make(chan int, 1)
	go func() {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				start := time.Now()
				for n := 900000 + i*10000; n < 900000+(i+1)*10000; n++ {
					if _, err := cache.Get("key" + strconv.Itoa(i)); err != nil {
						fmt.Println(err)
					}
				}
				fmt.Printf("get 10000 items cost %v\n", time.Since(start))
			}(i)
		}
		wg.Wait()
		consumerCh <- 1
	}()

	go func() {
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				start := time.Now()
				for n := 1000000 + i*100000; n < 1000000+(i+1)*100000; n++ {
					if err := cache.Set("key"+strconv.Itoa(i), n); err != nil {
						fmt.Println(err)
					}
				}
				fmt.Printf("set 100000 items cost %v\n", time.Since(start))
			}(i)
		}
		wg.Wait()
		producerCh <- 1
	}()
	<-consumerCh
	<-producerCh
}
