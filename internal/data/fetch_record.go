package data

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type FetchRecordModel struct {
	Id        uint64     `gorm:"primarykey"`
	TaskUUID  string     `gorm:"type:varchar(36);not null"`
	Url       string     `gorm:"type:varchar(200);not null"`
	State     FetchState `gorm:"type:uint"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type FetchState uint64

const (
	FetchStateNormal  FetchState = 1
	FetchStateParsing FetchState = 2
	FetchStateFinish  FetchState = 3
)

func (FetchRecordModel) TableName() string {
	return "fetch_records"
}

type FetchRecordRepo interface {
	Exist(ctx context.Context, taskUUID string, url string) (bool, error)
	Scan(ctx context.Context, offsetID uint64, limit int) ([]FetchRecordModel, error)
	UpdateState(ctx context.Context, recordID uint64, state FetchState) error
	BatchCreate(ctx context.Context, list []FetchRecordModel) error
}

func NewFetchRecordRepo(data *Data) FetchRecordRepo {
	return &fetchRecordRepo{db: data.db}
}

type fetchRecordRepo struct {
	db *gorm.DB
}

func (r fetchRecordRepo) Exist(ctx context.Context, taskUUID string, url string) (bool, error) {
	var count int64
	err := r.db.Model(&FetchRecordModel{}).Where("task_uuid = ?", taskUUID).Where("url = ?", url).Count(&count).Error
	return count > 0, err
}

func (r fetchRecordRepo) Scan(ctx context.Context, offsetID uint64, limit int) ([]FetchRecordModel, error) {
	records := make([]FetchRecordModel, 0)
	err := r.db.Model(&FetchRecordModel{}).Where("id > ?", offsetID).Limit(limit).Find(&records).Error
	return records, err
}

func (r fetchRecordRepo) UpdateState(ctx context.Context, recordID uint64, state FetchState) error {
	err := r.db.Model(&FetchRecordModel{}).Where("id = ?", recordID).Update("state", state).Error
	return err
}

func (r fetchRecordRepo) BatchCreate(ctx context.Context, list []FetchRecordModel) error {
	return r.db.Create(&list).Error
}
