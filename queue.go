package gojobqueue

type Queue chan job

type job struct {
	transact func() error
	rollback func(error)
}

func (q Queue) AddJob(transact func() error, rollback func(error)) (err error) {
	j := job{transact, rollback}
	defer func() {
		r := recover()
		if r != nil {
			err = r.(error)
		}
	}()
	q <- j
	return
}

func (q Queue) Close() {
	close(q)
}

func (q Queue) StartWorking() {
	go workJobs(q)
}

func workJobs(q Queue) {
	for j := range q {
		err := j.transact()
		if err != nil {
			j.rollback(err)
		}
	}
}