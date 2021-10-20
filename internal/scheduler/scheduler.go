// Package scheduler @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/10/18 9:15 下午
package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/comeonjy/go-kit/pkg/xsync"
	"github.com/google/wire"
	"spider/configs"
	"spider/internal/data"
	"spider/internal/scheduler/worker"
)

var ProviderSet = wire.NewSet(NewScheduler)

type Scheduler struct {
	conf         configs.Interface
	logger       *xlog.Logger
	taskRepo     data.TaskRepo
	recordRepo   data.FetchRecordRepo
	resourceRepo data.ResourceRepo
	workChan     chan Request
	out          chan Resource
	concurrent   int
}

func NewScheduler(cfg configs.Interface, logger *xlog.Logger, taskRepo data.TaskRepo, recordRepo data.FetchRecordRepo, resourceRepo data.ResourceRepo) *Scheduler {
	return &Scheduler{
		conf:         cfg,
		logger:       logger,
		taskRepo:     taskRepo,
		recordRepo:   recordRepo,
		resourceRepo: resourceRepo,
		workChan:     make(chan Request),
		out:          make(chan Resource),
		concurrent:   2,
	}
}

type Request struct {
	TaskUUID string `json:"task_uuid"`
	Url      string `json:"url"`
}

type Resource struct {
	Request
	Data string   `json:"data"`
	Urls []string `json:"urls"`
}

func (s *Scheduler) Run() error {
	finish := make(chan struct{})

	xsync.NewGroup().Go(func(ctx context.Context) error {
		return s.Save(ctx)
	})

	xsync.NewGroup().Go(func(ctx context.Context) error {
		finish <- struct{}{}
		return s.Worker(ctx)
	})

	ticker := time.NewTicker(time.Second * 10)
	var ctx context.Context
	var cancel context.CancelFunc
	var flag bool
	for {
		select {
		case <-ticker.C:
			flag = s.conf.Get().RunSpider == "true"
			if !flag {
				if cancel != nil {
					cancel()
					flag = false
				}
				continue
			}
			flag = true
		case <-finish:
			ctx, cancel = context.WithCancel(context.Background())
			xsync.NewGroup(xsync.WithContext(ctx)).Go(func(ctx context.Context) error {
				if err := s.Start(ctx); err != nil {
					s.logger.Error(ctx, err.Error())
					cancel()
					return err
				}
				for {
					if flag {
						break
					}
					time.Sleep(time.Second)
				}
				finish <- struct{}{}
				return nil
			})
		}
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	tasks, err := s.taskRepo.TakeN(ctx, s.concurrent)
	if err != nil {
		s.logger.Error(ctx, err.Error())
		return err
	}
	g := xsync.NewGroup(xsync.WithContext(ctx))
	for _, one := range tasks {
		if one.State == data.TaskStateNormal {
			s.workChan <- Request{
				TaskUUID: one.UUID,
				Url:      one.Entrance,
			}
			if err := s.taskRepo.UpdateState(ctx, one.UUID, data.TaskStateWorking); err != nil {
				return err
			}
		}
		g.Go(func(ctx context.Context) error {
			return s.Scan(ctx, one)
		})
	}
	g.Wait()
	return nil
}

func (s *Scheduler) Worker(ctx context.Context) error {
	for v := range s.workChan {
		info, urls, err := worker.Work(v.Url)
		if err != nil {
			return err
		}
		s.out <- Resource{
			Request: v,
			Data:    info,
			Urls:    urls,
		}
		log.Println("fetch :", v.Url)
	}
	return nil
}

func (s *Scheduler) Save(ctx context.Context) error {
	for v := range s.out {

		if err := s.resourceRepo.Insert(ctx, &data.ResourceModel{
			TaskUUID: v.TaskUUID,
			Url:      v.Url,
			Content:  v.Data,
		}); err != nil {
			s.logger.Error(ctx, err.Error())
		}
		records := make([]data.FetchRecordModel, 0)
		for _, url := range v.Urls {
			if !s.IsExist(v.TaskUUID, v.Url) {
				records = append(records, data.FetchRecordModel{
					TaskUUID: v.TaskUUID,
					Url:      url,
					State:    data.FetchStateNormal,
				})
			}
		}
		if err := s.recordRepo.BatchCreate(ctx, records); err != nil {
			s.logger.Error(ctx, err.Error())
		}

	}
	return nil
}

func (s *Scheduler) IsExist(taskUUID string, url string) bool {
	isExist, _ := s.recordRepo.Exist(context.Background(), taskUUID, url)
	return isExist
}

func (s *Scheduler) Scan(ctx context.Context, task data.TaskModel) error {
	for {
		finishTime := time.Now()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			list, err := s.recordRepo.Scan(ctx, task.FetchOffset, 10)
			if err != nil {
				s.logger.Error(ctx, err.Error())
				time.Sleep(time.Second * 10)
				continue
			}
			if len(list) == 0 {
				if time.Now().Sub(finishTime) > time.Minute {
					return s.taskRepo.UpdateState(ctx, task.UUID, data.TaskStateFinish)
				}
				time.Sleep(time.Second * 10)
				continue
			} else {
				finishTime = time.Now()
			}
			for _, v := range list {
				if err := s.recordRepo.UpdateState(ctx, v.Id, data.FetchStateParsing); err != nil {
					s.logger.Error(ctx, err.Error())
					continue
				}
				if err := s.taskRepo.SetOffset(ctx, v.TaskUUID, v.Id); err != nil {
					s.logger.Error(ctx, err.Error())
					continue
				}
				task.FetchOffset = v.Id
				s.workChan <- Request{
					TaskUUID: v.TaskUUID,
					Url:      v.Url,
				}
			}
		}
	}

}
