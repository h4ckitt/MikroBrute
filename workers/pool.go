package workers

/*import "fmt"

type pool struct {
	jobs    <-chan string
	workers []*worker
	done    chan *worker
}

func NewPool(num int, jobsChan <-chan string) *pool {
	p := &pool{
		jobs:    jobsChan,
		workers: make([]*worker, num),
		done:    make(chan *worker),
	}
	for i := 1; i <= num; i++ {
		w := worker{
			name:     fmt.Sprintf("Worker %d", i),
			lastHash: "",
			errorJob: "",
		}

		p.workers = append(p.workers, &w)
	}

	return p
}

func (p *pool) Start() {
	go func() {
		for {
			select {
			case job := <-p.jobs:
				if len(p.workers) == 0 {
					continue
				}
				worker := p.workers[0]
				if len(p.workers) > 1 {
					p.workers = p.workers[1:]
				}
				go worker.ListenAndExecute(job, p.done)

			case worker := <-p.done:
				p.workers = append(p.workers, worker)
			}
		}
	}()
}
*/
