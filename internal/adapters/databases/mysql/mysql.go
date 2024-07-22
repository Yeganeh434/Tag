package mysql

import (
	"log"
	"tag_project/internal/domain/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var TagDB Database

type Database struct {
	DB *gorm.DB
}

func InitialDatabase() {
	dsn := "root:Yeganeh-2004@tcp(localhost:3306)/tag_project?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error connecting to the database")
	}
	TagDB.DB = gormDB
	err=TagDB.DB.AutoMigrate(&entity.Tag{}, &entity.Taxonomy{})
	if err != nil {
		log.Printf("error in migrating: %v", err)
		return
	}
}

