package data

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type TaskModel struct {
	Id          uint64    `gorm:"primarykey"`
	UUID        string    `gorm:"type:varchar(36);not null"`
	Name        string    `gorm:"type:varchar(50);not null"`
	UserUUID    string    `gorm:"type:varchar(36);not null"`
	Entrance    string    `gorm:"type:varchar(200);not null"`
	FetchOffset uint64    `gorm:"type:uint"`
	State       TaskState `gorm:"type:uint"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type TaskState uint64

const (
	TaskStateNormal  TaskState = 1
	TaskStateWorking TaskState = 2
	TaskStatePause   TaskState = 3
	TaskStateFinish  TaskState = 4
)

func (TaskModel) TableName() string {
	return "tasks"
}

type TaskRepo interface {
	Get(ctx context.Context, id int) (*TaskModel, error)
	TakeOne(ctx context.Context) (*TaskModel, error)
	TakeN(ctx context.Context, num int) ([]TaskModel, error)
	SetOffset(ctx context.Context, taskUUID string, offset uint64) error
	UpdateState(ctx context.Context, taskUUID string, state TaskState) error
}

func NewTaskRepo(data *Data) TaskRepo {
	return &taskRepo{db: data.db}
}

type taskRepo struct {
	db *gorm.DB
}

func (r taskRepo) Get(ctx context.Context, id int) (*TaskModel, error) {
	return &TaskModel{}, nil
}

func (r taskRepo) TakeOne(ctx context.Context) (*TaskModel, error) {
	task := TaskModel{}
	err := r.db.Model(&TaskModel{}).Where("state in (?,?)", TaskStateNormal, TaskStateWorking).First(&task).Error
	return &task, err
}
func (r taskRepo) TakeN(ctx context.Context, num int) ([]TaskModel, error) {
	tasks := make([]TaskModel, 0)
	err := r.db.Model(&TaskModel{}).Where("state in (?,?)", TaskStateNormal, TaskStateWorking).Order("id desc").Find(&tasks).Error
	return tasks, err
}
func (r taskRepo) SetOffset(ctx context.Context, taskUUID string, offset uint64) error {
	return r.db.Model(&TaskModel{}).Where("uuid", taskUUID).Update("fetch_offset", offset).Error
}
func (r taskRepo) UpdateState(ctx context.Context, taskUUID string, state TaskState) error {
	return r.db.Model(&TaskModel{}).Where("uuid = ?", taskUUID).Update("state", state).Error
}
