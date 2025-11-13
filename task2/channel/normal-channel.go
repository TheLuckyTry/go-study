package main

import (
	"fmt"
	"sync"
)

func send(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		ch <- i
		fmt.Printf("发送: %d\n", i)
	}
	close(ch)
}

func receive(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		fmt.Printf("接收到: %d\n", v)
	}
}

func main() {
	ch := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2) // 需要等待两个goroutine

	go send(ch, &wg)
	go receive(ch, &wg)

	wg.Wait()
}
