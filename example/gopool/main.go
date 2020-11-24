package main

import (
	"fmt"
	"soil/gopool"
	"time"
)

type Test1 struct{}

func(t Test1) Do() {
	//fmt.Println("test1...")
	time.Sleep(time.Millisecond*100)
	return
}

type Test2 struct {
	Number int64
}


func(t Test2)Do(){
	fmt.Println("---->number:",t.Number)
	time.Sleep(time.Millisecond*100)
	return
}

func Test(pool *gopool.Gopool) {

	for i:=0;i<10000;i++ {
		var t Test1
		pool.Push(t)
	}
	pool.Stop()
}

func Test11(pool *gopool.Gopool) {
	time.Sleep(time.Second)
}

func main() {
	pool := gopool.NewGopool(100)
	go Test(pool)
	go Test11(pool)
	pool.Loop()
}
