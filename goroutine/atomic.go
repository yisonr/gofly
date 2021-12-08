package main

/*
 * 原子操作是并发编程中"最小的且不可并行化"的操作, 通常如果多个并发体对同一个
 * 共享资源进行的操作是原子操作, 那么同一时刻最多只能有一个并发体对该资源进行
 * 操作; 从线程角度看, 在当前线程修改共享资源期间, 其他线程是不能访问该资源的,
 * 原子操作对多线程并发编程模型来说, 不会发生有别于单线程的意外情况, 共享资源
 * 的完整性可以得到保证.
 *
 * 一般, 原子操作都是通过"互斥"访问来保证的, 通常是由特殊的 CPU 指令提供保护,
 * 如果仅仅想模拟粗粒度的原子操作, 可以借助于 sync.Mutex 实现:
 */

import (
	"fmt"
	"sync"
)

var total struct {
	sync.Mutex
	value int
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i <= 100; i++ {
		total.Lock()
		total.value += i
		total.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go worker(&wg)
	go worker(&wg)
	wg.Wait()

	fmt.Println(total.value)
}