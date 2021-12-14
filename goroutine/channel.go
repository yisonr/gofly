package main

import "fmt"

/*
 * 通道(channel)是在 goroutine 之间进行同步的主要方法, 在无缓冲的通道上的每一
 * 次发送都有与其对应的接收操作相匹配, 发送和接收操作通常发生在不同的 goroutine
 * 上(在同一个 goroutine 上执行两个操作很容易导致死锁);
 *
 */

// 在无缓冲的通道上的发送操作总在对应的接收操作完成前发生
var done = make(chan bool)
var msg string

func aGoroutine() {
	msg = "hello, world"
	done <- true // close(done)
	// 若在关闭通道后继续从中接收数据, 接收者会收到该通道返回的零值;
	// 所以在此例中可使用 close(done) 关闭通道代替 done <- true, 作用相同
}

func main1() {
	go aGoroutine()
	<-done
	fmt.Println(msg)
}

// 对于从无缓冲通道进行的接收, 发生在对该通道进行的发送完成之前
/* 基于以上规则, 交换两个 goroutine 中的接收和发送操作也是可以的(很危险)   */
func bGoroutine() {
	msg = "hello, world"
	<-done
}

func main2() {
	go bGoroutine()
	done <- true
	fmt.Println(msg)
}

// 也可以保证打印出 "hello, world", 因为在 main 中 done <- true 发送完成前,
// 后台 goroutine 的接收已经开始, 这保证了 msg 的赋值操作被执行; 也就是说
// 对无缓冲通道, 接收方和发送方都准备好后才会接收和发送; 但若该通道为带缓冲
// 的 done = make(bool, 1), 则 main 中的 done <- true 发送操作将不会被后台
// goroutine 的 <-done 接收操作阻塞, 该程序将无法打印出期望的结果.

/*
 * 对于带缓冲(缓冲大小为C)的通道, 对于通道中的第K个接收操作发生在第K+C个发送
 * 操作完成之前, 如果将C设置为0自然就对应无缓冲的通道, 也就是第K个接收完成在
 * 第K个发送完成之前, 因为无缓冲的通道只能同步发1个, 所以就简化为前面无缓冲
 * 通道的规则: 对于从无缓冲通道进行的接收, 发生在对该通道进行的发送完成之前.
 *
 */

// 可以通过控制通道的缓冲区的大小控制并发执行的 goroutine 的最大数目
var limit = make(chan int, 3)
var work []func()

func main() {
	for _, w := range work {
		go func(w func()) {
			limit <- 1
			w()
			<-limit
			// 在 goroutine 中对 limit 的发送和接收操作限制了同时执行的
			// goroutine 的数量
		}(w)
	}

	// <-make(chan int)
	// select {}
	for {
	}
	// select {} 是一个空的通道选择语句, 会导致死锁, <- make(chan int) 也一样;
	// for{} 会出现死循环

}
