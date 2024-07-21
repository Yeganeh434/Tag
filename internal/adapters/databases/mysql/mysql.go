package mysql

import (
	"log"
	"tag_project/internal/domain/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var TagDB Database

type Database struct {
	db *gorm.DB
}

func InitialDatabase() {
	dsn := "root:Yeganeh-2004@tcp(localhost:3306)/tag_project?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error connecting to the database")
	}
	TagDB.db = gormDB
	err=TagDB.db.AutoMigrate(&entity.Tag{}, &entity.Taxonomy{},&Counter{})
	if err != nil {
		log.Printf("error in migrating: %v", err)
		return
	}
}

type Counter struct{
	Count int `gorm:"primary_key"`
}

func (d *Database) GetCounter() (int,error) {
	var c Counter
	result:=d.db.First(&c)
	if result.Error!=nil {
		return 0,result.Error
	}
	if result.RowsAffected==0 {
		c.Count=1
		result=d.db.Create(&c)
		if result.Error!=nil{
			return 0 ,result.Error
		}
		return c.Count,nil
	}
	c.Count+=1
	result=d.db.Save(&c)
	if result.Error!=nil{
		return 0 ,result.Error
	}
	return c.Count,nil
}
