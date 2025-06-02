package entity

import (
	"gorm.io/gorm"
	"time"
)

type Result struct {
	Id           int            `gorm:"column:id;type:int;primaryKey;autoIncrement;not null"`
	ProjectId    int            `gorm:"column:project_id;type:int"`
	ProjectKey   string         `gorm:"column:project_key;type:varchar(255)"`
	Rule         string         `gorm:"column:rule;type:varchar(255)"`
	Path         string         `gorm:"column:path;type:varchar(255)"`
	Line         string         `gorm:"column:line;type:varchar(255)"`
	Title        string         `gorm:"column:title;type:text"`
	Description  string         `gorm:"column:description;type:text"`
	Severity     string         `gorm:"column:severity;type:varchar(255)"`
	Type         string         `gorm:"column:type;type:varchar(255)"`
	References   string         `gorm:"column:references;type:text"`
	LastFoundAt  string         `gorm:"column:last_found_at;type:varchar(255)"`
	StatusResult int            `gorm:"column:status_result;type:int"`
	ScanVersion  int            `gorm:"column:scan_version;type:int"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;->"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;default:null;->"`
}

func (Result) TableName() string {
	return "results"
}
