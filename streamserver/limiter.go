package main

import "log"

/*
	每当有一个request访问了路由的handler，就会新建一个goroutine
	所以token bucket中的计数不能是简单的变量，得是线程安全才行
	在这里不使用共享内存加锁的方法，而是用共享channel
 */


 /*
 	有时候多个request同时连接，超过bucketsize，但没出现Reached the rate limitation.
 	这个跟H5在浏览器的处理方式有关系，如果默认缓存开启的话，浏览器会自动去缓存中找已经缓冲好的视频。
 	因此直接会结束你的新request，但是确实是在播放。
  */
type ConnLimiter struct {
	bucketSize int
	bucketChan chan int
}

//相当于构造函数
func NewConnLimiter(bs int) *ConnLimiter {
	return &ConnLimiter {
		bucketSize: bs,
		bucketChan: make(chan int, bs),
	}
}

func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucketChan) >= cl.bucketSize {
		log.Printf("Reached the rate limitation.")
		return false
	}
	cl.bucketChan <- 1 //随便写个数字
	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	c := <- cl.bucketChan
	log.Printf("Released: %d", c)
}