package migration

import (
	"gorm.io/gorm"
	"sast-integration/app/dbo/entity"
)

func DoMigration(db *gorm.DB) {
	db.AutoMigrate(&entity.Project{})
	db.AutoMigrate(&entity.ProjectAuth{})
}
