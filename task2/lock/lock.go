package main

import (
	"fmt"
	"sync"
)

/*
*
编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ： sync.Mutex 的使用、并发数据安全。
*/
func incrementNum(mu *sync.Mutex, num *int) {
	for i := 0; i < 1000; i++ {
		mu.Lock()
		(*num)++
		mu.Unlock()
	}
}
func main() {
	num := 0
	mu := sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			incrementNum(&mu, &num)
		}()
	}
	wg.Wait()
	fmt.Println(num)
}
