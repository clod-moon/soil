package gopool

type Worker struct {
  task_chan *chan Task
  stop_chan chan bool
  timeout int //毫秒
}

func NewWorker(task_chan *chan Task,timeout int) *Worker{
	return &Worker{
		task_chan:task_chan,
		stop_chan:make(chan bool),
		timeout:   timeout,
	}
}

func (w *Worker) Work() {
	for {
		Loop:
		select {
			case task,ok:= <- *w.task_chan:
				if !ok {
					goto End
				}
				task.Do()
				goto Loop
			case <- w.stop_chan:
					goto End
		default:
			break
		}
	}
	End:
}

