package data

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type ResourceModel struct {
	Id        uint64 `gorm:"primarykey"`
	TaskUUID  string `gorm:"type:varchar(36);not null"`
	Url       string `gorm:"type:varchar(200);not null"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ResourceModel) TableName() string {
	return "resource"
}

type ResourceRepo interface {
	Get(ctx context.Context, id int) (*ResourceModel, error)
	Insert(ctx context.Context, resource *ResourceModel) error
}

func NewResourceRepo(data *Data) ResourceRepo {
	return &resourceRepo{db: data.db}
}

type resourceRepo struct {
	db *gorm.DB
}

func (r resourceRepo) Get(ctx context.Context, id int) (*ResourceModel, error) {
	return &ResourceModel{}, nil
}

func (r resourceRepo) Insert(ctx context.Context, resource *ResourceModel) error {
	return r.db.Create(&resource).Error
}
