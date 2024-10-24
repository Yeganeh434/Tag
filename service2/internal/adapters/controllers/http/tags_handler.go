package http

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	nats "github.com/nats-io/nats.go"
)

type TagInfo struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

type CustomError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func EditTag(c *gin.Context) {
	var tag TagInfo
	err := c.BindJSON(&tag)
	if err != nil {
		log.Printf("error binding json:%v", err)
		c.Status(400)
		return
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Printf("error connecting NATS:%v", err)
		c.Status(400)
		return
	}
	defer nc.Close()

	// js, err := nc.JetStream()
	// if err != nil {
	// 	log.Printf("error connecting JetStream:%v", err)
	// 	c.Status(400)
	// 	return
	// }

	// _, err = js.AddStream(&nats.StreamConfig{
	// 	Name:     "EVENTS",
	// 	Subjects: []string{"events.*"},
	// })
	// if err != nil {
	// 	log.Printf("error setting stream configuration:%v", err)
	// 	c.Status(400)
	// 	return
	// }

	data, err := json.Marshal(tag)
	if err != nil {
		log.Printf("error marshaling:%v", err)
		c.Status(400)
		return
	}

	// In this mode the Request/Reply pattern is not used.(using Publish/Subscribe pattern)
	// err=nc.Publish("events.EditTag",data)
	// if err!=nil{
	// 	log.Printf("error publishing:%v",err)
	// 	c.Status(400)
	// 	return
	// }

	msg, err := nc.Request("EditTag", data, nats.DefaultTimeout) //jetstream:"events.EditTag"
	if err != nil {
		log.Printf("error sending request to NATS:%v", err)
		c.Status(400)
		return
	}

	var customError CustomError
	err = json.Unmarshal(msg.Data, &customError)
	if err != nil {
		log.Printf("error unmarshaling:%v", err)
		c.Status(400)
		return
	}
	
	if customError.Code == 400 {
		c.JSON(400, customError.Error)
		return
	}
	c.JSON(200, "tag edited successfully")
}
