package scheduler

import (
	log "github.com/sirupsen/logrus"

	"spider/engine"
)

type QueuedScheduler struct {
	requestChan chan engine.Resource
	workChan    chan chan engine.Resource
}

func (q *QueuedScheduler) WorkChan() chan engine.Resource {
	in := make(chan engine.Resource)
	return in
}

func (q *QueuedScheduler) WorkReady(in chan engine.Resource) {
	q.workChan <- in
}

func (q *QueuedScheduler) Submit(r engine.Resource) {
	q.requestChan <- r
}



func (q *QueuedScheduler) Worker(out chan engine.Result)  {
	in:=q.WorkChan()
	go func() {
		for {
			q.WorkReady(in)
			resource := <-in
			body, err := engine.Download(resource.Url)
			if err != nil {
				log.Error(err)
			}
			result, err := resource.FetchFunc(body)
			if err != nil {
				log.Error(err)
			}
			out <- result
		}
	}()
}

func (q *QueuedScheduler) Run() {
	reqQueued := make([]engine.Resource, 0)
	workQueued := make([]chan engine.Resource, 0)

	go func() {
		for {
			activeReq := engine.Resource{}
			activeWork := make(chan engine.Resource)
			if len(reqQueued) >= 1 && len(workQueued) >= 1 {
				activeReq = reqQueued[0]
				activeWork = workQueued[0]

			}
			select {
			case r := <-q.requestChan:
				reqQueued = append(reqQueued, r)
			case w := <-q.workChan:
				workQueued = append(workQueued, w)
			case activeWork <- activeReq:
				reqQueued = reqQueued[1:]
				workQueued = workQueued[1:]
			}
		}

	}()

}
