package gopool

import (
	"sync"
)

type Gopool struct {
	WorkNumber int
	task_chan chan Task
	Workers []*Worker
	wg sync.WaitGroup
}

func NewGopool(worker_number int) *Gopool {
	return &Gopool{WorkNumber:worker_number,task_chan:make(chan Task,10000)}
}

func (g *Gopool) Loop() {
	for i:=0;i<g.WorkNumber;i++ {
		g.Workers= append(g.Workers,NewWorker(&g.task_chan,0))
	}
	for _,v := range  g.Workers {
		g.wg.Add(1)
		go func(worker *Worker) {
			worker.Work()
			g.wg.Done()
		}(v)
	}

	g.wg.Wait()
}

func (g *Gopool)Stop () {
	close(g.task_chan)
}

func (g *Gopool)StopNow() {
	for _,v:= range g.Workers {
		v.stop_chan <- true
	}
}

func (g *Gopool)Push(t Task) {
	g.task_chan<- t
}