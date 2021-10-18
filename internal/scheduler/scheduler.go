// Package scheduler @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/10/18 9:15 下午
package scheduler

import (
	"context"
	"time"

	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/comeonjy/go-kit/pkg/xsync"
	"github.com/google/wire"
	"spider/configs"
	"spider/internal/data"
)

var ProviderSet = wire.NewSet(NewScheduler)

type Scheduler struct {
	conf         configs.Interface
	logger       *xlog.Logger
	taskRepo     data.TaskRepo
	recordRepo   data.FetchRecordRepo
	resourceRepo data.ResourceRepo
}

func NewScheduler(cfg configs.Interface, logger *xlog.Logger, taskRepo data.TaskRepo, recordRepo data.FetchRecordRepo, resourceRepo data.ResourceRepo) *Scheduler {
	return &Scheduler{
		conf:         cfg,
		logger:       logger,
		taskRepo:     taskRepo,
		recordRepo:   recordRepo,
		resourceRepo: resourceRepo,
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
	workChan := make(chan Request)
	out := make(chan Resource)
	finish := make(chan struct{}, 1)
	finish <- struct{}{}

	var err error

	xsync.NewGroup().Go(func(ctx context.Context) error {
		return s.Save(ctx, out)
	})

	xsync.NewGroup().Go(func(ctx context.Context) error {
		return s.Worker(ctx, workChan, out)
	})

	ticker := time.NewTicker(time.Second * 10)
	for {
		ctx, cancel := context.WithCancel(context.Background())
		g := xsync.NewGroup(xsync.WithContext(ctx))
		select {
		case <-ticker.C:
			if !s.conf.Get().RunSpider {
				cancel()
				continue
			}
		case <-finish:
			ctx, cancel = context.WithCancel(context.Background())

			one := &data.TaskModel{}
			one, err = s.taskRepo.TakeOne(ctx)
			if err != nil {
				s.logger.Error(ctx, err.Error())
				time.Sleep(time.Second * 10)
				continue
			}
			if one.State == data.TaskStateNormal {
				workChan <- Request{
					TaskUUID: one.UUID,
					Url:      one.Entrance,
				}
			}
			g.Go(func(ctx context.Context) error {
				return s.Scan(ctx, one, workChan, finish)
			})
		}
	}
}

func (s *Scheduler) Worker(ctx context.Context, workChan chan Request, out chan Resource) error {
	for v := range workChan {
		time.Sleep(time.Second)
		out <- Resource{
			Request: v,
			Data:    "test",
			Urls:    []string{v.Url},
		}
	}
	return nil
}

func (s *Scheduler) Save(ctx context.Context, out chan Resource) error {
	for v := range out {
		if !IsExist(v.Url) {
			if err := s.resourceRepo.Insert(ctx, &data.ResourceModel{
				TaskUUID: v.TaskUUID,
				Url:      v.Url,
				Content:  v.Data,
			}); err != nil {
				s.logger.Error(ctx, err.Error())
			}
			records := make([]data.FetchRecordModel, 0)
			for _, url := range v.Urls {
				records = append(records, data.FetchRecordModel{
					TaskUUID: v.TaskUUID,
					Url:      url,
					State:    data.FetchStateNormal,
				})
			}
			if err := s.recordRepo.BatchCreate(ctx, records); err != nil {
				s.logger.Error(ctx, err.Error())
			}
		}
	}
	return nil
}

func IsExist(url string) bool {
	return false
}

func (s *Scheduler) Scan(ctx context.Context, task *data.TaskModel, workChan chan Request, finish chan struct{}) error {
	for {
		finishTime := time.Now()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			list, err := s.recordRepo.Scan(ctx, task.FetchOffset, 10)
			if err != nil {
				s.logger.Error(context.TODO(), err.Error())
				time.Sleep(time.Second * 10)
				continue
			}
			if len(list) == 0 {
				if time.Now().Sub(finishTime) > time.Minute {
					finish <- struct{}{}
					return nil
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
				if err := s.taskRepo.SetOffset(ctx, v.Id); err != nil {
					s.logger.Error(ctx, err.Error())
					continue
				}
				task.FetchOffset = v.Id
				workChan <- Request{
					TaskUUID: v.TaskUUID,
					Url:      v.Url,
				}
			}
		}
	}

}
