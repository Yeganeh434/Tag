package usecases

import (
	"log"
	"time"

	"github.com/sony/sonyflake"
)

var flake *sonyflake.Sonyflake

func init() {
	setting := sonyflake.Settings{
		StartTime: time.Now(),
	}
	flake = sonyflake.NewSonyflake(setting)
	if flake == nil {
		log.Printf("sonyflake not created")
		return
	}
}
func GenerateID() (uint64, error) {
	id, err := flake.NextID()
	if err != nil {
		return 0, err
	}
	return id, nil
}
